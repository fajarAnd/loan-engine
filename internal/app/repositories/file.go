package repositories

import (
	"context"
	"fmt"
	"github.com/fajar-andriansyah/loan-engine/internal/app/database"
	"github.com/google/uuid"
	"time"
)

type FileRepository interface {
	UpdateLoanSurveyInfo(ctx context.Context, loanID uuid.UUID, validatorID uuid.UUID, surveyDate time.Time, fileURL, surveyNotes string) error
	GetLoanCurrentState(ctx context.Context, loanID uuid.UUID) (string, error)
}

type fileRepository struct {
	db database.Querier
}

func NewFileRepository(db database.Querier) FileRepository {
	return &fileRepository{
		db: db,
	}
}

func (r *fileRepository) GetLoanCurrentState(ctx context.Context, loanID uuid.UUID) (string, error) {
	query := `SELECT current_state FROM loans WHERE id = $1`

	var currentState string
	err := r.db.QueryRow(ctx, query, loanID).Scan(&currentState)
	if err != nil {
		return "", fmt.Errorf("failed to get loan state: %w", err)
	}

	return currentState, nil
}

func (r *fileRepository) UpdateLoanSurveyInfo(ctx context.Context, loanID uuid.UUID, validatorID uuid.UUID, surveyDate time.Time, fileURL, surveyNotes string) error {
	query := `
		UPDATE loans 
		SET field_validator_employee_id = $2,
		    survey_date = $3,
		    field_visit_proof_url = $4,
		    survey_notes = $5,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND current_state = 'PROPOSED'
	`

	if db, ok := r.db.(database.Executor); ok {
		result, err := db.Exec(ctx, query, loanID, validatorID, surveyDate, fileURL, surveyNotes)
		if err != nil {
			return fmt.Errorf("failed to update loan survey info: %w", err)
		}

		if result.RowsAffected() == 0 {
			return fmt.Errorf("loan not found or not in PROPOSED state")
		}
	} else {
		return fmt.Errorf("database does not support Exec operation")
	}

	return nil
}
