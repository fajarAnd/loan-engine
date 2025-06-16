package usecase

import (
	"context"
	mocksPdf "github.com/fajar-andriansyah/loan-engine/internal/app/mocks/pdf"
	mocksRepo "github.com/fajar-andriansyah/loan-engine/internal/app/mocks/repositories"
	"github.com/fajar-andriansyah/loan-engine/internal/app/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// State Transition (APPROVED -> FUNDING)
func TestCreateInvestment_FirstInvestmentTransitionsToFunding(t *testing.T) {
	mockRepo := mocksRepo.NewInvestmentRepository(t)
	mockPdfGen := mocksPdf.NewPDFGenerator(t)
	investmentUsecase := NewInvestmentUsecase(mockRepo, mockPdfGen)

	loanID := uuid.New()
	investorID := uuid.New()

	req := &models.CreateInvestmentRequest{
		InvestmentAmount: 2000000, // First investment
	}

	loanInfo := &models.LoanInvestmentInfo{
		ID:              loanID,
		PrincipalAmount: 5000000,
		ROIRate:         8,
		CurrentState:    "APPROVED", // Ready for first investment
		TotalInvested:   0,          // No previous investments
	}

	mockRepo.On("GetLoanForInvestment", mock.Anything, loanID).Return(loanInfo, nil)
	mockRepo.On("CheckExistingInvestment", mock.Anything, loanID, investorID).Return(false, nil)
	mockRepo.On("GetInvestorName", mock.Anything, investorID).Return("Test Investor", nil)
	mockRepo.On("CreateInvestment", mock.Anything, mock.AnythingOfType("*models.Investment")).Return(nil)

	mockPdfGen.On("GenerateInvestmentAgreement",
		mock.AnythingOfType("*models.Investment"),
		mock.AnythingOfType("*models.LoanInvestmentInfo"),
		"Test Investor").Return("/uploads/agreements/investment_agreement.pdf", nil)

	mockRepo.On("UpdateLoanState", mock.Anything, loanID, "FUNDING").Return(nil)

	result, err := investmentUsecase.CreateInvestment(context.Background(), loanID.String(), investorID.String(), req)

	assert.NoError(t, err)
	assert.Equal(t, "FUNDING", result.LoanCurrentState)
	assert.Equal(t, float64(2000000), result.TotalInvestedAmount)
	assert.Equal(t, float64(3000000), result.RemainingAmount) // 5M - 2M = 3M
	assert.Equal(t, "/uploads/agreements/investment_agreement.pdf", result.AgreementURL)
}

// State Transition (FUNDING -> INVESTED)
func TestCreateInvestment_FullInvestmentTransitionsToInvested(t *testing.T) {
	mockRepo := mocksRepo.NewInvestmentRepository(t)
	mockPdfGen := mocksPdf.NewPDFGenerator(t)
	investmentUsecase := NewInvestmentUsecase(mockRepo, mockPdfGen)

	loanID := uuid.New()
	investorID := uuid.New()

	req := &models.CreateInvestmentRequest{
		InvestmentAmount: 2000000, // This completes the funding
	}

	loanInfo := &models.LoanInvestmentInfo{
		ID:              loanID,
		PrincipalAmount: 5000000,
		ROIRate:         8,
		CurrentState:    "FUNDING",
		TotalInvested:   3000000, // Already has 3M, need 2M more
	}

	mockRepo.On("GetLoanForInvestment", mock.Anything, loanID).Return(loanInfo, nil)
	mockRepo.On("CheckExistingInvestment", mock.Anything, loanID, investorID).Return(false, nil)
	mockRepo.On("GetInvestorName", mock.Anything, investorID).Return("Test Investor", nil)
	mockRepo.On("CreateInvestment", mock.Anything, mock.AnythingOfType("*models.Investment")).Return(nil)

	mockPdfGen.On("GenerateInvestmentAgreement",
		mock.AnythingOfType("*models.Investment"),
		mock.AnythingOfType("*models.LoanInvestmentInfo"),
		"Test Investor").Return("/uploads/agreements/investment_agreement.pdf", nil)

	mockRepo.On("UpdateLoanState", mock.Anything, loanID, "INVESTED").Return(nil)

	result, err := investmentUsecase.CreateInvestment(context.Background(), loanID.String(), investorID.String(), req)

	assert.NoError(t, err)
	assert.Equal(t, "INVESTED", result.LoanCurrentState)
	assert.Equal(t, float64(5000000), result.TotalInvestedAmount) // Fully funded
	assert.Equal(t, float64(0), result.RemainingAmount)           // No remaining amount
}

func TestCreateInvestment_ROICalculation(t *testing.T) {
	mockRepo := mocksRepo.NewInvestmentRepository(t)
	mockPdfGen := mocksPdf.NewPDFGenerator(t)
	investmentUsecase := NewInvestmentUsecase(mockRepo, mockPdfGen)

	loanID := uuid.New()
	investorID := uuid.New()

	req := &models.CreateInvestmentRequest{
		InvestmentAmount: 1000000, // 1M investment
	}

	loanInfo := &models.LoanInvestmentInfo{
		ID:              loanID,
		PrincipalAmount: 5000000,
		ROIRate:         12, // 12% ROI rate
		CurrentState:    "APPROVED",
		TotalInvested:   0,
	}

	mockRepo.On("GetLoanForInvestment", mock.Anything, loanID).Return(loanInfo, nil)
	mockRepo.On("CheckExistingInvestment", mock.Anything, loanID, investorID).Return(false, nil)
	mockRepo.On("GetInvestorName", mock.Anything, investorID).Return("Test Investor", nil)

	mockRepo.On("CreateInvestment", mock.Anything, mock.MatchedBy(func(investment *models.Investment) bool {
		expectedReturn := float64(1000000) * (float64(12) / 100) // 1M * 12% = 120K
		return investment.ExpectedReturn == expectedReturn
	})).Return(nil)

	mockPdfGen.On("GenerateInvestmentAgreement",
		mock.AnythingOfType("*models.Investment"),
		mock.AnythingOfType("*models.LoanInvestmentInfo"),
		"Test Investor").Return("/uploads/agreements/investment_agreement.pdf", nil)

	mockRepo.On("UpdateLoanState", mock.Anything, loanID, "FUNDING").Return(nil)

	result, err := investmentUsecase.CreateInvestment(context.Background(), loanID.String(), investorID.String(), req)

	assert.NoError(t, err)
	assert.Equal(t, float64(120000), result.ExpectedReturn)
}

func TestCreateInvestment_PreventOverInvestment(t *testing.T) {
	mockRepo := mocksRepo.NewInvestmentRepository(t)
	mockPdfGen := mocksPdf.NewPDFGenerator(t)
	investmentUsecase := NewInvestmentUsecase(mockRepo, mockPdfGen)

	loanID := uuid.New()
	investorID := uuid.New()

	req := &models.CreateInvestmentRequest{
		InvestmentAmount: 3000000,
	}

	loanInfo := &models.LoanInvestmentInfo{
		ID:              loanID,
		PrincipalAmount: 5000000,
		ROIRate:         8,
		CurrentState:    "FUNDING",
		TotalInvested:   3000000,
	}

	mockRepo.On("GetLoanForInvestment", mock.Anything, loanID).Return(loanInfo, nil)
	mockRepo.On("CheckExistingInvestment", mock.Anything, loanID, investorID).Return(false, nil)

	result, err := investmentUsecase.CreateInvestment(context.Background(), loanID.String(), investorID.String(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "investment amount exceeds remaining loan amount")
}

func TestCreateInvestment_PreventDuplicateInvestment(t *testing.T) {
	mockRepo := mocksRepo.NewInvestmentRepository(t)
	mockPdfGen := mocksPdf.NewPDFGenerator(t)
	investmentUsecase := NewInvestmentUsecase(mockRepo, mockPdfGen)

	loanID := uuid.New()
	investorID := uuid.New()

	req := &models.CreateInvestmentRequest{
		InvestmentAmount: 2000000,
	}

	loanInfo := &models.LoanInvestmentInfo{
		ID:              loanID,
		PrincipalAmount: 5000000,
		ROIRate:         8,
		CurrentState:    "APPROVED",
		TotalInvested:   0,
	}

	mockRepo.On("GetLoanForInvestment", mock.Anything, loanID).Return(loanInfo, nil)
	mockRepo.On("CheckExistingInvestment", mock.Anything, loanID, investorID).Return(true, nil)

	result, err := investmentUsecase.CreateInvestment(context.Background(), loanID.String(), investorID.String(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "investor has already invested in this loan")
}

func TestCreateInvestment_RequiresApprovedOrFundingState(t *testing.T) {
	mockRepo := mocksRepo.NewInvestmentRepository(t)
	mockPdfGen := mocksPdf.NewPDFGenerator(t)
	investmentUsecase := NewInvestmentUsecase(mockRepo, mockPdfGen)

	loanID := uuid.New()
	investorID := uuid.New()

	req := &models.CreateInvestmentRequest{
		InvestmentAmount: 2000000,
	}

	loanInfo := &models.LoanInvestmentInfo{
		ID:              loanID,
		PrincipalAmount: 5000000,
		ROIRate:         8,
		CurrentState:    "PROPOSED", // Invalid state for investment
		TotalInvested:   0,
	}

	mockRepo.On("GetLoanForInvestment", mock.Anything, loanID).Return(loanInfo, nil)

	result, err := investmentUsecase.CreateInvestment(context.Background(), loanID.String(), investorID.String(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "loan must be in APPROVED or FUNDING state")
}

func TestCreateInvestment_InvalidUUIDs(t *testing.T) {
	mockRepo := mocksRepo.NewInvestmentRepository(t)
	mockPdfGen := mocksPdf.NewPDFGenerator(t)
	investmentUsecase := NewInvestmentUsecase(mockRepo, mockPdfGen)

	req := &models.CreateInvestmentRequest{
		InvestmentAmount: 2000000,
	}

	result, err := investmentUsecase.CreateInvestment(context.Background(), "invalid-loan-id", uuid.New().String(), req)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid loan ID")

	// Test invalid investor ID
	result, err = investmentUsecase.CreateInvestment(context.Background(), uuid.New().String(), "invalid-investor-id", req)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid investor ID")
}
