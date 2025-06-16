package controller

import (
	"encoding/json"
	"fmt"
	models2 "github.com/fajar-andriansyah/loan-engine/internal/app/models"
	"github.com/fajar-andriansyah/loan-engine/internal/app/usecase"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type AuthController struct {
	authUsecase usecase.AuthUsecase
	validator   *validator.Validate
}

func NewAuthController(authUsecase usecase.AuthUsecase) *AuthController {
	return &AuthController{
		authUsecase: authUsecase,
		validator:   validator.New(),
	}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req models2.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		c.sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := c.validator.Struct(&req); err != nil {
		log.Error().Err(err).Msg("Validation failed")
		c.sendValidationErrorResponse(w, err)
		return
	}

	resp, err := c.authUsecase.Login(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Str("user_type", req.UserType).Msg("Login failed")
		
		if err.Error() == "invalid credentials" || err.Error() == "authentication failed" {
			c.sendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials", map[string]string{
				"error_code": "INVALID_CREDENTIALS",
			})
		} else {
			c.sendErrorResponse(w, http.StatusInternalServerError, "Internal server error", map[string]string{
				"error_code": "INTERNAL_ERROR",
			})
		}
		return
	}

	log.Info().Str("email", req.Email).Str("user_type", req.UserType).Msg("Login successful")
	c.sendSuccessResponse(w, http.StatusOK, "Login successful", resp)
}

func (c *AuthController) sendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := models2.Response[interface{}]{
		Data: map[string]interface{}{
			"success": true,
			"message": message,
			"data":    data,
		},
	}

	json.NewEncoder(w).Encode(resp)
}

func (c *AuthController) sendErrorResponse(w http.ResponseWriter, statusCode int, message string, extra map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorData := map[string]interface{}{
		"success": false,
		"message": message,
	}

	for k, v := range extra {
		errorData[k] = v
	}

	response := models2.Response[interface{}]{
		Data: errorData,
	}

	json.NewEncoder(w).Encode(response)
}

func (c *AuthController) sendValidationErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	var errors []map[string]string
	for _, err := range err.(validator.ValidationErrors) {
		fieldError := map[string]string{
			"field":   err.Field(),
			"message": getValidationMessage(err),
		}
		errors = append(errors, fieldError)
	}

	response := models2.Response[interface{}]{
		Data: map[string]interface{}{
			"success": false,
			"message": "Validation error",
			"errors":  errors,
		},
	}

	json.NewEncoder(w).Encode(response)
}

func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}
