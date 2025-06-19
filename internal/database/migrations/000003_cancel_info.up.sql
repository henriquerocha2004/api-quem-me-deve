CREATE TABLE IF NOT EXISTS cancel_info (
    id CHAR(26) PRIMARY KEY,
    reason TEXT NOT NULL,
    cancel_date TIMESTAMP,
    cancelled_by CHAR(26) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    debt_id CHAR(26) NOT NULL REFERENCES debts(id)
);

CREATE INDEX idx_cancel_info_cancel_date ON cancel_info(cancel_date);
CREATE INDEX idx_cancel_info_cancelled_by ON cancel_info(cancelled_by);