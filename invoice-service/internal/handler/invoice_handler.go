package handler

import (
	"context"
	"time"

	"github.com/emzola/numer/invoiceservice/internal/models"
	"github.com/emzola/numer/invoiceservice/internal/service"
	pb "github.com/emzola/numer/invoiceservice/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InvoiceHandler struct {
	service *service.InvoiceService
	pb.UnimplementedInvoiceServiceServer
}

func NewInvoiceHandler(service *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{service: service}
}

func (h *InvoiceHandler) CreateInvoice(ctx context.Context, req *pb.CreateInvoiceRequest) (*pb.CreateInvoiceResponse, error) {
	invoice := &models.Invoice{
		UserID:             req.UserId,
		CustomerID:         req.CustomerId,
		IssueDate:          req.IssueDate.AsTime(),
		DueDate:            req.DueDate.AsTime(),
		Currency:           req.Currency,
		DiscountPercentage: req.DiscountPercentage,
		AccountName:        req.AccountName,
		AccountNumber:      req.AccountNumber,
		BankName:           req.BankName,
		RoutingNumber:      req.RoutingNumber,
		Note:               req.Note,
	}

	// Add invoice items
	for _, item := range req.Items {
		invoice.Items = append(invoice.Items, &models.InvoiceItem{
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	invoice, err := h.service.CreateInvoice(ctx, invoice)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create invoice")
	}

	return &pb.CreateInvoiceResponse{InvoiceId: invoice.ID}, nil
}

func (h *InvoiceHandler) GetInvoice(ctx context.Context, req *pb.GetInvoiceRequest) (*pb.GetInvoiceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	invoice, err := h.service.GetInvoice(ctx, req.InvoiceId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "invoice not found")
	}

	return &pb.GetInvoiceResponse{Invoice: models.ConvertInvoiceToProto(invoice)}, nil
}

func (h *InvoiceHandler) UpdateInvoice(ctx context.Context, req *pb.UpdateInvoiceRequest) (*pb.UpdateInvoiceResponse, error) {
	invoice, err := h.service.GetInvoice(ctx, req.InvoiceId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "invoice not found")
	}

	invoice.Status = req.Status
	invoice.IssueDate = req.IssueDate.AsTime()
	invoice.DueDate = req.DueDate.AsTime()
	invoice.Currency = req.Currency
	invoice.DiscountPercentage = req.DiscountPercentage
	invoice.AccountName = req.AccountName
	invoice.AccountName = req.AccountNumber
	invoice.BankName = req.BankName
	invoice.RoutingNumber = req.RoutingNumber
	invoice.Note = req.Note

	// Update invoice items
	invoice.Items = []*models.InvoiceItem{}
	for _, itemReq := range req.Items {
		invoice.Items = append(invoice.Items, &models.InvoiceItem{
			Description: itemReq.Description,
			Quantity:    itemReq.Quantity,
			UnitPrice:   itemReq.UnitPrice,
		})
	}

	// Call service to update invoice
	err = h.service.UpdateInvoice(ctx, invoice)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update invoice")
	}

	return &pb.UpdateInvoiceResponse{Invoice: models.ConvertInvoiceToProto(invoice)}, nil
}

func (s *InvoiceHandler) ListInvoices(ctx context.Context, req *pb.ListInvoicesRequest) (*pb.ListInvoicesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	invoices, nextPageToken, err := s.service.ListInvoicesByUserID(ctx, req.UserId, int(req.PageSize), req.PageToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list invoices")
	}

	protoInvoices := make([]*pb.Invoice, len(invoices))
	for i, invoice := range invoices {
		protoInvoices[i] = models.ConvertInvoiceToProto(invoice)
	}

	return &pb.ListInvoicesResponse{
		Invoices:      protoInvoices,
		NextPageToken: nextPageToken,
	}, nil
}
