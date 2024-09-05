package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"

	"github.com/emzola/numer/invoiceservice/pkg/model"
)

var (
	ErrInvoiceNumberIncrement = errors.New("failed to increment invoice number")
	ErrNotFound               = errors.New("not found")
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

func (r *InvoiceRepository) GetByID(ctx context.Context, invoiceID string) (*model.Invoice, error) {
	var invoice model.Invoice
	query := `SELECT invoice_id, user_id, customer_id, invoice_number, status, issue_date, due_date, billing_currency,
	                 discount_percentage, subtotal, discount, total, account_name, account_number, bank_name, routing_number, note
	          FROM invoices WHERE invoice_id = $1`

	err := r.db.QueryRowContext(ctx, query, invoiceID).Scan(
		&invoice.InvoiceID, &invoice.UserID, &invoice.CustomerID, &invoice.InvoiceNumber,
		&invoice.Status, &invoice.IssueDate, &invoice.DueDate, &invoice.BillingCurrency,
		&invoice.DiscountPercentage, &invoice.Subtotal, &invoice.Discount, &invoice.Total,
		&invoice.PaymentInformation.AccountName, &invoice.PaymentInformation.AccountNumber,
		&invoice.PaymentInformation.BankName, &invoice.PaymentInformation.RoutingNumber, &invoice.Note,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Query for invoice items
	itemsQuery := `SELECT description, quantity, unit_price FROM invoice_items WHERE invoice_id = $1`
	rows, err := r.db.QueryContext(ctx, itemsQuery, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Populate invoice items
	var items []model.InvoiceItem
	for rows.Next() {
		var item model.InvoiceItem
		if err := rows.Scan(&item.Description, &item.Quantity, &item.UnitPrice); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	invoice.Items = items

	return &invoice, nil
}

func (r *InvoiceRepository) ListByUserID(ctx context.Context, userID string, pageSize int, pageToken string) ([]*model.Invoice, string, error) {
	var invoices []*model.Invoice
	var offset int

	// Decode the pageToken into an offset (cursor-based pagination)
	if pageToken != "" {
		offset = decodePageToken(pageToken)
	}

	query := `SELECT invoice_id, user_id, customer_id, invoice_number, status, issue_date, due_date, billing_currency,
	                 discount_percentage, subtotal, discount, total, account_name, account_number, bank_name, routing_number, note
	          FROM invoices WHERE user_id = $1 ORDER BY issue_date DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	for rows.Next() {
		var invoice model.Invoice
		err := rows.Scan(
			&invoice.InvoiceID, &invoice.UserID, &invoice.CustomerID, &invoice.InvoiceNumber, &invoice.Status,
			&invoice.IssueDate, &invoice.DueDate, &invoice.BillingCurrency, &invoice.DiscountPercentage,
			&invoice.Subtotal, &invoice.Discount, &invoice.Total, &invoice.PaymentInformation.AccountName,
			&invoice.PaymentInformation.AccountNumber, &invoice.PaymentInformation.BankName,
			&invoice.PaymentInformation.RoutingNumber, &invoice.Note,
		)
		if err != nil {
			return nil, "", err
		}

		// Fetch invoice items for each invoice (if necessary)
		invoice.Items, err = r.fetchInvoiceItems(ctx, invoice.InvoiceID)
		if err != nil {
			return nil, "", err
		}

		invoices = append(invoices, &invoice)
	}

	// Generate nextPageToken if there are more rows
	var nextPageToken string
	if len(invoices) == pageSize {
		nextPageToken = encodePageToken(offset + pageSize)
	}

	return invoices, nextPageToken, nil
}

// fetchInvoiceItems fetches invoice items by invoice ID
func (r *InvoiceRepository) fetchInvoiceItems(ctx context.Context, invoiceID string) ([]model.InvoiceItem, error) {
	var items []model.InvoiceItem
	query := `SELECT description, quantity, unit_price FROM invoice_items WHERE invoice_id = $1`
	rows, err := r.db.QueryContext(ctx, query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.InvoiceItem
		err := rows.Scan(&item.Description, &item.Quantity, &item.UnitPrice)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// encodePageToken encodes the integer offset as a base64 pageToken.
func encodePageToken(offset int) string {
	// Convert the integer offset to a string
	offsetStr := strconv.Itoa(offset)

	// Encode the string in base64
	return base64.URLEncoding.EncodeToString([]byte(offsetStr))
}

// decodePageToken decodes the base64 pageToken to an integer offset.
func decodePageToken(pageToken string) int {
	// Decode the base64-encoded string
	offsetBytes, err := base64.URLEncoding.DecodeString(pageToken)
	if err != nil {
		return 0
	}

	// Convert the decoded bytes to an integer
	offset, err := strconv.Atoi(string(offsetBytes))
	if err != nil {
		return 0
	}

	return offset
}
