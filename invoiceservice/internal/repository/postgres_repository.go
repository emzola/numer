package repository

import (
	"database/sql"

	"github.com/emzola/numer/invoiceservice/pkg/model"
)

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) Create(invoice model.Invoice) error {
	return nil
}
