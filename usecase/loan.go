package usecase

import (
	"context"
	"fmt"
	"github.com/fajar-andriansyah/loan-engine/internal/constants"
	"github.com/fajar-andriansyah/loan-engine/internal/pdf"
	"time"

	"github.com/fajar-andriansyah/loan-engine/models"
	"github.com/fajar-andriansyah/loan-engine/repositories"
	"github.com/google/uuid"
)

type LoanUsecase interface {
	CreateLoanProposal(ctx context.Context, req *models.CreateLoanRequest, borrowerID string) (*models.LoanResponse, error)
	ApproveLoan(ctx context.Context, loanID string, approvingEmployeeID string, req *models.ApproveLoanRequest) (*models.ApproveLoanResponse, error)
	DisburseLoan(ctx context.Context, loanID string, fieldOfficerID string, req *models.DisburseLoanRequest, signedAgreementURL string) (*models.DisburseLoanResponse, error)
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
		CurrentState:    constants.PROPOSED,
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

func (u *loanUsecase) ApproveLoan(ctx context.Context, loanID string, approvingEmployeeID string, req *models.ApproveLoanRequest) (*models.ApproveLoanResponse, error) {
	loanUUID, err := uuid.Parse(loanID)
	if err != nil {
		return nil, fmt.Errorf("invalid loan ID")
	}

	employeeUUID, err := uuid.Parse(approvingEmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee ID")
	}

	loan, err := u.loanRepo.GetLoanForApproval(ctx, loanUUID)
	if err != nil {
		return nil, err // Repository already handles "loan not found" and "survey not completed"
	}

	if loan.CurrentState != constants.PROPOSED {
		return nil, fmt.Errorf("loan must be in proposed state")
	}

	// Generate loan agreement PDF
	agreementURL, err := pdf.GenerateLoanAgreement(loan)
	if err != nil {
		return nil, fmt.Errorf("failed to generate agreement: %w", err)
	}

	err = u.loanRepo.ApproveLoan(ctx, loanUUID, employeeUUID, req.ApprovalNotes, agreementURL)
	if err != nil {
		return nil, err
	}

	response, err := u.loanRepo.GetApprovedLoan(ctx, loanUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get approved loan data: %w", err)
	}

	return response, nil
}

func (u *loanUsecase) DisburseLoan(ctx context.Context, loanID string, fieldOfficerID string, req *models.DisburseLoanRequest, signedAgreementURL string) (*models.DisburseLoanResponse, error) {
	loanUUID, err := uuid.Parse(loanID)
	if err != nil {
		return nil, fmt.Errorf("invalid loan ID")
	}

	officerUUID, err := uuid.Parse(fieldOfficerID)
	if err != nil {
		return nil, fmt.Errorf("invalid officer ID")
	}

	loan, err := u.loanRepo.GetLoanForDisbursement(ctx, loanUUID)
	if err != nil {
		return nil, err
	}

	if loan.CurrentState != constants.INVESTED {
		return nil, fmt.Errorf("loan must be in invested state")
	}

	err = u.loanRepo.DisburseLoan(ctx, loanUUID, officerUUID, signedAgreementURL, req.DisbursementNotes)
	if err != nil {
		return nil, err
	}

	response, err := u.loanRepo.GetDisbursedLoan(ctx, loanUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get disbursed loan data: %w", err)
	}

	return response, nil
}
