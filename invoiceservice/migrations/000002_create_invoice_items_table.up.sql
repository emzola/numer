CREATE TABLE IF NOT EXISTS invoice_items (
    id SERIAL PRIMARY KEY,
    invoice_id UUID REFERENCES invoices(invoice_id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    quantity INT NOT NULL,
    unit_price BIGINT NOT NULL
);