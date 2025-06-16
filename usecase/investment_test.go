package usecase

import (
	"context"
	"testing"

	"github.com/fajar-andriansyah/loan-engine/mocks/repositories"
	"github.com/fajar-andriansyah/loan-engine/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test Business Logic: First Investment State Transition (APPROVED -> FUNDING)
func TestCreateInvestment_FirstInvestmentTransitionsToFunding(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewInvestmentRepository(t)
	investmentUsecase := NewInvestmentUsecase(mockRepo)

	loanID := uuid.New()
	investorID := uuid.New()

	req := &models.CreateInvestmentRequest{
		InvestmentAmount: 2000000, // First investment
	}

	// Mock loan in APPROVED state with no investments
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
	// Critical: State must transition to FUNDING on first investment
	mockRepo.On("UpdateLoanState", mock.Anything, loanID, "FUNDING").Return(nil)

	// Act
	result, err := investmentUsecase.CreateInvestment(context.Background(), loanID.String(), investorID.String(), req)

	// Assert - Critical business rule: First investment triggers FUNDING state
	assert.NoError(t, err)
	assert.Equal(t, "FUNDING", result.LoanCurrentState)
	assert.Equal(t, float64(2000000), result.TotalInvestedAmount)
	assert.Equal(t, float64(3000000), result.RemainingAmount) // 5M - 2M = 3M
}

// Test Business Logic: Full Investment State Transition (FUNDING -> INVESTED)
func TestCreateInvestment_FullInvestmentTransitionsToInvested(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewInvestmentRepository(t)
	investmentUsecase := NewInvestmentUsecase(mockRepo)

	loanID := uuid.New()
	investorID := uuid.New()

	req := &models.CreateInvestmentRequest{
		InvestmentAmount: 2000000, // This completes the funding
	}

	// Mock loan in FUNDING state with partial investment
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
	// Critical: State must transition to INVESTED when fully funded
	mockRepo.On("UpdateLoanState", mock.Anything, loanID, "INVESTED").Return(nil)

	// Act
	result, err := investmentUsecase.CreateInvestment(context.Background(), loanID.String(), investorID.String(), req)

	// Assert - Critical business rule: Full funding triggers INVESTED state
	assert.NoError(t, err)
	assert.Equal(t, "INVESTED", result.LoanCurrentState)
	assert.Equal(t, float64(5000000), result.TotalInvestedAmount) // Fully funded
	assert.Equal(t, float64(0), result.RemainingAmount)           // No remaining amount
}

// Test Business Logic: ROI Calculation
func TestCreateInvestment_ROICalculationAccuracy(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewInvestmentRepository(t)
	investmentUsecase := NewInvestmentUsecase(mockRepo)

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
		// Critical: Expected return must be calculated correctly
		expectedReturn := 1000000 * 0.12 // 1M * 12% = 120K
		return investment.ExpectedReturn == expectedReturn
	})).Return(nil)
	mockRepo.On("UpdateLoanState", mock.Anything, loanID, "FUNDING").Return(nil)

	// Act
	result, err := investmentUsecase.CreateInvestment(context.Background(), loanID.String(), investorID.String(), req)

	// Assert - Critical business rule: ROI calculation must be accurate
	assert.NoError(t, err)
	assert.Equal(t, float64(120000), result.ExpectedReturn) // 1M * 12% = 120K
}

// Test Business Logic: Investment Amount Validation
func TestCreateInvestment_PreventOverInvestment(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewInvestmentRepository(t)
	investmentUsecase := NewInvestmentUsecase(mockRepo)

	loanID := uuid.New()
	investorID := uuid.New()

	req := &models.CreateInvestmentRequest{
		InvestmentAmount: 3000000, // Exceeds remaining 2M
	}

	loanInfo := &models.LoanInvestmentInfo{
		ID:              loanID,
		PrincipalAmount: 5000000,
		ROIRate:         8,
		CurrentState:    "FUNDING",
		TotalInvested:   3000000, // Only 2M remaining
	}

	mockRepo.On("GetLoanForInvestment", mock.Anything, loanID).Return(loanInfo, nil)
	mockRepo.On("CheckExistingInvestment", mock.Anything, loanID, investorID).Return(false, nil)

	// Act
	result, err := investmentUsecase.CreateInvestment(context.Background(), loanID.String(), investorID.String(), req)

	// Assert - Critical business rule: Cannot exceed principal amount
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "investment amount exceeds remaining loan amount")
}

// Test Business Logic: Duplicate Investment Prevention
func TestCreateInvestment_PreventDuplicateInvestment(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewInvestmentRepository(t)
	investmentUsecase := NewInvestmentUsecase(mockRepo)

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
	// Critical: Check returns true (investor already invested)
	mockRepo.On("CheckExistingInvestment", mock.Anything, loanID, investorID).Return(true, nil)

	// Act
	result, err := investmentUsecase.CreateInvestment(context.Background(), loanID.String(), investorID.String(), req)

	// Assert - Critical business rule: One investment per investor per loan
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "investor has already invested in this loan")
}

// Test Business Logic: Investment State Validation
func TestCreateInvestment_RequiresApprovedOrFundingState(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewInvestmentRepository(t)
	investmentUsecase := NewInvestmentUsecase(mockRepo)

	loanID := uuid.New()
	investorID := uuid.New()

	req := &models.CreateInvestmentRequest{
		InvestmentAmount: 2000000,
	}

	// Test with PROPOSED state (invalid)
	loanInfo := &models.LoanInvestmentInfo{
		ID:              loanID,
		PrincipalAmount: 5000000,
		ROIRate:         8,
		CurrentState:    "PROPOSED", // Invalid state for investment
		TotalInvested:   0,
	}

	mockRepo.On("GetLoanForInvestment", mock.Anything, loanID).Return(loanInfo, nil)

	// Act
	result, err := investmentUsecase.CreateInvestment(context.Background(), loanID.String(), investorID.String(), req)

	// Assert - Critical business rule: Only APPROVED or FUNDING loans can receive investments
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "loan must be in APPROVED or FUNDING state")
}

// Test Business Logic: Invalid UUID Validation
func TestCreateInvestment_InvalidUUIDs(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewInvestmentRepository(t)
	investmentUsecase := NewInvestmentUsecase(mockRepo)

	req := &models.CreateInvestmentRequest{
		InvestmentAmount: 2000000,
	}

	// Test invalid loan ID
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
