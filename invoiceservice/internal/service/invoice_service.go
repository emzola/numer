package service

import (
	"context"
	"errors"

	"github.com/emzola/numer/invoiceservice/internal/repository"
	"github.com/emzola/numer/invoiceservice/pkg/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrNotFound               = errors.New("not found")
	ErrInvalidRequest         = errors.New("the request is invalid")
	ErrInvoiceNumberIncrement = errors.New("failed to increment invoice number")
)

// InvoiceRepository defines the interface for interacting with the invoice data store.
type InvoiceRepository interface {
	IncrementInvoiceNumber(ctx context.Context) (string, error)
	Create(ctx context.Context, invoice model.Invoice) error
	GetByID(ctx context.Context, invoiceID string) (*model.Invoice, error)
	ListByUserID(ctx context.Context, userID string, pageSize int, pageToken string) ([]*model.Invoice, string, error)
}

type InvoiceService struct {
	repo InvoiceRepository
}

func NewInvoiceService(repo InvoiceRepository) *InvoiceService {
	return &InvoiceService{repo: repo}
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, invoice model.Invoice) (*model.Invoice, error) {
	// Generate invoice ID
	invoice.InvoiceID = uuid.New().String()

	// Calculate invoice subtotal, discount and total
	subtotal, discount, total := calculateInvoiceValues(invoice.Items, invoice.DiscountPercentage) // e.g. 1000 discountPercentage == 10% discount
	invoice.Subtotal = subtotal
	invoice.Discount = discount
	invoice.Total = total

	// Add invoice status
	invoice.Status = model.PENDING

	// Generate a new invoice number
	invoiceNumber, err := s.repo.IncrementInvoiceNumber(ctx)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInvoiceNumberIncrement):
			return nil, ErrInvoiceNumberIncrement
		default:
			return nil, err
		}
	}
	invoice.InvoiceNumber = invoiceNumber

	// Save to repository
	err = s.repo.Create(ctx, invoice)
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (s *InvoiceService) GetInvoice(ctx context.Context, invoiceID string) (*model.Invoice, error) {
	// Call repository to fetch the invoice by its ID
	invoice, err := s.repo.GetByID(ctx, invoiceID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return invoice, nil
}

func (s *InvoiceService) ListInvoices(ctx context.Context, userID string, pageSize int, pageToken string) ([]*model.Invoice, string, error) {
	// Call repository to fetch paginated invoices
	invoices, nextPageToken, err := s.repo.ListByUserID(ctx, userID, pageSize, pageToken)
	if err != nil {
		return nil, "", err
	}

	return invoices, nextPageToken, nil
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
	return decimal.NewFromInt(percentage).Div(decimal.NewFromInt(100))
}

// Example function to calculate the subtotal, discount, and total.
func calculateInvoiceValues(items []model.InvoiceItem, discountPercentage int64) (subtotal int64, discount int64, total int64) {
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
