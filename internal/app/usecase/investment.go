package usecase

import (
	"context"
	"fmt"
	"github.com/fajar-andriansyah/loan-engine/internal/app/constants"
	"github.com/fajar-andriansyah/loan-engine/internal/app/models"
	"github.com/fajar-andriansyah/loan-engine/internal/app/repositories"
	"github.com/fajar-andriansyah/loan-engine/internal/pkg/pdf"
	"time"

	"github.com/google/uuid"
)

type InvestmentUsecase interface {
	CreateInvestment(ctx context.Context, loanID, investorID string, req *models.CreateInvestmentRequest) (*models.InvestmentResponse, error)
}

type investmentUsecase struct {
	investmentRepo repositories.InvestmentRepository
	pdfGenerator   pdf.PDFGenerator
}

func NewInvestmentUsecase(investmentRepo repositories.InvestmentRepository, pdfGenerator pdf.PDFGenerator) InvestmentUsecase {
	return &investmentUsecase{
		investmentRepo: investmentRepo,
		pdfGenerator:   pdfGenerator,
	}
}

func (u *investmentUsecase) CreateInvestment(ctx context.Context, loanID, investorID string, req *models.CreateInvestmentRequest) (*models.InvestmentResponse, error) {
	loanUUID, err := uuid.Parse(loanID)
	if err != nil {
		return nil, fmt.Errorf("invalid loan ID")
	}

	investorUUID, err := uuid.Parse(investorID)
	if err != nil {
		return nil, fmt.Errorf("invalid investor ID")
	}

	loan, err := u.investmentRepo.GetLoanForInvestment(ctx, loanUUID)
	if err != nil {
		return nil, err
	}

	if loan.CurrentState != constants.APPROVED && loan.CurrentState != constants.FUNDING {
		return nil, fmt.Errorf("loan must be in APPROVED or FUNDING state")
	}

	exists, err := u.investmentRepo.CheckExistingInvestment(ctx, loanUUID, investorUUID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("investor has already invested in this loan")
	}

	remainingAmount := loan.PrincipalAmount - loan.TotalInvested
	if req.InvestmentAmount > remainingAmount {
		return nil, fmt.Errorf("investment amount exceeds remaining loan amount")
	}

	expectedReturn := req.InvestmentAmount * (loan.ROIRate / 100)

	investorName, err := u.investmentRepo.GetInvestorName(ctx, investorUUID)
	if err != nil {
		return nil, err
	}

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
	agreementURL, err := u.pdfGenerator.GenerateInvestmentAgreement(investment, loan, investorName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate investment agreement: %w", err)
	}

	newTotalInvested := loan.TotalInvested + req.InvestmentAmount
	newRemainingAmount := loan.PrincipalAmount - newTotalInvested

	var newState string

	if newRemainingAmount == 0 {
		// Loan = principal - move to invested
		newState = constants.INVESTED
	} else if loan.CurrentState == constants.APPROVED {
		// First investment - transition to FUNDING
		newState = constants.FUNDING
	} else {
		newState = constants.FUNDING
	}

	if newState != loan.CurrentState {
		err = u.investmentRepo.UpdateLoanState(ctx, loanUUID, newState)
		if err != nil {
			return nil, err
		}
	}

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
