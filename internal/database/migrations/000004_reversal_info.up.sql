CREATE TABLE IF NOT EXISTS reversal_info (
    id CHAR(26) PRIMARY KEY,
    reason TEXT NOT NULL,
    reversal_date TIMESTAMP,
    reversed_by CHAR(26) NOT NULL,
    reversed_installment_qtd INTEGER NOT NULL DEFAULT 0,
    cancelled_installment_qtd INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para otimizar consultas
CREATE INDEX idx_reversal_info_reversal_date ON reversal_info(reversal_date);
CREATE INDEX idx_reversal_info_reversed_by ON reversal_info(reversed_by);