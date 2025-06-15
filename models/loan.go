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
