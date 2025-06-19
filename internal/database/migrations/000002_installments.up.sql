CREATE TABLE IF NOT EXISTS installments (
    id CHAR(26) PRIMARY KEY,
    description TEXT NOT NULL,
    value DECIMAL(12,2) NOT NULL,
    due_date TIMESTAMP,
    deb_date TIMESTAMP,
    status VARCHAR(255) NOT NULL,
    payment_date TIMESTAMP,
    payment_method VARCHAR(255),
    number INTEGER NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    debt_id CHAR(26) NOT NULL REFERENCES debts(id)
);

CREATE INDEX idx_installments_number ON installments(number);
CREATE INDEX idx_installments_deb_date ON installments(deb_date);
CREATE INDEX idx_installments_status ON installments(status);