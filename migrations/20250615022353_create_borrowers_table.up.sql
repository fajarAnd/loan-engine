CREATE TABLE borrowers (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           full_name VARCHAR(100) NOT NULL,
                           identity_number VARCHAR(20) NOT NULL UNIQUE,
                           phone_number VARCHAR(15) NOT NULL,
                           email VARCHAR(100),
                           address TEXT,
                           date_of_birth DATE,
                           occupation VARCHAR(50),
                           monthly_income DECIMAL(15,2),
                           bank_account_number VARCHAR(50),
                           bank_name VARCHAR(100),
                           account_holder_name VARCHAR(100),
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_borrowers_identity_number ON borrowers(identity_number);
CREATE INDEX idx_borrowers_phone_number ON borrowers(phone_number);
CREATE INDEX idx_borrowers_email ON borrowers(email);