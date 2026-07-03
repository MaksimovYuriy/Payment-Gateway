-- +goose Up
CREATE TABLE payment_attempts (
    id SERIAL PRIMARY KEY,
    payment_id INT NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    bank_id INT NOT NULL REFERENCES banks(id) ON DELETE RESTRICT,
    status VARCHAR(20) NOT NULL DEFAULT 'created',
    external_payment_id VARCHAR(255),
    error_message TEXT,
    error_code VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_payment_attempts_status 
    CHECK (status IN ('created', 'processing', 'succeeded', 'failed', 'timeout', 'cancelled'))
);

-- +goose Down
DROP TABLE IF EXISTS payment_attempts;
