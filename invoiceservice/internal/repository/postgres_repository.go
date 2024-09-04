package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/emzola/numer/invoiceservice/pkg/model"
)

var (
	ErrInvoiceNumberIncrement = errors.New("failed to increment invoice number")
)

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) IncrementInvoiceNumber(ctx context.Context) (string, error) {
	var invoiceNumber int
	query := `
		UPDATE invoice_number_sequence 
		SET current_value = current_value + 1 
		RETURNING current_value
	`
	err := r.db.QueryRowContext(ctx, query).Scan(&invoiceNumber)
	if err != nil {
		return "", ErrInvoiceNumberIncrement
	}

	return fmt.Sprintf("%06d", invoiceNumber), nil
}

func (r *InvoiceRepository) Create(ctx context.Context, invoice model.Invoice) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
        INSERT INTO invoices (
            invoice_id, user_id, customer_id, invoice_number, status,
            issue_date, due_date, billing_currency, discount_percentage, 
            subtotal, discount, total, account_name, account_number, 
            bank_name, routing_number, note
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
        )`,
		invoice.InvoiceID, invoice.UserID, invoice.CustomerID, invoice.InvoiceNumber, invoice.Status,
		invoice.IssueDate, invoice.DueDate, invoice.BillingCurrency, invoice.DiscountPercentage,
		invoice.Subtotal, invoice.Discount, invoice.Total, invoice.PaymentInformation.AccountName,
		invoice.PaymentInformation.AccountNumber, invoice.PaymentInformation.BankName,
		invoice.PaymentInformation.RoutingNumber, invoice.Note,
	)
	if err != nil {
		return err
	}

	for _, item := range invoice.Items {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO invoice_items (
                invoice_id, description, quantity, unit_price
            ) VALUES (
                $1, $2, $3, $4
            )`,
			invoice.InvoiceID, item.Description, item.Quantity, item.UnitPrice,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
