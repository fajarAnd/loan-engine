CREATE TABLE loan_state_histories (
                                      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                      loan_id UUID NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
                                      changed_by_employee_id UUID REFERENCES employees(id) ON DELETE SET NULL,
                                      previous_state VARCHAR(50),
                                      new_state VARCHAR(50) NOT NULL,
                                      change_reason TEXT,
                                      changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_loan_state_histories_loan_id ON loan_state_histories(loan_id);
CREATE INDEX idx_loan_state_histories_changed_by ON loan_state_histories(changed_by_employee_id);
CREATE INDEX idx_loan_state_histories_changed_at ON loan_state_histories(changed_at);