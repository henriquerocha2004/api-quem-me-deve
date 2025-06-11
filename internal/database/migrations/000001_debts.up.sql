CREATE TABLE debts (
    id CHAR(26) PRIMARY KEY,
    description TEXT NOT NULL,
    total_value NUMERIC(12,2) NOT NULL,
    due_date TIMESTAMP,
    installments_quantity INTEGER NOT NULL,
    debt_date TIMESTAMP,
    status TEXT NOT NULL,
    user_client_id CHAR(26) NOT NULL,
    product_ids CHAR(26)[] DEFAULT '{}',
    service_ids CHAR(26)[] DEFAULT '{}',
    finished_at TIMESTAMP
);