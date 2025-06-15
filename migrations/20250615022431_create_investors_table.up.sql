CREATE TABLE investors (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           full_name VARCHAR(100) NOT NULL,
                           identity_number VARCHAR(20) NOT NULL UNIQUE,
                           email VARCHAR(100) NOT NULL UNIQUE,
                           phone_number VARCHAR(15) NOT NULL,
                           address TEXT,
                           bank_account VARCHAR(50),
                           is_active BOOLEAN DEFAULT true,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_investors_identity_number ON investors(identity_number);
CREATE INDEX idx_investors_email ON investors(email);
CREATE INDEX idx_investors_is_active ON investors(is_active);