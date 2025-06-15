CREATE TYPE loan_state_enum AS ENUM (
    'PROPOSED',
    'APPROVED',
    'INVESTED',
    'DISBURSED',
    'REJECTED'
    );

CREATE TABLE loans (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       borrower_id UUID NOT NULL REFERENCES borrowers(id) ON DELETE RESTRICT,
                       principal_amount DECIMAL(15,2) NOT NULL CHECK (principal_amount > 0),
                       interest_rate DECIMAL(5,2) NOT NULL CHECK (interest_rate >= 0),
                       roi_rate DECIMAL(5,2) NOT NULL CHECK (roi_rate >= 0),
                       loan_term_month INTEGER NOT NULL CHECK (loan_term_month > 0),
                       current_state loan_state_enum NOT NULL DEFAULT 'PROPOSED',
                       loan_agreement_pdf_url TEXT,

                       field_validator_employee_id UUID REFERENCES employees(id) ON DELETE SET NULL,
                       survey_date DATE,
                       field_visit_proof_url TEXT,
                       survey_notes TEXT,
                       approving_employee_id UUID REFERENCES employees(id) ON DELETE SET NULL,
                       approval_date DATE,
                       approval_notes TEXT,

    -- Disbursement fields
                       field_officer_employee_id UUID REFERENCES employees(id) ON DELETE SET NULL,
                       disbursement_date DATE,
                       signed_agreement_url TEXT,
                       disbursement_notes TEXT,

                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_loans_borrower_id ON loans(borrower_id);
CREATE INDEX idx_loans_current_state ON loans(current_state);
CREATE INDEX idx_loans_field_validator_employee_id ON loans(field_validator_employee_id);
CREATE INDEX idx_loans_field_officer_employee_id ON loans(field_officer_employee_id);
CREATE INDEX idx_loans_created_at ON loans(created_at);