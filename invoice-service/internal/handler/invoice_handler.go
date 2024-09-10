package handler

import (
	"context"

	"github.com/emzola/numer/invoice-service/internal/models"
	"github.com/emzola/numer/invoice-service/internal/service"
	publisher "github.com/emzola/numer/invoice-service/internal/service/rabbitmq"
	pb "github.com/emzola/numer/invoice-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InvoiceHandler struct {
	service   *service.InvoiceService
	publisher *publisher.Publisher
	pb.UnimplementedInvoiceServiceServer
}

func NewInvoiceHandler(service *service.InvoiceService, publisher *publisher.Publisher) *InvoiceHandler {
	return &InvoiceHandler{
		service:   service,
		publisher: publisher,
	}
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

	invoice, err := h.service.CreateInvoice(ctx, invoice)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Publish activity to rabbitMQ
	message := map[string]interface{}{
		"invoice_id":  invoice.ID,
		"user_id":     invoice.UserID,
		"action":      "Invoice creation",
		"description": "Created invoice" + " " + invoice.InvoiceNumber,
	}

	h.publisher.Publish(message)

	return &pb.CreateInvoiceResponse{InvoiceId: invoice.ID}, nil
}

func (h *InvoiceHandler) GetInvoice(ctx context.Context, req *pb.GetInvoiceRequest) (*pb.GetInvoiceResponse, error) {
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateInvoiceResponse{
		InvoiceId: invoice.ID,
		Message:   "invoice successfully updated",
	}, nil
}

func (s *InvoiceHandler) ListInvoices(ctx context.Context, req *pb.ListInvoicesRequest) (*pb.ListInvoicesResponse, error) {
	invoices, nextPageToken, err := s.service.ListInvoicesByUserID(ctx, req.UserId, int(req.PageSize), req.PageToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
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

func (s *InvoiceHandler) GetDueInvoices(ctx context.Context, req *pb.GetDueInvoicesRequest) (*pb.GetDueInvoicesResponse, error) {
	invoices, err := s.service.GetDueInvoices(ctx, req.Threshold)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	protoInvoices := make([]*pb.Invoice, len(invoices))
	for i, invoice := range invoices {
		protoInvoices[i] = models.ConvertInvoiceToProto(invoice)
	}

	return &pb.GetDueInvoicesResponse{Invoices: protoInvoices}, nil
}
