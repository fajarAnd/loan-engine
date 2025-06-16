package models

import (
	_ "fmt"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	UserType string `json:"user_type" validate:"required,oneof=employee borrower investor"`
}

type LoginResponse struct {
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	ExpiresIn   int      `json:"expires_in"`
	User        UserInfo `json:"user"`
}

type UserInfo struct {
	ID       string      `json:"id"`
	UserType string      `json:"user_type"`
	Profile  interface{} `json:"profile"`
}

type EmployeeProfile struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	EmployeeRole string `json:"employee_role"`
	Department   string `json:"department"`
	IsActive     bool   `json:"is_active"`
}

type BorrowerProfile struct {
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	IdentityNumber string `json:"identity_number"`
	Occupation     string `json:"occupation"`
}

type InvestorProfile struct {
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	IdentityNumber string `json:"identity_number"`
	IsActive       bool   `json:"is_active"`
}

type JWTClaims struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
	Role     string `json:"role,omitempty"`
	jwt.RegisteredClaims
}
