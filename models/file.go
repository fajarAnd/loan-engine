package models

import (
	"github.com/google/uuid"
	"time"
)

type UploadDocumentRequest struct {
	LoanID      string `form:"loan_id" validate:"required,uuid"`
	SurveyDate  string `form:"survey_date" validate:"required"`
	SurveyNotes string `form:"survey_notes"`
}

type UploadDocumentResponse struct {
	LoanID                   uuid.UUID `json:"loan_id"`
	FileName                 string    `json:"file_name"`
	FileURL                  string    `json:"file_url"`
	FileType                 string    `json:"file_type"`
	FieldValidatorEmployeeID uuid.UUID `json:"field_validator_employee_id"`
	SurveyDate               string    `json:"survey_date"`
	SurveyNotes              string    `json:"survey_notes,omitempty"`
	UploadedAt               time.Time `json:"uploaded_at"`
}
