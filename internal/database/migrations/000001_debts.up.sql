CREATE TABLE debts (
    id CHAR(26) PRIMARY KEY,
    description TEXT NOT NULL,
    total_value DECIMAL(12,2) NOT NULL,
    due_date TIMESTAMP,
    installments_quantity INTEGER NOT NULL,
    debt_date TIMESTAMP,
    status VARCHAR(255) NOT NULL,
    user_client_id CHAR(26) NOT NULL,
    product_ids CHAR(26)[] DEFAULT '{}',
    service_ids CHAR(26)[] DEFAULT '{}',
    finished_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_debts_status ON debts(status);
CREATE INDEX idx_debts_user_client_id ON debts(user_client_id);
CREATE INDEX idx_debts_debt_date ON debts(debt_date);
CREATE INDEX idx_debts_finished_at ON debts(finished_at);