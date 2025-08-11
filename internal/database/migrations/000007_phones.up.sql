CREATE TABLE phones (
    id CHAR(26) PRIMARY KEY,
    number VARCHAR(20) NOT NULL,
    description VARCHAR(255),
    owner_id CHAR(26) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_phones_owner_id ON phones(owner_id);
CREATE INDEX idx_phones_number ON phones(number);