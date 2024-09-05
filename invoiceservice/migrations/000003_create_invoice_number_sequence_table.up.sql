CREATE TABLE IF NOT EXISTS invoice_number_sequence (
    id SERIAL PRIMARY KEY,
    current_value INT NOT NULL
);

-- Initialize with a starting value
INSERT INTO invoice_number_sequence (current_value) VALUES (100000);