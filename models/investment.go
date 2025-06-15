package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateInvestmentRequest struct {
	InvestmentAmount float64 `json:"investment_amount" validate:"required,gt=0"`
}

type InvestmentResponse struct {
	ID                  uuid.UUID `json:"id"`
	LoanID              uuid.UUID `json:"loan_id"`
	InvestorID          uuid.UUID `json:"investor_id"`
	InvestmentAmount    float64   `json:"investment_amount"`
	ExpectedReturn      float64   `json:"expected_return"`
	InvestmentDate      string    `json:"investment_date"`
	LoanCurrentState    string    `json:"loan_current_state"`
	TotalInvestedAmount float64   `json:"total_invested_amount"`
	RemainingAmount     float64   `json:"remaining_amount"`
	AgreementURL        string    `json:"agreement_url,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
}

type Investment struct {
	ID               uuid.UUID `json:"id"`
	LoanID           uuid.UUID `json:"loan_id"`
	InvestorID       uuid.UUID `json:"investor_id"`
	InvestmentAmount float64   `json:"investment_amount"`
	ExpectedReturn   float64   `json:"expected_return"`
	InvestmentDate   time.Time `json:"investment_date"`
	CreatedAt        time.Time `json:"created_at"`
}

type LoanInvestmentInfo struct {
	ID              uuid.UUID `json:"id"`
	PrincipalAmount float64   `json:"principal_amount"`
	ROIRate         float64   `json:"roi_rate"`
	CurrentState    string    `json:"current_state"`
	TotalInvested   float64   `json:"total_invested"`
}
