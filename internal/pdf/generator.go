// internal/pdf/generator.go
package pdf

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fajar-andriansyah/loan-engine/models"
	"github.com/jung-kurt/gofpdf"
)

func GenerateLoanAgreement(loan *models.LoanForApproval) (string, error) {
	// Create agreements directory if not exists
	agreementDir := "uploads/agreements"
	if err := os.MkdirAll(agreementDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create agreement directory: %w", err)
	}

	// Generate filename
	fileName := fmt.Sprintf("loan_agreement_%s.pdf", loan.ID.String())
	filePath := filepath.Join(agreementDir, fileName)

	// Calculate total payment
	totalAmount := loan.PrincipalAmount * (1 + loan.InterestRate/100)
	monthlyPayment := totalAmount / float64(loan.LoanTermMonth)

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "LOAN AGREEMENT")
	pdf.Ln(15)

	// Basic info
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Loan ID: %s", loan.ID.String()))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Date: %s", time.Now().Format("2006-01-02")))
	pdf.Ln(15)

	// Borrower information
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "BORROWER INFORMATION:")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Name: %s", loan.BorrowerName))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Borrower ID: %s", loan.BorrowerID.String()))
	pdf.Ln(15)

	// Loan details
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "LOAN DETAILS:")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Principal Amount: Rp %.2f", loan.PrincipalAmount))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Interest Rate: %.2f%% per annum", loan.InterestRate))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("ROI Rate: %.2f%% per annum", loan.ROIRate))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Loan Term: %d months", loan.LoanTermMonth))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Total Amount: Rp %.2f", totalAmount))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Monthly Payment: Rp %.2f", monthlyPayment))
	pdf.Ln(15)

	// Terms and conditions
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "TERMS AND CONDITIONS:")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, "1. The borrower agrees to repay the loan in monthly installments")
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("2. Payment schedule: %d months total", loan.LoanTermMonth))
	pdf.Ln(8)
	pdf.Cell(0, 8, "3. Interest is calculated on flat rate basis")
	pdf.Ln(8)
	pdf.Cell(0, 8, "4. Late payment may incur additional charges")
	pdf.Ln(15)

	// Approval information
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "APPROVAL INFORMATION:")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Survey Date: %s", loan.SurveyDate.Format("2006-01-02")))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Field Validator: %s", loan.FieldValidatorEmployeeID.String()))
	pdf.Ln(15)

	pdf.Cell(0, 8, fmt.Sprintf("This agreement is generated on %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(20)

	// Signature section
	pdf.Cell(90, 8, "_________________________")
	pdf.Cell(90, 8, "_________________________")
	pdf.Ln(8)
	pdf.Cell(90, 8, "Borrower Signature")
	pdf.Cell(90, 8, "Amartha Representative")

	// Save PDF
	err := pdf.OutputFileAndClose(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Return URL path
	return fmt.Sprintf("/uploads/agreements/%s", fileName), nil
}

func GenerateInvestmentAgreement(investment *models.Investment, loan *models.LoanInvestmentInfo, investorName string) (string, error) {
	// Create agreements directory if not exists
	agreementDir := "uploads/agreements"
	if err := os.MkdirAll(agreementDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create agreement directory: %w", err)
	}

	// Generate filename
	fileName := fmt.Sprintf("investment_agreement_%s_%s.pdf", investment.LoanID.String(), investment.InvestorID.String())
	filePath := filepath.Join(agreementDir, fileName)

	// Calculate return details
	investmentPeriod := 12 // months (could be from loan data)
	totalReturn := investment.InvestmentAmount + investment.ExpectedReturn
	monthlyReturn := totalReturn / float64(investmentPeriod)

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "INVESTMENT AGREEMENT")
	pdf.Ln(15)

	// Basic info
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Agreement ID: %s", investment.ID.String()))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Date: %s", time.Now().Format("2006-01-02")))
	pdf.Ln(15)

	// Investment information
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "INVESTMENT DETAILS:")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Investor Name: %s", investorName))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Investor ID: %s", investment.InvestorID.String()))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Loan ID: %s", investment.LoanID.String()))
	pdf.Ln(15)

	// Financial details
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "FINANCIAL TERMS:")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Investment Amount: Rp %.2f", investment.InvestmentAmount))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Expected Return: Rp %.2f", investment.ExpectedReturn))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("ROI Rate: %.2f%% per annum", loan.ROIRate))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Investment Period: %d months", investmentPeriod))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Total Return: Rp %.2f", totalReturn))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Monthly Return: Rp %.2f", monthlyReturn))
	pdf.Ln(15)

	// Terms and conditions
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "TERMS AND CONDITIONS:")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, "1. This investment is for peer-to-peer lending facilitated by Amartha")
	pdf.Ln(8)
	pdf.Cell(0, 8, "2. Returns are subject to borrower's ability to repay the loan")
	pdf.Ln(8)
	pdf.Cell(0, 8, "3. Amartha acts as facilitator and is not liable for borrower default")
	pdf.Ln(8)
	pdf.Cell(0, 8, "4. Returns will be paid monthly as per the loan repayment schedule")
	pdf.Ln(8)
	pdf.Cell(0, 8, "5. This agreement is governed by Indonesian law")
	pdf.Ln(15)

	// Investment date
	pdf.Cell(0, 8, fmt.Sprintf("Investment Date: %s", investment.InvestmentDate.Format("2006-01-02")))
	pdf.Ln(8)
	pdf.Cell(0, 8, fmt.Sprintf("Agreement generated on: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(20)

	// Signature section
	pdf.Cell(90, 8, "_________________________")
	pdf.Cell(90, 8, "_________________________")
	pdf.Ln(8)
	pdf.Cell(90, 8, "Investor Signature")
	pdf.Cell(90, 8, "Amartha Representative")

	// Save PDF
	err := pdf.OutputFileAndClose(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Return URL path
	return fmt.Sprintf("/uploads/agreements/%s", fileName), nil
}
