package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/fajar-andriansyah/loan-engine/models"
	"github.com/fajar-andriansyah/loan-engine/repositories"
	"github.com/google/uuid"
)

type LoanUsecase interface {
	CreateLoanProposal(ctx context.Context, req *models.CreateLoanRequest, borrowerID string) (*models.LoanResponse, error)
}

type loanUsecase struct {
	loanRepo repositories.LoanRepository
}

func NewLoanUsecase(loanRepo repositories.LoanRepository) LoanUsecase {
	return &loanUsecase{
		loanRepo: loanRepo,
	}
}

func (u *loanUsecase) CreateLoanProposal(ctx context.Context, req *models.CreateLoanRequest, borrowerID string) (*models.LoanResponse, error) {
	borrowerUUID, err := uuid.Parse(borrowerID)
	if err != nil {
		return nil, fmt.Errorf("invalid borrower ID: %w", err)
	}

	loanID := uuid.New()
	now := time.Now()

	loan := &models.Loan{
		ID:              loanID,
		BorrowerID:      borrowerUUID,
		PrincipalAmount: req.PrincipalAmount,
		InterestRate:    req.InterestRate,
		ROIRate:         req.ROIRate,
		LoanTermMonth:   req.LoanTermMonth,
		CurrentState:    "PROPOSED",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := u.loanRepo.CreateLoan(ctx, loan); err != nil {
		return nil, fmt.Errorf("failed to create loan: %w", err)
	}

	response := &models.LoanResponse{
		ID:              loan.ID,
		BorrowerID:      loan.BorrowerID,
		PrincipalAmount: loan.PrincipalAmount,
		InterestRate:    loan.InterestRate,
		ROIRate:         loan.ROIRate,
		LoanTermMonth:   loan.LoanTermMonth,
		CurrentState:    loan.CurrentState,
		CreatedAt:       loan.CreatedAt,
	}

	return response, nil
}
