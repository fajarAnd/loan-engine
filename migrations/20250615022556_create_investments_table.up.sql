CREATE TABLE investments (
                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                             loan_id UUID NOT NULL REFERENCES loans(id) ON DELETE RESTRICT,
                             investor_id UUID NOT NULL REFERENCES investors(id) ON DELETE RESTRICT,
                             investment_amount DECIMAL(15,2) NOT NULL CHECK (investment_amount > 0),
                             expected_return DECIMAL(15,2) NOT NULL CHECK (expected_return >= 0),
                             investment_date DATE NOT NULL DEFAULT CURRENT_DATE,
                             agreement_url TEXT,
                             agreement_signed BOOLEAN DEFAULT false,
                             agreement_signed_date DATE,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                             UNIQUE(loan_id, investor_id)
);

CREATE INDEX idx_investments_loan_id ON investments(loan_id);
CREATE INDEX idx_investments_investor_id ON investments(investor_id);
CREATE INDEX idx_investments_investment_date ON investments(investment_date);
CREATE INDEX idx_investments_agreement_signed ON investments(agreement_signed);