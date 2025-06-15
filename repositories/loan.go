package repositories

import (
	"context"
	"fmt"
	"github.com/fajar-andriansyah/loan-engine/infrastructure/database"
	"github.com/fajar-andriansyah/loan-engine/models"
)

type LoanRepository interface {
	CreateLoan(ctx context.Context, loan *models.Loan) error
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
	} else {
		return fmt.Errorf("database does not support Exec operation")
	}

	return nil
}
