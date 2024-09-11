package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/emzola/numer/invoice-service/internal/models"
	"github.com/shopspring/decimal"
)

var (
	ErrNotFound               = errors.New("not found")
	ErrInvalidRequest         = errors.New("the request is invalid")
	ErrInvoiceNumberIncrement = errors.New("failed to increment invoice number")
)

type invoiceRepository interface {
	CreateInvoice(ctx context.Context, invoice *models.Invoice) error
	GetInvoiceByID(ctx context.Context, invoiceID int64) (*models.Invoice, error)
	UpdateInvoice(ctx context.Context, invoice *models.Invoice) error
	ListInvoicesByUserID(ctx context.Context, userID int64, pageSize int, pageToken string) ([]*models.Invoice, string, error)
	GetDueInvoices(ctx context.Context, daysBeforeDue int32) ([]*models.Invoice, error)
	IncrementInvoiceNumber(ctx context.Context) (int64, error)
}

type InvoiceService struct {
	repo invoiceRepository
}

func NewInvoiceService(repo invoiceRepository) *InvoiceService {
	return &InvoiceService{repo: repo}
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	invoiceNumber, err := s.repo.IncrementInvoiceNumber(ctx)
	if err != nil {
		return nil, err
	}
	invoice.InvoiceNumber = fmt.Sprintf("%06d", invoiceNumber)
	invoice.Status = "draft"

	// Calculate invoice amounts
	subtotal, discount, total := calculateInvoiceAmounts(invoice.Items, invoice.DiscountPercentage) // e.g. 1000 discountPercentage == 10% discount
	invoice.Subtotal = subtotal
	invoice.DiscountAmount = discount
	invoice.Total = total

	err = s.repo.CreateInvoice(ctx, invoice)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *InvoiceService) GetInvoice(ctx context.Context, invoiceID int64) (*models.Invoice, error) {
	invoice, err := s.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	return invoice, nil
}

func (s *InvoiceService) UpdateInvoice(ctx context.Context, invoice *models.Invoice) error {
	// Recalculate invoice amounts
	subtotal, discount, total := calculateInvoiceAmounts(invoice.Items, invoice.DiscountPercentage)
	invoice.Subtotal = subtotal
	invoice.DiscountAmount = discount
	invoice.Total = total

	err := s.repo.UpdateInvoice(ctx, invoice)
	if err != nil {
		return err
	}
	return nil
}

func (s *InvoiceService) ListInvoicesByUserID(ctx context.Context, userID int64, pageSize int, pageToken string) ([]*models.Invoice, string, error) {
	invoices, nextPageToken, err := s.repo.ListInvoicesByUserID(ctx, userID, pageSize, pageToken)
	if err != nil {
		return nil, "", err
	}

	return invoices, nextPageToken, nil
}

func (s *InvoiceService) GetDueInvoices(ctx context.Context, threshold int32) ([]*models.Invoice, error) {
	invoices, err := s.repo.GetDueInvoices(ctx, threshold)
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

// ConvertDecimalToCents converts a decimal.Decimal to int64 (cents).
func ConvertDecimalToCents(d decimal.Decimal) int64 {
	cents := d.Mul(decimal.NewFromInt(100))
	return cents.IntPart()
}

// ConvertCentsToDecimal converts int64 (cents) to decimal.Decimal.
func ConvertCentsToDecimal(cents int64) decimal.Decimal {
	return decimal.NewFromInt(cents).Div(decimal.NewFromInt(100))
}

// ConvertPercentageToDecimal converts an int64 percentage (in hundredths) to decimal.Decimal.
func ConvertPercentageToDecimal(percentage int64) decimal.Decimal {
	return decimal.NewFromInt(percentage).Div(decimal.NewFromInt(10000))
}

// calculateInvoiceValues calculates the subtotal, discount, and total of invoice items.
func calculateInvoiceAmounts(items []*models.InvoiceItem, discountPercentage int64) (subtotal int64, discount int64, total int64) {
	subtotalDecimal := decimal.Zero
	for _, item := range items {
		itemTotal := decimal.NewFromInt(int64(item.Quantity)).Mul(ConvertCentsToDecimal(item.UnitPrice))
		subtotalDecimal = subtotalDecimal.Add(itemTotal)
	}

	discountDecimal := subtotalDecimal.Mul(ConvertPercentageToDecimal(discountPercentage))
	totalDecimal := subtotalDecimal.Sub(discountDecimal)

	subtotal = ConvertDecimalToCents(subtotalDecimal)
	discount = ConvertDecimalToCents(discountDecimal)
	total = ConvertDecimalToCents(totalDecimal)

	return
}
