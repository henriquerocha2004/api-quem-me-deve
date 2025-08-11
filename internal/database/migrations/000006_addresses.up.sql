CREATE TABLE addresses (
    id CHAR(26) PRIMARY KEY,
    street VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    zip_code VARCHAR(20) NOT NULL,
    neighborhood VARCHAR(100) NOT NULL,
    owner_id CHAR(26) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_addresses_owner_id ON addresses(owner_id);
CREATE INDEX idx_addresses_zip_code ON addresses(zip_code);