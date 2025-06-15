package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fajar-andriansyah/loan-engine/infrastructure/http/middleware"
	"github.com/fajar-andriansyah/loan-engine/internal/helpers"
	"github.com/fajar-andriansyah/loan-engine/models"
	"github.com/fajar-andriansyah/loan-engine/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type InvestmentController struct {
	investmentUsecase usecase.InvestmentUsecase
	validator         *validator.Validate
}

func NewInvestmentController(investmentUsecase usecase.InvestmentUsecase) *InvestmentController {
	return &InvestmentController{
		investmentUsecase: investmentUsecase,
		validator:         validator.New(),
	}
}

func (c *InvestmentController) CreateInvestment(w http.ResponseWriter, r *http.Request) {
	// Get loan ID from URL parameter
	loanID := chi.URLParam(r, "id")
	if loanID == "" {
		c.sendErrorResponse(w, http.StatusBadRequest, "Loan ID is required", map[string]string{
			"error_code": "MISSING_LOAN_ID",
		})
		return
	}

	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user from context")
		c.sendErrorResponse(w, http.StatusUnauthorized, "User context not found", nil)
		return
	}

	// Parse request body
	var req models.CreateInvestmentRequest
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

	response, err := c.investmentUsecase.CreateInvestment(r.Context(), loanID, user.UserID, &req)
	if err != nil {
		log.Error().Err(err).
			Str("loan_id", loanID).
			Str("investor_id", user.UserID).
			Float64("investment_amount", req.InvestmentAmount).
			Msg("Failed to create investment")

		c.handleInvestmentError(w, err)
		return
	}

	// Determine response message
	message := "Investment created successfully"
	if response.LoanCurrentState == "INVESTED" {
		message = "Investment created successfully. Loan is now fully invested!"
	}

	log.Info().
		Str("loan_id", loanID).
		Str("investor_id", user.UserID).
		Str("investment_id", response.ID.String()).
		Float64("investment_amount", req.InvestmentAmount).
		Str("loan_state", response.LoanCurrentState).
		Msg("Investment created successfully")

	c.sendSuccessResponse(w, http.StatusCreated, message, response)
}

func (c *InvestmentController) handleInvestmentError(w http.ResponseWriter, err error) {
	errMsg := err.Error()

	switch {
	case errMsg == "loan not found":
		c.sendErrorResponse(w, http.StatusNotFound, "Loan not found", map[string]string{
			"error_code": "LOAN_NOT_FOUND",
		})
	case errMsg == "invalid loan ID" || errMsg == "invalid investor ID":
		c.sendErrorResponse(w, http.StatusBadRequest, errMsg, map[string]string{
			"error_code": "INVALID_ID",
		})
	case strings.Contains(errMsg, "loan must be in APPROVED or FUNDING state"):
		c.sendErrorResponse(w, http.StatusConflict, errMsg, map[string]string{
			"error_code": "INVALID_LOAN_STATE",
		})
	case errMsg == "investor has already invested in this loan":
		c.sendErrorResponse(w, http.StatusConflict, errMsg, map[string]string{
			"error_code": "DUPLICATE_INVESTMENT",
		})
	case strings.Contains(errMsg, "investment amount exceeds remaining"):
		c.sendErrorResponse(w, http.StatusConflict, errMsg, map[string]string{
			"error_code": "INVESTMENT_EXCEEDS_REMAINING",
		})
	default:
		c.sendErrorResponse(w, http.StatusInternalServerError, "Failed to create investment", map[string]string{
			"error_code": "INTERNAL_ERROR",
		})
	}
}

func (c *InvestmentController) sendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
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

func (c *InvestmentController) sendErrorResponse(w http.ResponseWriter, statusCode int, message string, extra map[string]string) {
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

func (c *InvestmentController) sendValidationErrorResponse(w http.ResponseWriter, err error) {
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
