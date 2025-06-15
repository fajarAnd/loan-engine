package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/fajar-andriansyah/loan-engine/internal/pdf"
	"github.com/fajar-andriansyah/loan-engine/models"
	"github.com/fajar-andriansyah/loan-engine/repositories"
	"github.com/google/uuid"
)

type InvestmentUsecase interface {
	CreateInvestment(ctx context.Context, loanID, investorID string, req *models.CreateInvestmentRequest) (*models.InvestmentResponse, error)
}

type investmentUsecase struct {
	investmentRepo repositories.InvestmentRepository
}

func NewInvestmentUsecase(investmentRepo repositories.InvestmentRepository) InvestmentUsecase {
	return &investmentUsecase{
		investmentRepo: investmentRepo,
	}
}

func (u *investmentUsecase) CreateInvestment(ctx context.Context, loanID, investorID string, req *models.CreateInvestmentRequest) (*models.InvestmentResponse, error) {
	// Parse UUIDs
	loanUUID, err := uuid.Parse(loanID)
	if err != nil {
		return nil, fmt.Errorf("invalid loan ID")
	}

	investorUUID, err := uuid.Parse(investorID)
	if err != nil {
		return nil, fmt.Errorf("invalid investor ID")
	}

	// Get loan information
	loan, err := u.investmentRepo.GetLoanForInvestment(ctx, loanUUID)
	if err != nil {
		return nil, err
	}

	// Check loan state (must be APPROVED or FUNDING)
	if loan.CurrentState != "APPROVED" && loan.CurrentState != "FUNDING" {
		return nil, fmt.Errorf("loan must be in APPROVED or FUNDING state")
	}

	// Check if investor already invested in this loan
	exists, err := u.investmentRepo.CheckExistingInvestment(ctx, loanUUID, investorUUID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("investor has already invested in this loan")
	}

	// Calculate remaining amount
	remainingAmount := loan.PrincipalAmount - loan.TotalInvested

	// Check if investment amount exceeds remaining
	if req.InvestmentAmount > remainingAmount {
		return nil, fmt.Errorf("investment amount exceeds remaining loan amount")
	}

	// Calculate expected return based on ROI rate
	expectedReturn := req.InvestmentAmount * (loan.ROIRate / 100)

	// Get investor name for agreement
	investorName, err := u.investmentRepo.GetInvestorName(ctx, investorUUID)
	if err != nil {
		return nil, err
	}

	// Create investment
	now := time.Now()
	investment := &models.Investment{
		ID:               uuid.New(),
		LoanID:           loanUUID,
		InvestorID:       investorUUID,
		InvestmentAmount: req.InvestmentAmount,
		ExpectedReturn:   expectedReturn,
		InvestmentDate:   now,
		CreatedAt:        now,
	}

	err = u.investmentRepo.CreateInvestment(ctx, investment)
	if err != nil {
		return nil, err
	}

	// Generate individual investment agreement PDF
	agreementURL, err := pdf.GenerateInvestmentAgreement(investment, loan, investorName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate investment agreement: %w", err)
	}

	// Calculate new totals
	newTotalInvested := loan.TotalInvested + req.InvestmentAmount
	newRemainingAmount := loan.PrincipalAmount - newTotalInvested

	// Update loan state based on funding status
	var newState string

	if newRemainingAmount == 0 {
		// Fully funded - transition to INVESTED
		newState = "INVESTED"
	} else if loan.CurrentState == "APPROVED" {
		// First investment - transition to FUNDING
		newState = "FUNDING"
	} else {
		// Already in FUNDING, stay in FUNDING
		newState = "FUNDING"
	}

	// Update loan state if changed
	if newState != loan.CurrentState {
		err = u.investmentRepo.UpdateLoanState(ctx, loanUUID, newState)
		if err != nil {
			return nil, err
		}
	}

	// Prepare response
	response := &models.InvestmentResponse{
		ID:                  investment.ID,
		LoanID:              investment.LoanID,
		InvestorID:          investment.InvestorID,
		InvestmentAmount:    investment.InvestmentAmount,
		ExpectedReturn:      investment.ExpectedReturn,
		InvestmentDate:      investment.InvestmentDate.Format("2006-01-02"),
		LoanCurrentState:    newState,
		TotalInvestedAmount: newTotalInvested,
		RemainingAmount:     newRemainingAmount,
		AgreementURL:        agreementURL,
		CreatedAt:           investment.CreatedAt,
	}

	return response, nil
}
