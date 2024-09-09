-- +goose Up
CREATE TABLE IF NOT EXISTS invoice_number_sequence (
    id SERIAL PRIMARY KEY,
    current_value INT NOT NULL
);

INSERT INTO invoice_number_sequence (current_value) VALUES (100000);

-- +goose Down
DROP TABLE IF EXISTS invoice_number_sequence;
