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
