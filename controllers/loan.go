package controller

import (
	"encoding/json"
	"github.com/fajar-andriansyah/loan-engine/internal/helpers"
	"net/http"

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

	user, err := middleware.GetUserFromContext(r.Context())
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
