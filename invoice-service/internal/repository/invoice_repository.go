package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"strconv"

	"github.com/emzola/numer/invoice-service/internal/models"
)

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) IncrementInvoiceNumber(ctx context.Context) (int64, error) {
	query := `
		UPDATE invoice_number_sequence 
		SET current_value = current_value + 1 
		RETURNING current_value
	`
	var invoiceNumber int64
	err := r.db.QueryRowContext(ctx, query).Scan(&invoiceNumber)
	if err != nil {
		return 0, err
	}

	return invoiceNumber, nil
}

func (r *InvoiceRepository) CreateInvoice(ctx context.Context, invoice *models.Invoice) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert invoice details
	query := `
		INSERT INTO invoices (user_id, customer_id, invoice_number, status,	issue_date, due_date, currency, subtotal, 
			discount_percentage, discount_amount, total, account_name, account_number, bank_name, routing_number, note)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query,
		invoice.UserID, invoice.CustomerID, invoice.InvoiceNumber, invoice.Status, invoice.IssueDate, invoice.DueDate,
		invoice.Currency, invoice.Subtotal, invoice.DiscountPercentage, invoice.DiscountAmount, invoice.Total,
		invoice.AccountName, invoice.AccountNumber, invoice.BankName, invoice.RoutingNumber, invoice.Note).Scan(
		&invoice.ID, &invoice.CreatedAt, &invoice.UpdatedAt)
	if err != nil {
		return err
	}

	// Insert invoice items
	for _, item := range invoice.Items {
		itemQuery := `
			INSERT INTO invoice_items (invoice_id, description, quantity, unit_price)
			VALUES ($1, $2, $3, $4)
			RETURNING id`
		err := tx.QueryRowContext(ctx, itemQuery, invoice.ID, item.Description, item.Quantity, item.UnitPrice).Scan(&item.ID)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *InvoiceRepository) GetInvoiceByID(ctx context.Context, invoiceID int64) (*models.Invoice, error) {
	// Fetch invoice
	query := `
		SELECT id, user_id, customer_id, invoice_number, status, issue_date, due_date, currency, subtotal, 
			discount_percentage, discount_amount, total, account_name, account_number, bank_name, routing_number, note, 
			created_at, updated_at
		FROM invoices
		WHERE id = $1`

	var invoice models.Invoice

	err := r.db.QueryRowContext(ctx, query, invoiceID).Scan(
		&invoice.ID, &invoice.UserID, &invoice.CustomerID, &invoice.InvoiceNumber, &invoice.Status, &invoice.IssueDate,
		&invoice.DueDate, &invoice.Currency, &invoice.Subtotal, &invoice.DiscountPercentage, &invoice.DiscountAmount, &invoice.Total,
		&invoice.AccountName, &invoice.AccountNumber, &invoice.BankName, &invoice.RoutingNumber, &invoice.Note,
		&invoice.CreatedAt, &invoice.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Fetch associated invoice items
	itemQuery := `
		SELECT id, description, quantity, unit_price
		FROM invoice_items 
		WHERE invoice_id = $1`
	rows, err := r.db.QueryContext(ctx, itemQuery, invoice.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.InvoiceItem
	for rows.Next() {
		var item models.InvoiceItem
		if err := rows.Scan(&item.ID, &item.Description, &item.Quantity, &item.UnitPrice); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	invoice.Items = items

	return &invoice, nil
}

func (r *InvoiceRepository) UpdateInvoice(ctx context.Context, invoice *models.Invoice) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update invoice details
	updateInvoiceQuery := `
		UPDATE invoices 
		SET status = $1, issue_date = $2, due_date = $3, currency = $4, discount_percentage = $5, account_name = $6, 
		account_number = $7, bank_name = $8, routing_number = $9, note = $10, updated_at = NOW()
		WHERE id = $11`
	_, err = tx.ExecContext(ctx, updateInvoiceQuery,
		invoice.Status, invoice.IssueDate, invoice.DueDate, invoice.Currency, invoice.DiscountPercentage,
		invoice.AccountName, invoice.AccountNumber, invoice.BankName, invoice.RoutingNumber, invoice.Note, invoice.ID)
	if err != nil {
		return err
	}

	// Delete old invoice items
	deleteItemsQuery := `DELETE FROM invoice_items WHERE invoice_id = $1`
	_, err = tx.ExecContext(ctx, deleteItemsQuery, invoice.ID)
	if err != nil {
		return err
	}

	// Insert updated invoice items
	for _, item := range invoice.Items {
		insertItemQuery := `
			INSERT INTO invoice_items (invoice_id, description, quantity, unit_price)
			VALUES ($1, $2, $3, $4)`
		_, err := tx.ExecContext(ctx, insertItemQuery, invoice.ID, item.Description, item.Quantity, item.UnitPrice)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *InvoiceRepository) ListInvoicesByUserID(ctx context.Context, userID int64, pageSize int, pageToken string) ([]*models.Invoice, string, error) {
	var invoices []*models.Invoice
	var offset int

	// Decode the pageToken into an offset (cursor-based pagination)
	if pageToken != "" {
		offset = decodePageToken(pageToken)
	}

	query := `
		SELECT id, user_id, customer_id, invoice_number, status, issue_date, due_date, currency, subtotal, 
			discount_percentage, discount_amount, total, account_name, account_number, bank_name, routing_number, note, 
			created_at, updated_at
	    FROM invoices 
		WHERE user_id = $1 
		ORDER BY issue_date DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	for rows.Next() {
		var invoice models.Invoice
		err := rows.Scan(
			&invoice.ID, &invoice.UserID, &invoice.CustomerID, &invoice.InvoiceNumber, &invoice.Status, &invoice.IssueDate,
			&invoice.DueDate, &invoice.Currency, &invoice.Subtotal, &invoice.DiscountPercentage, &invoice.DiscountAmount, &invoice.Total,
			&invoice.AccountName, &invoice.AccountNumber, &invoice.BankName, &invoice.RoutingNumber, &invoice.Note, &invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			return nil, "", err
		}

		// Fetch invoice items for each invoice
		invoice.Items, err = r.fetchInvoiceItems(ctx, invoice.ID)
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

func (r *InvoiceRepository) fetchInvoiceItems(ctx context.Context, invoiceID int64) ([]*models.InvoiceItem, error) {
	var items []*models.InvoiceItem
	query := `
		SELECT id, description, quantity, unit_price 
		FROM invoice_items 
		WHERE invoice_id = $1`
	rows, err := r.db.QueryContext(ctx, query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.InvoiceItem
		err := rows.Scan(&item.ID, &item.Description, &item.Quantity, &item.UnitPrice)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
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
