-- +goose Up
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    merchant_id INT NOT NULL REFERENCES merchants(id) ON DELETE RESTRICT,
    bank_id INT REFERENCES banks(id) ON DELETE RESTRICT,
    order_id VARCHAR(255) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    currency CHAR(3) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'created',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (merchant_id, order_id),
    CONSTRAINT chk_status CHECK (status IN ('created', 'processing', 'completed', 'failed', 'cancelled'))
);

-- +goose Down
DROP TABLE IF EXISTS payments;
