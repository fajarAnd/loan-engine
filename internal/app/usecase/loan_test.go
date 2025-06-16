package usecase

import (
	"context"
	"fmt"
	mocksPdf "github.com/fajar-andriansyah/loan-engine/internal/app/mocks/pdf"
	mocksRepo "github.com/fajar-andriansyah/loan-engine/internal/app/mocks/repositories"
	"github.com/fajar-andriansyah/loan-engine/internal/app/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateLoanProposal_InitialStateIsProposed(t *testing.T) {
	mockRepo := mocksRepo.NewLoanRepository(t)
	mockPdfGen := mocksPdf.NewPDFGenerator(t)
	loanUsecase := NewLoanUsecase(mockRepo, mockPdfGen)

	borrowerID := uuid.New()
	req := &models.CreateLoanRequest{
		PrincipalAmount: 5000000,
		InterestRate:    10,
		ROIRate:         8,
		LoanTermMonth:   12,
	}

	mockRepo.On("CreateLoan", mock.Anything, mock.MatchedBy(func(loan *models.Loan) bool {
		return loan.CurrentState == "PROPOSED" &&
			loan.BorrowerID == borrowerID &&
			loan.PrincipalAmount == 5000000
	})).Return(nil)

	result, err := loanUsecase.CreateLoanProposal(context.Background(), req, borrowerID.String())

	assert.NoError(t, err)
	assert.Equal(t, "PROPOSED", result.CurrentState)
	assert.Equal(t, borrowerID, result.BorrowerID)
	assert.Equal(t, float64(5000000), result.PrincipalAmount)
}

func TestApproveLoan_RequiresSurveyCompletion(t *testing.T) {
	mockRepo := mocksRepo.NewLoanRepository(t)
	mockPdfGen := mocksPdf.NewPDFGenerator(t)
	loanUsecase := NewLoanUsecase(mockRepo, mockPdfGen)

	loanID := uuid.New()
	employeeID := uuid.New()

	req := &models.ApproveLoanRequest{
		ApprovalNotes: "Approval attempt",
	}

	mockRepo.On("GetLoanForApproval", mock.Anything, loanID).Return(nil, fmt.Errorf("survey not completed"))

	result, err := loanUsecase.ApproveLoan(context.Background(), loanID.String(), employeeID.String(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "survey not completed")
}

func TestApproveLoan_PreventInvalidStateTransition(t *testing.T) {
	mockRepo := mocksRepo.NewLoanRepository(t)
	mockPdfGen := mocksPdf.NewPDFGenerator(t)
	loanUsecase := NewLoanUsecase(mockRepo, mockPdfGen)

	loanID := uuid.New()
	employeeID := uuid.New()

	req := &models.ApproveLoanRequest{
		ApprovalNotes: "Approval",
	}

	loanForApproval := &models.LoanForApproval{
		ID:           loanID,
		CurrentState: "APPROVED",
		SurveyDate:   time.Now(),
	}

	mockRepo.On("GetLoanForApproval", mock.Anything, loanID).Return(loanForApproval, nil)

	result, err := loanUsecase.ApproveLoan(context.Background(), loanID.String(), employeeID.String(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "loan must be in proposed state")
}

func TestApproveLoan_InvalidLoanID(t *testing.T) {
	mockRepo := mocksRepo.NewLoanRepository(t)
	mockPdfGen := mocksPdf.NewPDFGenerator(t)
	loanUsecase := NewLoanUsecase(mockRepo, mockPdfGen)

	employeeID := uuid.New()

	req := &models.ApproveLoanRequest{
		ApprovalNotes: "Approval attempt",
	}

	result, err := loanUsecase.ApproveLoan(context.Background(), "invalid-uuid", employeeID.String(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid loan ID")
}

func TestApproveLoan_SuccessfulApprovalWithPDFGeneration(t *testing.T) {
	mockRepo := mocksRepo.NewLoanRepository(t)
	mockPdfGen := mocksPdf.NewPDFGenerator(t)
	loanUsecase := NewLoanUsecase(mockRepo, mockPdfGen)

	loanID := uuid.New()
	employeeID := uuid.New()
	borrowerID := uuid.New()

	req := &models.ApproveLoanRequest{
		ApprovalNotes: "Loan approved after field verification",
	}

	loanForApproval := &models.LoanForApproval{
		ID:                       loanID,
		BorrowerID:               borrowerID,
		BorrowerName:             "Test Borrower",
		PrincipalAmount:          5000000,
		InterestRate:             10,
		ROIRate:                  8,
		LoanTermMonth:            12,
		CurrentState:             "PROPOSED",
		FieldValidatorEmployeeID: uuid.New(),
		SurveyDate:               time.Now(),
	}

	expectedAgreementURL := "/uploads/agreements/loan_agreement_" + loanID.String() + ".pdf"

	approvedLoanResponse := &models.ApproveLoanResponse{
		ID:                       loanID,
		BorrowerID:               borrowerID,
		PrincipalAmount:          5000000,
		InterestRate:             10,
		ROIRate:                  8,
		LoanTermMonth:            12,
		CurrentState:             "APPROVED",
		ApprovalDate:             time.Now().Format("2006-01-02"),
		ApprovingEmployeeID:      employeeID,
		ApprovalNotes:            req.ApprovalNotes,
		LoanAgreementPDFURL:      expectedAgreementURL,
		FieldValidatorEmployeeID: loanForApproval.FieldValidatorEmployeeID,
		SurveyDate:               loanForApproval.SurveyDate.Format("2006-01-02"),
		UpdatedAt:                time.Now(),
	}

	mockRepo.On("GetLoanForApproval", mock.Anything, loanID).Return(loanForApproval, nil)
	mockPdfGen.On("GenerateLoanAgreement", loanForApproval).Return(expectedAgreementURL, nil)
	mockRepo.On("ApproveLoan", mock.Anything, loanID, employeeID, req.ApprovalNotes, expectedAgreementURL).Return(nil)
	mockRepo.On("GetApprovedLoan", mock.Anything, loanID).Return(approvedLoanResponse, nil)
	
	result, err := loanUsecase.ApproveLoan(context.Background(), loanID.String(), employeeID.String(), req)

	assert.NoError(t, err)
	assert.Equal(t, "APPROVED", result.CurrentState)
	assert.Equal(t, employeeID, result.ApprovingEmployeeID)
	assert.Equal(t, expectedAgreementURL, result.LoanAgreementPDFURL)
	assert.Equal(t, req.ApprovalNotes, result.ApprovalNotes)
}

func TestDisburseLoan_RequiresInvestedState(t *testing.T) {
	mockRepo := mocksRepo.NewLoanRepository(t)
	mockPdfGen := mocksPdf.NewPDFGenerator(t)
	loanUsecase := NewLoanUsecase(mockRepo, mockPdfGen)

	loanID := uuid.New()
	officerID := uuid.New()

	req := &models.DisburseLoanRequest{
		DisbursementNotes: "Money disbursed to borrower",
	}

	loan := &models.Loan{
		ID:           loanID,
		CurrentState: "APPROVED",
	}

	mockRepo.On("GetLoanForDisbursement", mock.Anything, loanID).Return(loan, nil)

	result, err := loanUsecase.DisburseLoan(context.Background(), loanID.String(), officerID.String(), req, "/signed-agreement.pdf")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "loan must be in invested state")
}

// Valid Disbursement Flow
func TestDisburseLoan_SuccessfulStateTransition(t *testing.T) {
	mockRepo := mocksRepo.NewLoanRepository(t)
	mockPdfGen := mocksPdf.NewPDFGenerator(t)
	loanUsecase := NewLoanUsecase(mockRepo, mockPdfGen)

	loanID := uuid.New()
	officerID := uuid.New()
	signedAgreementURL := "/uploads/agreements/signed_agreement_" + loanID.String() + ".pdf"

	req := &models.DisburseLoanRequest{
		DisbursementNotes: "Money disbursed successfully",
	}

	loan := &models.Loan{
		ID:           loanID,
		CurrentState: "INVESTED",
	}

	disbursedLoan := &models.DisburseLoanResponse{
		ID:                     loanID,
		CurrentState:           "DISBURSED",
		FieldOfficerEmployeeID: officerID,
		SignedAgreementURL:     signedAgreementURL,
		DisbursementNotes:      req.DisbursementNotes,
	}

	mockRepo.On("GetLoanForDisbursement", mock.Anything, loanID).Return(loan, nil)
	mockRepo.On("DisburseLoan", mock.Anything, loanID, officerID, signedAgreementURL, req.DisbursementNotes).Return(nil)
	mockRepo.On("GetDisbursedLoan", mock.Anything, loanID).Return(disbursedLoan, nil)

	result, err := loanUsecase.DisburseLoan(context.Background(), loanID.String(), officerID.String(), req, signedAgreementURL)

	assert.NoError(t, err)
	assert.Equal(t, "DISBURSED", result.CurrentState)
	assert.Equal(t, officerID, result.FieldOfficerEmployeeID)
	assert.Equal(t, signedAgreementURL, result.SignedAgreementURL)
}

func TestDisburseLoan_InvalidEmployeeID(t *testing.T) {
	// Arrange
	mockRepo := mocksRepo.NewLoanRepository(t)
	mockPdfGen := mocksPdf.NewPDFGenerator(t)
	loanUsecase := NewLoanUsecase(mockRepo, mockPdfGen)

	loanID := uuid.New()

	req := &models.DisburseLoanRequest{
		DisbursementNotes: "Money disbursed",
	}

	result, err := loanUsecase.DisburseLoan(context.Background(), loanID.String(), "invalid-officer-id", req, "/agreement.pdf")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid officer ID")
}

func TestApproveLoan_InvalidEmployeeID(t *testing.T) {
	mockRepo := mocksRepo.NewLoanRepository(t)
	mockPdfGen := mocksPdf.NewPDFGenerator(t)
	loanUsecase := NewLoanUsecase(mockRepo, mockPdfGen)

	loanID := uuid.New()

	req := &models.ApproveLoanRequest{
		ApprovalNotes: "Approval attempt",
	}

	result, err := loanUsecase.ApproveLoan(context.Background(), loanID.String(), "invalid-employee-uuid", req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid employee ID")
}
