package service

import (
	"context"
	"strconv"

	"github.com/emzola/numer/invoiceservice/pkg/model"
	"github.com/google/uuid"
)

// InvoiceRepository defines the interface for interacting with the invoice data store.
type InvoiceRepository interface {
	Create(invoice model.Invoice) error
}

type InvoiceService struct {
	repo              InvoiceRepository
	nextInvoiceNumber int
}

func NewInvoiceService(repo InvoiceRepository) *InvoiceService {
	return &InvoiceService{
		repo:              repo,
		nextInvoiceNumber: 100000, // Temporary for now
	}
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, invoice model.Invoice) (*model.Invoice, error) {
	// Generate invoice ID
	invoice.InvoiceID = uuid.New().String()

	// Calculate invoice subtotal, discount and total
	subtotal, discount, total := calculateInvoice(invoice.Items, invoice.DiscountPercentage)
	invoice.Subtotal = subtotal
	invoice.Discount = discount
	invoice.Total = total

	// Add invoice status
	invoice.Status = model.PENDING

	// TEMPORARY IMPLEMENTATION: Generate a new invoice number
	// TODO: concurrency safe invoice number increment
	invoiceNumber := strconv.Itoa(s.nextInvoiceNumber)
	s.nextInvoiceNumber++
	invoice.InvoiceNumber = invoiceNumber

	// Save to repository
	err := s.repo.Create(invoice)
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

// calculateInvoice calculates and returns the subtotal, discount and total of invoice items.
func calculateInvoice(items []model.InvoiceItem, discountPercentage float64) (float64, float64, float64) {
	var subtotal float64
	for _, item := range items {
		subtotal += float64(item.Quantity) * item.UnitPrice
	}
	discount := subtotal * (discountPercentage / 100)
	total := subtotal - discount
	return subtotal, discount, total
}
