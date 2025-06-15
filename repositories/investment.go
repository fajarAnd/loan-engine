package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fajar-andriansyah/loan-engine/infrastructure/database"
	"github.com/fajar-andriansyah/loan-engine/models"
	"github.com/google/uuid"
)

type InvestmentRepository interface {
	GetLoanForInvestment(ctx context.Context, loanID uuid.UUID) (*models.LoanInvestmentInfo, error)
	CheckExistingInvestment(ctx context.Context, loanID, investorID uuid.UUID) (bool, error)
	CreateInvestment(ctx context.Context, investment *models.Investment) error
	UpdateLoanState(ctx context.Context, loanID uuid.UUID, newState string) error
	GetTotalInvestedAmount(ctx context.Context, loanID uuid.UUID) (float64, error)
	GetInvestorName(ctx context.Context, investorID uuid.UUID) (string, error)
}

type investmentRepository struct {
	db database.Querier
}

func NewInvestmentRepository(db database.Querier) InvestmentRepository {
	return &investmentRepository{
		db: db,
	}
}

func (r *investmentRepository) GetLoanForInvestment(ctx context.Context, loanID uuid.UUID) (*models.LoanInvestmentInfo, error) {
	query := `
		SELECT 
			l.id, l.principal_amount, l.roi_rate, l.current_state,
			COALESCE(SUM(i.investment_amount), 0) as total_invested
		FROM loans l
		LEFT JOIN investments i ON l.id = i.loan_id
		WHERE l.id = $1
		GROUP BY l.id, l.principal_amount, l.roi_rate, l.current_state
	`

	var loan models.LoanInvestmentInfo
	err := r.db.QueryRow(ctx, query, loanID).Scan(
		&loan.ID,
		&loan.PrincipalAmount,
		&loan.ROIRate,
		&loan.CurrentState,
		&loan.TotalInvested,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("loan not found")
		}
		return nil, fmt.Errorf("failed to get loan: %w", err)
	}

	return &loan, nil
}

func (r *investmentRepository) CheckExistingInvestment(ctx context.Context, loanID, investorID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM investments WHERE loan_id = $1 AND investor_id = $2)`

	var exists bool
	err := r.db.QueryRow(ctx, query, loanID, investorID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check existing investment: %w", err)
	}

	return exists, nil
}

func (r *investmentRepository) CreateInvestment(ctx context.Context, investment *models.Investment) error {
	query := `
		INSERT INTO investments (
			id, loan_id, investor_id, investment_amount, expected_return,
			investment_date, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if db, ok := r.db.(database.Executor); ok {
		_, err := db.Exec(ctx, query,
			investment.ID,
			investment.LoanID,
			investment.InvestorID,
			investment.InvestmentAmount,
			investment.ExpectedReturn,
			investment.InvestmentDate,
			investment.CreatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to create investment: %w", err)
		}
	} else {
		return fmt.Errorf("database does not support Exec operation")
	}

	return nil
}

func (r *investmentRepository) UpdateLoanState(ctx context.Context, loanID uuid.UUID, newState string) error {
	query := `
		UPDATE loans 
		SET current_state = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	if db, ok := r.db.(database.Executor); ok {
		result, err := db.Exec(ctx, query, loanID, newState)
		if err != nil {
			return fmt.Errorf("failed to update loan state: %w", err)
		}

		if result.RowsAffected() == 0 {
			return fmt.Errorf("loan not found")
		}
	} else {
		return fmt.Errorf("database does not support Exec operation")
	}

	return nil
}

func (r *investmentRepository) GetTotalInvestedAmount(ctx context.Context, loanID uuid.UUID) (float64, error) {
	query := `SELECT COALESCE(SUM(investment_amount), 0) FROM investments WHERE loan_id = $1`

	var total float64
	err := r.db.QueryRow(ctx, query, loanID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total invested amount: %w", err)
	}

	return total, nil
}

func (r *investmentRepository) GetInvestorName(ctx context.Context, investorID uuid.UUID) (string, error) {
	query := `SELECT full_name FROM investors WHERE id = $1`

	var name string
	err := r.db.QueryRow(ctx, query, investorID).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("investor not found")
		}
		return "", fmt.Errorf("failed to get investor name: %w", err)
	}

	return name, nil
}
