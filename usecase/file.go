package usecase

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fajar-andriansyah/loan-engine/models"
	"github.com/fajar-andriansyah/loan-engine/repositories"
	"github.com/google/uuid"
)

type FileUsecase interface {
	UploadSurveyDocument(ctx context.Context, req *models.UploadDocumentRequest, file multipart.File, header *multipart.FileHeader, validatorID string) (*models.UploadDocumentResponse, error)
}

type fileUsecase struct {
	fileRepo repositories.FileRepository
}

func NewFileUsecase(fileRepo repositories.FileRepository) FileUsecase {
	return &fileUsecase{
		fileRepo: fileRepo,
	}
}

func (u *fileUsecase) UploadSurveyDocument(ctx context.Context, req *models.UploadDocumentRequest, file multipart.File, header *multipart.FileHeader, validatorID string) (*models.UploadDocumentResponse, error) {
	// Parse UUIDs
	loanUUID, err := uuid.Parse(req.LoanID)
	if err != nil {
		return nil, fmt.Errorf("invalid loan ID: %w", err)
	}

	validatorUUID, err := uuid.Parse(validatorID)
	if err != nil {
		return nil, fmt.Errorf("invalid validator ID: %w", err)
	}

	// Parse survey date
	surveyDate, err := time.Parse("2006-01-02", req.SurveyDate)
	if err != nil {
		return nil, fmt.Errorf("invalid survey date format, expected YYYY-MM-DD: %w", err)
	}

	// Check loan state
	currentState, err := u.fileRepo.GetLoanCurrentState(ctx, loanUUID)
	if err != nil {
		return nil, fmt.Errorf("loan not found: %w", err)
	}

	if currentState != "PROPOSED" {
		return nil, fmt.Errorf("loan must be in PROPOSED state, current state: %s", currentState)
	}

	// Validate file type
	if !isValidFileType(header.Filename) {
		return nil, fmt.Errorf("invalid file type, allowed: .jpg, .jpeg, .png, .pdf")
	}

	// Validate file size (10MB)
	if header.Size > 10*1024*1024 {
		return nil, fmt.Errorf("file size exceeds 10MB limit")
	}

	// Generate unique filename
	fileExt := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("survey_%s_%d%s", loanUUID.String(), time.Now().Unix(), fileExt)

	// Create uploads directory if not exists
	uploadDir := "uploads/survey_documents"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Save file to local storage
	filePath := filepath.Join(uploadDir, fileName)
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Generate file URL (assuming served from /uploads path)
	fileURL := fmt.Sprintf("/uploads/survey_documents/%s", fileName)

	// Update loan with survey information
	err = u.fileRepo.UpdateLoanSurveyInfo(ctx, loanUUID, validatorUUID, surveyDate, fileURL, req.SurveyNotes)
	if err != nil {
		// Clean up uploaded file on database error
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to update loan: %w", err)
	}

	return &models.UploadDocumentResponse{
		LoanID:                   loanUUID,
		FileName:                 fileName,
		FileURL:                  fileURL,
		FileType:                 header.Header.Get("Content-Type"),
		FieldValidatorEmployeeID: validatorUUID,
		SurveyDate:               req.SurveyDate,
		SurveyNotes:              req.SurveyNotes,
		UploadedAt:               time.Now(),
	}, nil
}

func isValidFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validTypes := []string{".jpg", ".jpeg", ".png", ".pdf"}

	for _, validType := range validTypes {
		if ext == validType {
			return true
		}
	}
	return false
}
