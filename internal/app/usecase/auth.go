package usecase

import (
	"context"
	"fmt"
	"github.com/fajar-andriansyah/loan-engine/internal/app/constants"
	"github.com/fajar-andriansyah/loan-engine/internal/app/models"
	"github.com/fajar-andriansyah/loan-engine/internal/app/repositories"
	"github.com/google/uuid"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
}

type authUsecase struct {
	authRepo  repositories.AuthRepository
	jwtSecret string
}

func NewAuthUsecase(authRepo repositories.AuthRepository, jwtSecret string) AuthUsecase {
	return &authUsecase{
		authRepo:  authRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *authUsecase) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	var userID uuid.UUID
	var profile interface{}
	var passwordHash string
	var role string
	var err error

	switch req.UserType {
	case constants.USER_EMPLOYEE:
		id, empProfile, hash, err := u.authRepo.GetEmployeeByEmail(ctx, req.Email)
		if err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
		userID = id
		profile = empProfile
		passwordHash = hash
		role = empProfile.EmployeeRole

	case constants.USER_BORROWER:
		id, borProfile, hash, err := u.authRepo.GetBorrowerByEmail(ctx, req.Email)
		if err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
		userID = id
		profile = borProfile
		passwordHash = hash

	case constants.USER_INVESTOR:
		id, invProfile, hash, err := u.authRepo.GetInvestorByEmail(ctx, req.Email)
		if err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
		userID = id
		profile = invProfile
		passwordHash = hash

	default:
		return nil, fmt.Errorf("invalid user type")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	expiresIn := 3600 // 1 hour
	claims := &models.JWTClaims{
		UserID:   userID.String(),
		UserType: req.UserType,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		User: models.UserInfo{
			ID:       userID.String(),
			UserType: req.UserType,
			Profile:  profile,
		},
	}, nil
}
