package controller

import (
	"encoding/json"
	"fmt"
	"github.com/fajar-andriansyah/loan-engine/internal/helpers"
	"github.com/go-chi/chi/v5"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fajar-andriansyah/loan-engine/infrastructure/http/middleware"
	"github.com/fajar-andriansyah/loan-engine/models"
	"github.com/fajar-andriansyah/loan-engine/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type LoanController struct {
	loanUsecase usecase.LoanUsecase
	validator   *validator.Validate
}

func NewLoanController(loanUsecase usecase.LoanUsecase) *LoanController {
	return &LoanController{
		loanUsecase: loanUsecase,
		validator:   validator.New(),
	}
}

func (c *LoanController) CreateLoanProposal(w http.ResponseWriter, r *http.Request) {
	var req models.CreateLoanRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		c.sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// Validate request
	if err := c.validator.Struct(&req); err != nil {
		log.Error().Err(err).Msg("Validation failed")
		c.sendValidationErrorResponse(w, err)
		return
	}

	user, err := middleware.GetUserFromCtx(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user from context")
		c.sendErrorResponse(w, http.StatusUnauthorized, "User context not found", nil)
		return
	}

	if user.UserType != "borrower" {
		log.Error().Str("user_type", user.UserType).Msg("Invalid user type for loan creation")
		c.sendErrorResponse(w, http.StatusForbidden, "Only borrowers can create loan proposals", nil)
		return
	}

	response, err := c.loanUsecase.CreateLoanProposal(r.Context(), &req, user.UserID)
	if err != nil {
		log.Error().Err(err).Str("borrower_id", user.UserID).Msg("Failed to create loan proposal")

		if err.Error() == "borrower not found" {
			c.sendErrorResponse(w, http.StatusNotFound, "Borrower not found", map[string]string{
				"error_code": "BORROWER_NOT_FOUND",
			})
		} else {
			c.sendErrorResponse(w, http.StatusInternalServerError, "Internal server error", map[string]string{
				"error_code": "INTERNAL_ERROR",
			})
		}
		return
	}

	log.Info().
		Str("borrower_id", user.UserID).
		Str("loan_id", response.ID.String()).
		Float64("principal_amount", response.PrincipalAmount).
		Msg("Loan proposal created successfully")

	c.sendSuccessResponse(w, http.StatusCreated, "Loan proposal created successfully", response)
}

func (c *LoanController) ApproveLoan(w http.ResponseWriter, r *http.Request) {
	loanID := chi.URLParam(r, "id")
	if loanID == "" {
		c.sendErrorResponse(w, http.StatusBadRequest, "Loan ID is required", nil)
		return
	}

	user, err := middleware.GetUserFromCtx(r.Context())
	if err != nil {
		c.sendErrorResponse(w, http.StatusUnauthorized, "User context not found", nil)
		return
	}

	// Parse request body
	var req models.ApproveLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		c.sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	response, err := c.loanUsecase.ApproveLoan(r.Context(), loanID, user.UserID, &req)
	if err != nil {
		log.Error().Err(err).Str("loan_id", loanID).Str("employee_id", user.UserID).Msg("Failed to approve loan")

		// Handle specific errors
		errMsg := err.Error()
		switch {
		case errMsg == "loan not found":
			c.sendErrorResponse(w, http.StatusNotFound, "Loan not found", nil)
		case errMsg == "survey not completed":
			c.sendErrorResponse(w, http.StatusConflict, "Survey must be completed before approval", nil)
		case errMsg == "loan must be in proposed state":
			c.sendErrorResponse(w, http.StatusConflict, "Loan must be in proposed state", nil)
		case errMsg == "invalid loan ID" || errMsg == "invalid employee ID":
			c.sendErrorResponse(w, http.StatusBadRequest, errMsg, nil)
		default:
			c.sendErrorResponse(w, http.StatusInternalServerError, "Failed to approve loan", nil)
		}
		return
	}

	log.Info().
		Str("loan_id", loanID).
		Str("employee_id", user.UserID).
		Str("new_state", response.CurrentState).
		Msg("Loan approved successfully")

	c.sendSuccessResponse(w, http.StatusOK, "Loan approved successfully", response)
}

func (c *LoanController) sendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.Response[interface{}]{
		Data: map[string]interface{}{
			"success": true,
			"message": message,
			"data":    data,
		},
	}

	json.NewEncoder(w).Encode(response)
}

func (c *LoanController) sendErrorResponse(w http.ResponseWriter, statusCode int, message string, extra map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorData := map[string]interface{}{
		"success": false,
		"message": message,
	}

	for k, v := range extra {
		errorData[k] = v
	}

	response := models.Response[interface{}]{
		Data: errorData,
	}

	json.NewEncoder(w).Encode(response)
}

func (c *LoanController) sendValidationErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	var errors []map[string]string
	for _, err := range err.(validator.ValidationErrors) {
		fieldError := map[string]string{
			"field":   err.Field(),
			"message": helpers.GetValidationMessage(err),
		}
		errors = append(errors, fieldError)
	}

	response := models.Response[interface{}]{
		Data: map[string]interface{}{
			"success": false,
			"message": "Validation error",
			"errors":  errors,
		},
	}

	json.NewEncoder(w).Encode(response)
}

func (c *LoanController) DisburseLoan(w http.ResponseWriter, r *http.Request) {
	loanID := chi.URLParam(r, "id")
	if loanID == "" {
		c.sendErrorResponse(w, http.StatusBadRequest, "Loan ID is required", nil)
		return
	}

	user, err := middleware.GetUserFromCtx(r.Context())
	if err != nil {
		c.sendErrorResponse(w, http.StatusUnauthorized, "User context not found", nil)
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		c.sendErrorResponse(w, http.StatusBadRequest, "Invalid form data", nil)
		return
	}

	// Get signed agreement file
	file, header, err := r.FormFile("signed_agreement")
	if err != nil {
		c.sendErrorResponse(w, http.StatusBadRequest, "Signed agreement file is required", nil)
		return
	}
	defer file.Close()

	// Validate file type
	if !isValidSignedAgreementFile(header.Filename) {
		c.sendErrorResponse(w, http.StatusBadRequest, "Invalid file type, allowed: .pdf, .jpg, .jpeg", nil)
		return
	}

	// Validate file size
	if header.Size > 10*1024*1024 {
		c.sendErrorResponse(w, http.StatusBadRequest, "File size exceeds 10MB limit", nil)
		return
	}

	// Parse form data
	req := &models.DisburseLoanRequest{
		DisbursementNotes: r.FormValue("disbursement_notes"),
	}

	// Save signed agreement file
	signedAgreementURL, err := saveSignedAgreement(file, header, loanID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save signed agreement")
		c.sendErrorResponse(w, http.StatusInternalServerError, "Failed to save signed agreement", nil)
		return
	}

	// Disburse loan
	response, err := c.loanUsecase.DisburseLoan(r.Context(), loanID, user.UserID, req, signedAgreementURL)
	if err != nil {
		log.Error().Err(err).Str("loan_id", loanID).Str("officer_id", user.UserID).Msg("Failed to disburse loan")

		// Handle specific errors
		errMsg := err.Error()
		switch {
		case errMsg == "loan not found":
			c.sendErrorResponse(w, http.StatusNotFound, "Loan not found", nil)
		case errMsg == "loan must be in invested state":
			c.sendErrorResponse(w, http.StatusConflict, "Loan must be in invested state", nil)
		case errMsg == "invalid loan ID" || errMsg == "invalid officer ID":
			c.sendErrorResponse(w, http.StatusBadRequest, errMsg, nil)
		default:
			c.sendErrorResponse(w, http.StatusInternalServerError, "Failed to disburse loan", nil)
		}
		return
	}

	log.Info().
		Str("loan_id", loanID).
		Str("officer_id", user.UserID).
		Str("new_state", response.CurrentState).
		Msg("Loan disbursed successfully")

	c.sendSuccessResponse(w, http.StatusOK, "Loan disbursed successfully", response)
}

func isValidSignedAgreementFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validTypes := []string{".pdf", ".jpg", ".jpeg"}

	for _, validType := range validTypes {
		if ext == validType {
			return true
		}
	}
	return false
}

func saveSignedAgreement(file multipart.File, header *multipart.FileHeader, loanID string) (string, error) {
	// Create directory
	agreementDir := "uploads/agreements"
	if err := os.MkdirAll(agreementDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate filename
	fileExt := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("signed_agreement_%s%s", loanID, fileExt)
	filePath := filepath.Join(agreementDir, fileName)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return URL
	return fmt.Sprintf("/uploads/agreements/%s", fileName), nil
}
