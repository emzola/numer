-- +goose Up
CREATE TABLE IF NOT EXISTS activities (
    id SERIAL PRIMARY KEY,
    invoice_id BIGINT,
    user_id BIGINT,
    action VARCHAR(255),
    description VARCHAR(255),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS activities;
