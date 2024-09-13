package service_test

import (
	"context"
	"testing"

	"github.com/emzola/numer/invoice-service/internal/models"
	"github.com/emzola/numer/invoice-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for invoiceRepository
type MockInvoiceRepository struct {
	mock.Mock
}

func (m *MockInvoiceRepository) CreateInvoice(ctx context.Context, invoice *models.Invoice) error {
	args := m.Called(ctx, invoice)
	return args.Error(0)
}

func (m *MockInvoiceRepository) GetInvoiceByID(ctx context.Context, invoiceID int64) (*models.Invoice, error) {
	args := m.Called(ctx, invoiceID)
	return args.Get(0).(*models.Invoice), args.Error(1)
}

func (m *MockInvoiceRepository) UpdateInvoice(ctx context.Context, invoice *models.Invoice) error {
	args := m.Called(ctx, invoice)
	return args.Error(0)
}

func (m *MockInvoiceRepository) ListInvoicesByUserID(ctx context.Context, userID int64, pageSize int, pageToken string) ([]*models.Invoice, string, error) {
	args := m.Called(ctx, userID, pageSize, pageToken)
	return args.Get(0).([]*models.Invoice), args.String(1), args.Error(2)
}

func (m *MockInvoiceRepository) GetDueInvoices(ctx context.Context, daysBeforeDue int32) ([]*models.Invoice, error) {
	args := m.Called(ctx, daysBeforeDue)
	return args.Get(0).([]*models.Invoice), args.Error(1)
}

func (m *MockInvoiceRepository) IncrementInvoiceNumber(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestCreateInvoice(t *testing.T) {
	mockRepo := new(MockInvoiceRepository)
	svc := service.NewInvoiceService(mockRepo)

	invoice := &models.Invoice{
		UserID:             1,
		CustomerID:         2,
		Currency:           "USD",
		DiscountPercentage: 1000, // 10%
		Items: []*models.InvoiceItem{
			{
				Description: "Item 1",
				Quantity:    2,
				UnitPrice:   10000, // $100.00
			},
		},
	}

	expectedInvoice := *invoice
	expectedInvoice.InvoiceNumber = "000001"
	expectedInvoice.Status = "draft"
	expectedInvoice.Subtotal = 20000      // $200.00
	expectedInvoice.DiscountAmount = 2000 // 10% discount
	expectedInvoice.Total = 18000         // $180.00

	mockRepo.On("IncrementInvoiceNumber", mock.Anything).Return(int64(1), nil)
	mockRepo.On("CreateInvoice", mock.Anything, mock.Anything).Return(nil)

	createdInvoice, err := svc.CreateInvoice(context.Background(), invoice)

	assert.NoError(t, err)
	assert.Equal(t, &expectedInvoice, createdInvoice)
	mockRepo.AssertExpectations(t)
}

func TestGetInvoice(t *testing.T) {
	mockRepo := new(MockInvoiceRepository)
	svc := service.NewInvoiceService(mockRepo)

	expectedInvoice := &models.Invoice{ID: 1}
	mockRepo.On("GetInvoiceByID", mock.Anything, int64(1)).Return(expectedInvoice, nil)

	invoice, err := svc.GetInvoice(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedInvoice, invoice)
	mockRepo.AssertExpectations(t)
}

func TestUpdateInvoice(t *testing.T) {
	mockRepo := new(MockInvoiceRepository)
	svc := service.NewInvoiceService(mockRepo)

	invoice := &models.Invoice{
		ID:                 1,
		UserID:             1,
		CustomerID:         2,
		Currency:           "USD",
		DiscountPercentage: 1000, // 10%
		Items: []*models.InvoiceItem{
			{
				Description: "Item 1",
				Quantity:    2,
				UnitPrice:   10000, // $100.00
			},
		},
	}

	expectedInvoice := *invoice
	expectedInvoice.Subtotal = 20000      // $200.00
	expectedInvoice.DiscountAmount = 2000 // 10% discount
	expectedInvoice.Total = 18000         // $180.00

	mockRepo.On("UpdateInvoice", mock.Anything, mock.Anything).Return(nil)

	err := svc.UpdateInvoice(context.Background(), invoice)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestListInvoicesByUserID(t *testing.T) {
	mockRepo := new(MockInvoiceRepository)
	svc := service.NewInvoiceService(mockRepo)

	invoices := []*models.Invoice{
		{ID: 1},
		{ID: 2},
	}
	mockRepo.On("ListInvoicesByUserID", mock.Anything, int64(1), 10, "").Return(invoices, "", nil)

	result, nextPageToken, err := svc.ListInvoicesByUserID(context.Background(), 1, 10, "")

	assert.NoError(t, err)
	assert.Equal(t, invoices, result)
	assert.Empty(t, nextPageToken)
	mockRepo.AssertExpectations(t)
}
