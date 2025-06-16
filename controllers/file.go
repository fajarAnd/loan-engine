package controller

import (
	"encoding/json"
	"github.com/fajar-andriansyah/loan-engine/infrastructure/http/middleware"
	"github.com/fajar-andriansyah/loan-engine/internal/helpers"
	"github.com/fajar-andriansyah/loan-engine/models"
	"github.com/fajar-andriansyah/loan-engine/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"net/http"
)

type FileController struct {
	fileUsecase usecase.FileUsecase
	validator   *validator.Validate
}

func NewFileController(fileUsecase usecase.FileUsecase) *FileController {
	return &FileController{
		fileUsecase: fileUsecase,
		validator:   validator.New(),
	}
}

func (c *FileController) UploadSurveyDocument(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromCtx(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user from context")
		c.sendErrorResponse(w, http.StatusUnauthorized, "User context not found", nil)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get file from form")
		c.sendErrorResponse(w, http.StatusBadRequest, "File is required", map[string]string{
			"error_code": "FILE_REQUIRED",
		})
		return
	}
	defer file.Close()

	// Check file size in header
	if header.Size > 10*1024*1024 {
		c.sendErrorResponse(w, http.StatusBadRequest, "File size exceeds maximum limit of 10MB", map[string]string{
			"error_code": "FILE_TOO_LARGE",
		})
		return
	}

	// Parse form data
	req := &models.UploadDocumentRequest{
		LoanID:      r.FormValue("loan_id"),
		SurveyDate:  r.FormValue("survey_date"),
		SurveyNotes: r.FormValue("survey_notes"),
	}

	// Validate request
	if err := c.validator.Struct(req); err != nil {
		log.Error().Err(err).Msg("Validation failed")
		c.sendValidationErrorResponse(w, err)
		return
	}

	// Upload document
	response, err := c.fileUsecase.UploadSurveyDocument(r.Context(), req, file, header, user.UserID)
	if err != nil {
		log.Error().Err(err).Str("loan_id", req.LoanID).Msg("Failed to upload survey document")

		// Handle specific errors
		errMsg := err.Error()
		switch {
		case contains(errMsg, "loan not found"):
			c.sendErrorResponse(w, http.StatusNotFound, "Loan not found", map[string]string{
				"error_code": "LOAN_NOT_FOUND",
			})
		case contains(errMsg, "not in PROPOSED state"):
			c.sendErrorResponse(w, http.StatusConflict, errMsg, map[string]string{
				"error_code": "INVALID_LOAN_STATE",
			})
		case contains(errMsg, "invalid file type"):
			c.sendErrorResponse(w, http.StatusBadRequest, errMsg, map[string]string{
				"error_code": "INVALID_FILE_TYPE",
			})
		case contains(errMsg, "file size exceeds"):
			c.sendErrorResponse(w, http.StatusBadRequest, errMsg, map[string]string{
				"error_code": "FILE_TOO_LARGE",
			})
		case contains(errMsg, "invalid survey date"):
			c.sendErrorResponse(w, http.StatusBadRequest, "Invalid survey date format, expected YYYY-MM-DD", map[string]string{
				"error_code": "INVALID_DATE_FORMAT",
			})
		default:
			c.sendErrorResponse(w, http.StatusInternalServerError, "Failed to upload document", map[string]string{
				"error_code": "UPLOAD_FAILED",
			})
		}
		return
	}

	log.Info().
		Str("loan_id", req.LoanID).
		Str("validator_id", user.UserID).
		Str("file_name", response.FileName).
		Msg("Survey document uploaded successfully")

	c.sendSuccessResponse(w, http.StatusCreated, "Survey document uploaded successfully", response)
}

func (c *FileController) sendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
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

func (c *FileController) sendErrorResponse(w http.ResponseWriter, statusCode int, message string, extra map[string]string) {
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

func (c *FileController) sendValidationErrorResponse(w http.ResponseWriter, err error) {
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

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && s[len(s)-len(substr):] == substr ||
		(len(s) > len(substr) && func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
