# Loan Engine

## Features
> For detail documentation workflow and features please refer to this page: **[Summary Requirements](doc/summary-requirements.md)**

- **Loan Management**
- **Multi-role Authentication**
- **Make investment**
- **PDF Generation**
- **Database Schema Management**

## Tech Stack
> For detail Technical documentation please refer to this page: **[Technical Doc](doc/tech-doc.md)**


- **Backend**: Go (Golang) with Chi Router
- **Database**: PostgreSQL with PGX driver
- **Authentication**: JWT tokens
- **PDF Generation**: gofpdf library
- **Migrations**: golang-migrate
- **Configuration**: Viper
- **Logging**: Zerolog

---

## Prerequisites

- Go 1.24.2 or higher
- PostgreSQL 13+
- Make

---

## How to Run

### 1. Clone Repository
```bash
git clone <repository-url>
```

### 2. Setup Configuration
Copy the example configuration and modify as needed:
```bash
cp config.example.yml config.yml
```

Update `config.yml` with your database credentials:
```yaml
database:
  host: localhost
  port: 5432
  user: your_db_user
  password: your_db_password
  name: loan_engine
  schema: public
  sslmode: disable
  max_open_conns: 25
  max_idle_conns: 10
  max_lifetime: 5m
  dsn: "postgres://your_db_user:your_db_password@localhost:5432/loan_engine"

jwt:
  secret: "your-super-secret-jwt-key-here"
```

### 3. Install Dependencies
```bash
go mod download
```

### 4. Database Setup

#### Install Migration Tool (one-time setup)
```bash
make install-migrate
```

#### Run Migrations
```bash
make migrate-up
```

#### Seed Database with Sample Data
```bash
make seed
```

### 5. Build Application
```bash
go build -o loan-engine main.go
```

### 6. Run Application
directly with Go:
```bash
make run
```

#### Health Check
Verify the application is running:
```bash
curl http://localhost:8080/__health
```

The server will start on `http://localhost:8080`


### 7. Run Unit Test
```bash
go run main.go
```


---

## Sample Users (From Seeding)

The seeding process creates sample users for testing:

### Employees
| Role | Email | Password | Full Name |
|------|-------|----------|-----------|
| FIELD_VALIDATOR | `validator@amartha.com` | `password123` | Ahmad Validator |
| FIELD_OFFICER | `officer@amartha.com` | `password123` | Sari Fieldofficer |
| ADMIN | `admin@amartha.com` | `password123` | Budi Administrator |

### Investors
| Email | Password | Full Name |
|-------|----------|-----------|
| `rina.investor@gmail.com` | `password123` | Rina Investor |
| `doni.kapital@gmail.com` | `password123` | Doni Kapital |

### Borrowers
| Email | Password | Full Name | Occupation |
|-------|----------|-----------|------------|
| `siti.peminjam@gmail.com` | `password123` | Siti Peminjam | Penjual Makanan |
| `joko.wirausaha@gmail.com` | `password123` | Joko Wirausaha | Penjual Pakaian |

---

## How to Test

### 1. API Testing
Use the provided HTTP test files in the `api/` directory:

- **Authentication**: `api/auth.http`
- **Loan Management**: `api/loan.http`
- **File Upload**: `api/file.http`
- **Loan Approval**: `api/approve_loan.http`
- **Investment**: `api/investment.http`
- **Disbursement**: `api/disburse_loan.http`

### 2. Complete E2E Workflow Test

```mermaid
flowchart LR
    LoginBorrower[1. Login Borrower] --> CreateLoan[2. Create Loan]
    CreateLoan --> LoginValidator[3. Login Validator]
    LoginValidator --> UploadSurvey[4. Upload Survey]
    UploadSurvey --> LoginOfficer[5. Login Officer]
    LoginOfficer --> ApproveLoan[6. Approve Loan]
    ApproveLoan --> LoginInvestor[7. Login Investor]
    LoginInvestor --> MakeInvestment[8. Make Investment]
    MakeInvestment --> DisburseLoan[9. Disburse Loan]
```

---

## What's Next

### Planned Features

#### 1. **Loan Portfolio Management**
- [ ] Get available loans for investment (`GET /api/v1/loans/available`)
- [ ] Get investor portfolio (`GET /api/v1/investors/{id}/portfolio`)
- [ ] Loan status tracking and history

#### 2. **API Enhancements**
- [ ] Comprehensive API documentation with Swagger

### Technical Debt & Improvements
- [ ] Add & comprehensive unit tests
- [ ] Optimize database queries
- [ ] Implement caching layer
- [ ] Implement graceful shutdown

---
