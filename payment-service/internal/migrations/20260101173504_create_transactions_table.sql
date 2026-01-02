-- +goose Up
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    reference TEXT NOT NULL UNIQUE,
    amount BIGINT NOT NULL,
    currency TEXT NOT NULL CHECK (currency IN ('ETB', 'USD')),
    status INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_transactions_status ON transactions(status);
-- +goose Down

DROP TABLE transactions;

DROP INDEX IF EXISTS idx_transactions_status;

CREATE INDEX idx_transactions_status ON transactions(status);
