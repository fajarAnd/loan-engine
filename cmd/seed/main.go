package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/fajar-andriansyah/loan-engine/config"
	database2 "github.com/fajar-andriansyah/loan-engine/infrastructure/database"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Employee struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	PhoneNumber  string    `json:"phone_number"`
	Role         string    `json:"employee_role"`
	Department   string    `json:"department"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Investor struct {
	ID             uuid.UUID `json:"id"`
	FullName       string    `json:"full_name"`
	IdentityNumber string    `json:"identity_number"`
	Email          string    `json:"email"`
	PhoneNumber    string    `json:"phone_number"`
	Address        string    `json:"address"`
	BankAccount    string    `json:"bank_account"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	PasswordHash   string    `json:"password_hash"`
}

type Borrower struct {
	ID                uuid.UUID `json:"id"`
	FullName          string    `json:"full_name"`
	IdentityNumber    string    `json:"identity_number"`
	PhoneNumber       string    `json:"phone_number"`
	Email             string    `json:"email"`
	Address           string    `json:"address"`
	DateOfBirth       time.Time `json:"date_of_birth"`
	Occupation        string    `json:"occupation"`
	MonthlyIncome     float64   `json:"monthly_income"`
	BankAccountNumber string    `json:"bank_account_number"`
	BankName          string    `json:"bank_name"`
	AccountHolderName string    `json:"account_holder_name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	PasswordHash      string    `json:"password_hash"`
}

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set log level to info for better output
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Connect to database
	dbConfig := database2.LoadDatabaseConfig()
	db, err := database2.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	log.Println(" Starting database seeding...")

	// Seed Employees
	if err := seedEmployees(ctx, db); err != nil {
		log.Fatalf("Failed to seed employees: %v", err)
	}

	// Seed Investors
	if err := seedInvestors(ctx, db); err != nil {
		log.Fatalf("Failed to seed investors: %v", err)
	}

	// Seed Borrowers
	if err := seedBorrowers(ctx, db); err != nil {
		log.Fatalf("Failed to seed borrowers: %v", err)
	}

	log.Println(" Database seeding completed successfully!")
}

func seedEmployees(ctx context.Context, db *sql.DB) error {
	log.Println(" Seeding employees...")

	employees := []Employee{
		{
			ID:           uuid.New(),
			Username:     "field_validator_001",
			Email:        "validator@amartha.com",
			PasswordHash: "$2a$10$WYmkE5HjRSYrJcUS9OGL/u9biq0iYc6GoUPiYLVd1UvO8hZPo98fO", // password123
			FullName:     "Ahmad Validator",
			PhoneNumber:  "+6281234567890",
			Role:         "FIELD_VALIDATOR",
			Department:   "Field Operations",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Username:     "field_officer_001",
			Email:        "officer@amartha.com",
			PasswordHash: "$2a$10$WYmkE5HjRSYrJcUS9OGL/u9biq0iYc6GoUPiYLVd1UvO8hZPo98fO", // password123
			FullName:     "Sari Fieldofficer",
			PhoneNumber:  "+6281234567891",
			Role:         "FIELD_OFFICER",
			Department:   "Field Operations",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Username:     "admin_001",
			Email:        "admin@amartha.com",
			PasswordHash: "$2a$10$WYmkE5HjRSYrJcUS9OGL/u9biq0iYc6GoUPiYLVd1UvO8hZPo98fO", // password123
			FullName:     "Budi Administrator",
			PhoneNumber:  "+6281234567892",
			Role:         "ADMIN",
			Department:   "IT Operations",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	query := `
		INSERT INTO employees (
			id, username, email, password_hash, full_name, phone_number,
			employee_role, department, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (username) DO NOTHING`

	for _, emp := range employees {
		_, err := db.ExecContext(ctx, query,
			emp.ID, emp.Username, emp.Email, emp.PasswordHash, emp.FullName,
			emp.PhoneNumber, emp.Role, emp.Department, emp.IsActive,
			emp.CreatedAt, emp.UpdatedAt,
		)
		if err != nil {
			return err
		}
		log.Printf("   Created employee: %s (%s)", emp.FullName, emp.Role)
	}

	return nil
}

func seedInvestors(ctx context.Context, db *sql.DB) error {
	log.Println(" Seeding investors...")

	investors := []Investor{
		{
			ID:             uuid.New(),
			FullName:       "Rina Investor",
			IdentityNumber: "3201234567890001",
			Email:          "rina.investor@gmail.com",
			PhoneNumber:    "+6281234567893",
			Address:        "Jl. Investor Raya No. 123, Jakarta Selatan",
			BankAccount:    "1234567890 - Bank Mandiri",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			PasswordHash:   "$2a$10$WYmkE5HjRSYrJcUS9OGL/u9biq0iYc6GoUPiYLVd1UvO8hZPo98fO", // password123
		},
		{
			ID:             uuid.New(),
			FullName:       "Doni Kapital",
			IdentityNumber: "3201234567890002",
			Email:          "doni.kapital@gmail.com",
			PhoneNumber:    "+6281234567894",
			Address:        "Jl. Kebon Jeruk No. 456, Jakarta Barat",
			BankAccount:    "0987654321 - Bank BCA",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			PasswordHash:   "$2a$10$WYmkE5HjRSYrJcUS9OGL/u9biq0iYc6GoUPiYLVd1UvO8hZPo98fO", // password123
		},
	}

	query := `
		INSERT INTO investors (
			id, full_name, identity_number, email, phone_number,
			address, bank_account, is_active, created_at, updated_at, password_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (identity_number) DO NOTHING`

	for _, inv := range investors {
		_, err := db.ExecContext(ctx, query,
			inv.ID, inv.FullName, inv.IdentityNumber, inv.Email,
			inv.PhoneNumber, inv.Address, inv.BankAccount,
			inv.IsActive, inv.CreatedAt, inv.UpdatedAt, inv.PasswordHash,
		)
		if err != nil {
			return err
		}
		log.Printf("  Created investor: %s", inv.FullName)
	}

	return nil
}

func seedBorrowers(ctx context.Context, db *sql.DB) error {
	log.Println(" Seeding borrowers...")

	borrowers := []Borrower{
		{
			ID:                uuid.New(),
			FullName:          "Siti Peminjam",
			IdentityNumber:    "3201234567890003",
			PhoneNumber:       "+6281234567895",
			Email:             "siti.peminjam@gmail.com",
			Address:           "Jl. Usaha Mikro No. 789, Tangerang Selatan",
			DateOfBirth:       time.Date(1985, 3, 15, 0, 0, 0, 0, time.UTC),
			Occupation:        "Penjual Makanan",
			MonthlyIncome:     2500000.00,
			BankAccountNumber: "1122334455",
			BankName:          "Bank BRI",
			AccountHolderName: "Siti Peminjam",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			PasswordHash:      "$2a$10$WYmkE5HjRSYrJcUS9OGL/u9biq0iYc6GoUPiYLVd1UvO8hZPo98fO", // password123
		},
		{
			ID:                uuid.New(),
			FullName:          "Joko Wirausaha",
			IdentityNumber:    "3201234567890004",
			PhoneNumber:       "+6281234567896",
			Email:             "joko.wirausaha@gmail.com",
			Address:           "Jl. UMKM Sejahtera No. 321, Bekasi",
			DateOfBirth:       time.Date(1980, 7, 22, 0, 0, 0, 0, time.UTC),
			Occupation:        "Penjual Pakaian",
			MonthlyIncome:     3000000.00,
			BankAccountNumber: "5566778899",
			BankName:          "Bank BNI",
			AccountHolderName: "Joko Wirausaha",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			PasswordHash:      "$2a$10$WYmkE5HjRSYrJcUS9OGL/u9biq0iYc6GoUPiYLVd1UvO8hZPo98fO", // password123
		},
	}

	query := `
		INSERT INTO borrowers (
			id, full_name, identity_number, phone_number, email, address,
			date_of_birth, occupation, monthly_income, bank_account_number,
			bank_name, account_holder_name, created_at, updated_at, password_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (identity_number) DO NOTHING`

	for _, borrower := range borrowers {
		_, err := db.ExecContext(ctx, query,
			borrower.ID, borrower.FullName, borrower.IdentityNumber,
			borrower.PhoneNumber, borrower.Email, borrower.Address,
			borrower.DateOfBirth, borrower.Occupation, borrower.MonthlyIncome,
			borrower.BankAccountNumber, borrower.BankName, borrower.AccountHolderName,
			borrower.CreatedAt, borrower.UpdatedAt, borrower.PasswordHash,
		)
		if err != nil {
			return err
		}
		log.Printf(" Created borrower: %s (%s)", borrower.FullName, borrower.Occupation)
	}

	return nil
}
