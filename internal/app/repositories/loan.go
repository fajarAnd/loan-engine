package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fajar-andriansyah/loan-engine/internal/app/constants"
	"github.com/fajar-andriansyah/loan-engine/internal/app/database"
	"github.com/fajar-andriansyah/loan-engine/internal/app/models"
	"github.com/google/uuid"
)

type LoanRepository interface {
	CreateLoan(ctx context.Context, loan *models.Loan) error
	GetLoanForApproval(ctx context.Context, loanID uuid.UUID) (*models.LoanForApproval, error)
	ApproveLoan(ctx context.Context, loanID, approvingEmployeeID uuid.UUID, approvalNotes, agreementURL string) error
	GetApprovedLoan(ctx context.Context, loanID uuid.UUID) (*models.ApproveLoanResponse, error)
	GetLoanForDisbursement(ctx context.Context, loanID uuid.UUID) (*models.Loan, error)
	DisburseLoan(ctx context.Context, loanID, fieldOfficerID uuid.UUID, signedAgreementURL, disbursementNotes string) error
	GetDisbursedLoan(ctx context.Context, loanID uuid.UUID) (*models.DisburseLoanResponse, error)
}

type loanRepository struct {
	db database.Querier
}

func NewLoanRepository(db database.Querier) LoanRepository {
	return &loanRepository{
		db: db,
	}
}

func (r *loanRepository) CreateLoan(ctx context.Context, loan *models.Loan) error {
	query := `
		INSERT INTO loans (
			id, borrower_id, principal_amount, interest_rate, roi_rate,
			loan_term_month, current_state, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	if db, ok := r.db.(database.Executor); ok {
		_, err := db.Exec(ctx, query,
			loan.ID,
			loan.BorrowerID,
			loan.PrincipalAmount,
			loan.InterestRate,
			loan.ROIRate,
			loan.LoanTermMonth,
			loan.CurrentState,
			loan.CreatedAt,
			loan.UpdatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to create loan: %w", err)
		}
	}

	return nil
}

func (r *loanRepository) GetLoanForApproval(ctx context.Context, loanID uuid.UUID) (*models.LoanForApproval, error) {
	query := `
		SELECT
			l.id, l.borrower_id, l.principal_amount, l.interest_rate, l.roi_rate,
			l.loan_term_month, l.current_state, l.field_validator_employee_id,
			l.survey_date, b.full_name as borrower_name
		FROM loans l
		JOIN borrowers b ON l.borrower_id = b.id
		WHERE l.id = $1
	`

	var loan models.LoanForApproval
	var surveyDate sql.NullTime
	var validatorID sql.NullString

	err := r.db.QueryRow(ctx, query, loanID).Scan(
		&loan.ID,
		&loan.BorrowerID,
		&loan.PrincipalAmount,
		&loan.InterestRate,
		&loan.ROIRate,
		&loan.LoanTermMonth,
		&loan.CurrentState,
		&validatorID,
		&surveyDate,
		&loan.BorrowerName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("loan not found")
		}
		return nil, fmt.Errorf("failed to get loan: %w", err)
	}

	if !validatorID.Valid || !surveyDate.Valid {
		return nil, fmt.Errorf("survey not completed")
	}

	loan.FieldValidatorEmployeeID = uuid.MustParse(validatorID.String)
	loan.SurveyDate = surveyDate.Time

	return &loan, nil
}

func (r *loanRepository) ApproveLoan(ctx context.Context, loanID, approvingEmployeeID uuid.UUID, approvalNotes, agreementURL string) error {
	query := `
		UPDATE loans 
		SET current_state = $5,
		    approving_employee_id = $2,
		    approval_date = CURRENT_DATE,
		    approval_notes = $3,
		    loan_agreement_pdf_url = $4,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND current_state = $6
	`

	if db, ok := r.db.(database.Executor); ok {
		result, err := db.Exec(ctx, query, loanID, approvingEmployeeID, approvalNotes, agreementURL, constants.APPROVED, constants.PROPOSED)
		if err != nil {
			return fmt.Errorf("failed to approve loan: %w", err)
		}

		if result.RowsAffected() == 0 {
			return fmt.Errorf("loan not found")
		}
	}

	return nil
}

func (r *loanRepository) GetApprovedLoan(ctx context.Context, loanID uuid.UUID) (*models.ApproveLoanResponse, error) {
	query := `
		SELECT 
			id, borrower_id, principal_amount, interest_rate, roi_rate,
			loan_term_month, current_state, approval_date, approving_employee_id,
			approval_notes, loan_agreement_pdf_url, field_validator_employee_id,
			survey_date, updated_at
		FROM loans 
		WHERE id = $1
	`

	var response models.ApproveLoanResponse
	var approvalDate sql.NullTime
	var approvalNotes sql.NullString
	var surveyDate sql.NullTime

	err := r.db.QueryRow(ctx, query, loanID).Scan(
		&response.ID,
		&response.BorrowerID,
		&response.PrincipalAmount,
		&response.InterestRate,
		&response.ROIRate,
		&response.LoanTermMonth,
		&response.CurrentState,
		&approvalDate,
		&response.ApprovingEmployeeID,
		&approvalNotes,
		&response.LoanAgreementPDFURL,
		&response.FieldValidatorEmployeeID,
		&surveyDate,
		&response.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get approved loan: %w", err)
	}

	if approvalDate.Valid {
		response.ApprovalDate = approvalDate.Time.Format("2006-01-02")
	}
	if approvalNotes.Valid {
		response.ApprovalNotes = approvalNotes.String
	}
	if surveyDate.Valid {
		response.SurveyDate = surveyDate.Time.Format("2006-01-02")
	}

	return &response, nil
}

func (r *loanRepository) GetLoanForDisbursement(ctx context.Context, loanID uuid.UUID) (*models.Loan, error) {
	query := `
		SELECT id, borrower_id, principal_amount, interest_rate, roi_rate,
		       loan_term_month, current_state, created_at, updated_at
		FROM loans 
		WHERE id = $1
	`

	var loan models.Loan
	err := r.db.QueryRow(ctx, query, loanID).Scan(
		&loan.ID,
		&loan.BorrowerID,
		&loan.PrincipalAmount,
		&loan.InterestRate,
		&loan.ROIRate,
		&loan.LoanTermMonth,
		&loan.CurrentState,
		&loan.CreatedAt,
		&loan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("loan not found")
		}
		return nil, fmt.Errorf("failed to get loan: %w", err)
	}

	return &loan, nil
}

func (r *loanRepository) DisburseLoan(ctx context.Context, loanID, fieldOfficerID uuid.UUID, signedAgreementURL, disbursementNotes string) error {
	query := `
		UPDATE loans 
		SET current_state = $5,
		    field_officer_employee_id = $2,
		    disbursement_date = CURRENT_DATE,
		    signed_agreement_url = $3,
		    disbursement_notes = $4,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND current_state = $6
	`

	if db, ok := r.db.(database.Executor); ok {
		result, err := db.Exec(ctx, query, loanID, fieldOfficerID, signedAgreementURL, disbursementNotes, constants.DISBURSED, constants.INVESTED)
		if err != nil {
			return fmt.Errorf("failed to disburse loan: %w", err)
		}

		if result.RowsAffected() == 0 {
			return fmt.Errorf("loan not found")
		}
	}
	return nil
}

func (r *loanRepository) GetDisbursedLoan(ctx context.Context, loanID uuid.UUID) (*models.DisburseLoanResponse, error) {
	query := `
		SELECT id, borrower_id, principal_amount, interest_rate, roi_rate,
		       loan_term_month, current_state, disbursement_date, 
		       field_officer_employee_id, signed_agreement_url, 
		       disbursement_notes, updated_at
		FROM loans 
		WHERE id = $1
	`

	var response models.DisburseLoanResponse
	var disbursementDate sql.NullTime
	var disbursementNotes sql.NullString

	err := r.db.QueryRow(ctx, query, loanID).Scan(
		&response.ID,
		&response.BorrowerID,
		&response.PrincipalAmount,
		&response.InterestRate,
		&response.ROIRate,
		&response.LoanTermMonth,
		&response.CurrentState,
		&disbursementDate,
		&response.FieldOfficerEmployeeID,
		&response.SignedAgreementURL,
		&disbursementNotes,
		&response.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get disbursed loan: %w", err)
	}

	if disbursementDate.Valid {
		response.DisbursementDate = disbursementDate.Time.Format("2006-01-02")
	}
	if disbursementNotes.Valid {
		response.DisbursementNotes = disbursementNotes.String
	}

	return &response, nil
}
