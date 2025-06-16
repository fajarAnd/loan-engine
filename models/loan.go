package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateLoanRequest struct {
	PrincipalAmount float64 `json:"principal_amount" validate:"required,gt=0"`
	InterestRate    float64 `json:"interest_rate" validate:"required,gte=0"`
	ROIRate         float64 `json:"roi_rate" validate:"required,gte=0"`
	LoanTermMonth   int     `json:"loan_term_month" validate:"required,gt=0"`
}

type LoanResponse struct {
	ID              uuid.UUID `json:"id"`
	BorrowerID      uuid.UUID `json:"borrower_id"`
	PrincipalAmount float64   `json:"principal_amount"`
	InterestRate    float64   `json:"interest_rate"`
	ROIRate         float64   `json:"roi_rate"`
	LoanTermMonth   int       `json:"loan_term_month"`
	CurrentState    string    `json:"current_state"`
	CreatedAt       time.Time `json:"created_at"`
}

type Loan struct {
	ID              uuid.UUID `json:"id"`
	BorrowerID      uuid.UUID `json:"borrower_id"`
	PrincipalAmount float64   `json:"principal_amount"`
	InterestRate    float64   `json:"interest_rate"`
	ROIRate         float64   `json:"roi_rate"`
	LoanTermMonth   int       `json:"loan_term_month"`
	CurrentState    string    `json:"current_state"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type LoanForApproval struct {
	ID                       uuid.UUID `json:"id"`
	BorrowerID               uuid.UUID `json:"borrower_id"`
	BorrowerName             string    `json:"borrower_name"`
	PrincipalAmount          float64   `json:"principal_amount"`
	InterestRate             float64   `json:"interest_rate"`
	ROIRate                  float64   `json:"roi_rate"`
	LoanTermMonth            int       `json:"loan_term_month"`
	CurrentState             string    `json:"current_state"`
	FieldValidatorEmployeeID uuid.UUID `json:"field_validator_employee_id"`
	SurveyDate               time.Time `json:"survey_date"`
}

type ApproveLoanRequest struct {
	ApprovalNotes string `json:"approval_notes"`
}

type ApproveLoanResponse struct {
	ID                       uuid.UUID `json:"id"`
	BorrowerID               uuid.UUID `json:"borrower_id"`
	PrincipalAmount          float64   `json:"principal_amount"`
	InterestRate             float64   `json:"interest_rate"`
	ROIRate                  float64   `json:"roi_rate"`
	LoanTermMonth            int       `json:"loan_term_month"`
	CurrentState             string    `json:"current_state"`
	ApprovalDate             string    `json:"approval_date"`
	ApprovingEmployeeID      uuid.UUID `json:"approving_employee_id"`
	ApprovalNotes            string    `json:"approval_notes,omitempty"`
	LoanAgreementPDFURL      string    `json:"loan_agreement_pdf_url"`
	FieldValidatorEmployeeID uuid.UUID `json:"field_validator_employee_id"`
	SurveyDate               string    `json:"survey_date"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type DisburseLoanRequest struct {
	DisbursementNotes string `form:"disbursement_notes"`
}

type DisburseLoanResponse struct {
	ID                     uuid.UUID `json:"id"`
	BorrowerID             uuid.UUID `json:"borrower_id"`
	PrincipalAmount        float64   `json:"principal_amount"`
	InterestRate           float64   `json:"interest_rate"`
	ROIRate                float64   `json:"roi_rate"`
	LoanTermMonth          int       `json:"loan_term_month"`
	CurrentState           string    `json:"current_state"`
	DisbursementDate       string    `json:"disbursement_date"`
	FieldOfficerEmployeeID uuid.UUID `json:"field_officer_employee_id"`
	SignedAgreementURL     string    `json:"signed_agreement_url"`
	DisbursementNotes      string    `json:"disbursement_notes,omitempty"`
	UpdatedAt              time.Time `json:"updated_at"`
}
