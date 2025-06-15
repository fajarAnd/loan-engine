package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fajar-andriansyah/loan-engine/infrastructure/database"
	"github.com/fajar-andriansyah/loan-engine/models"
	"github.com/google/uuid"
)

type AuthRepository interface {
	GetEmployeeByEmail(ctx context.Context, email string) (uuid.UUID, *models.EmployeeProfile, string, error)
	GetBorrowerByEmail(ctx context.Context, email string) (uuid.UUID, *models.BorrowerProfile, string, error)
	GetInvestorByEmail(ctx context.Context, email string) (uuid.UUID, *models.InvestorProfile, string, error)
}

type authRepository struct {
	db database.Querier
}

func NewAuthRepository(db database.Querier) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) GetEmployeeByEmail(ctx context.Context, email string) (uuid.UUID, *models.EmployeeProfile, string, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, employee_role, department, is_active
		FROM employees 
		WHERE email = $1 AND is_active = true
	`

	var id uuid.UUID
	var passwordHash string
	var profile models.EmployeeProfile

	err := r.db.QueryRow(ctx, query, email).Scan(
		&id,
		&profile.Username,
		&profile.Email,
		&passwordHash,
		&profile.FullName,
		&profile.EmployeeRole,
		&profile.Department,
		&profile.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, nil, "", fmt.Errorf("employee not found")
		}
		return uuid.Nil, nil, "", fmt.Errorf("failed to get employee: %w", err)
	}

	return id, &profile, passwordHash, nil
}

func (r *authRepository) GetBorrowerByEmail(ctx context.Context, email string) (uuid.UUID, *models.BorrowerProfile, string, error) {
	query := `
		SELECT id, full_name, email, phone_number, identity_number, occupation, password_hash
		FROM borrowers 
		WHERE email = $1
	`

	var id uuid.UUID
	var passwordHash string
	var profile models.BorrowerProfile

	err := r.db.QueryRow(ctx, query, email).Scan(
		&id,
		&profile.FullName,
		&profile.Email,
		&profile.PhoneNumber,
		&profile.IdentityNumber,
		&profile.Occupation,
		&passwordHash,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, nil, "", fmt.Errorf("borrower not found")
		}
		return uuid.Nil, nil, "", fmt.Errorf("failed to get borrower: %w", err)
	}

	return id, &profile, passwordHash, nil
}

func (r *authRepository) GetInvestorByEmail(ctx context.Context, email string) (uuid.UUID, *models.InvestorProfile, string, error) {
	query := `
		SELECT id, full_name, email, phone_number, identity_number, is_active, password_hash
		FROM investors 
		WHERE email = $1 AND is_active = true
	`

	var id uuid.UUID
	var passwordHash string
	var profile models.InvestorProfile

	err := r.db.QueryRow(ctx, query, email).Scan(
		&id,
		&profile.FullName,
		&profile.Email,
		&profile.PhoneNumber,
		&profile.IdentityNumber,
		&profile.IsActive,
		&passwordHash,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, nil, "", fmt.Errorf("investor not found")
		}
		return uuid.Nil, nil, "", fmt.Errorf("failed to get investor: %w", err)
	}

	return id, &profile, passwordHash, nil
}
