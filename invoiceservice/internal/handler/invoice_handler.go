package handler

import (
	"context"
	"errors"
	"time"

	invoicepb "github.com/emzola/numer/invoiceservice/genproto"
	"github.com/emzola/numer/invoiceservice/internal/service"
	"github.com/emzola/numer/invoiceservice/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InvoiceServiceServer is the gRPC server for the Invoice service.
type InvoiceServiceServer struct {
	invoicepb.UnimplementedInvoiceServiceServer
	service *service.InvoiceService
}

func NewInvoiceServiceServer(service *service.InvoiceService) *InvoiceServiceServer {
	return &InvoiceServiceServer{service: service}
}

func (s *InvoiceServiceServer) CreateInvoice(ctx context.Context, req *invoicepb.CreateInvoiceRequest) (*invoicepb.CreateInvoiceResponse, error) {
	if req == nil || req.UserId == "" || req.CustomerId == "" {
		return nil, status.Errorf(codes.InvalidArgument, service.ErrInvalidRequest.Error())
	}

	invoice := model.Invoice{
		UserID:             req.UserId,
		CustomerID:         req.CustomerId,
		IssueDate:          req.IssueDate,
		DueDate:            req.DueDate,
		BillingCurrency:    req.BillingCurrency,
		Items:              model.ConvertProtoItemsToInvoiceItems(req.Items),
		DiscountPercentage: req.DiscountPercentage,
		PaymentInformation: model.PaymentInformation{
			AccountName:   req.PaymentInformation.AccountName,
			AccountNumber: req.PaymentInformation.AccountNumber,
			BankName:      req.PaymentInformation.BankName,
			RoutingNumber: req.PaymentInformation.RoutingNumber,
		},
		Note: req.Note,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	createdInvoice, err := s.service.CreateInvoice(ctx, invoice)
	if err != nil {
		code, errMsg := mapToGRPCErrorCode(err), err.Error()
		return nil, status.Errorf(code, errMsg)
	}

	return &invoicepb.CreateInvoiceResponse{
		InvoiceId:     createdInvoice.InvoiceID,
		InvoiceNumber: createdInvoice.InvoiceNumber,
	}, nil
}

// mapToGRPCErrorCode maps domain-specific errors to gRPC status codes.
func mapToGRPCErrorCode(err error) codes.Code {
	switch {
	case errors.Is(err, service.ErrNotFound):
		return codes.NotFound
	case errors.Is(err, service.ErrInvalidRequest):
		return codes.InvalidArgument
	default:
		return codes.Internal
	}
}

func (s *InvoiceServiceServer) GetInvoice(ctx context.Context, req *invoicepb.GetInvoiceRequest) (*invoicepb.GetInvoiceResponse, error) {
	if req == nil || req.InvoiceId == "" {
		return nil, status.Errorf(codes.InvalidArgument, service.ErrInvalidRequest.Error())
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	invoice, err := s.service.GetInvoice(ctx, req.InvoiceId)
	if err != nil {
		code, errMsg := mapToGRPCErrorCode(err), err.Error()
		return nil, status.Errorf(code, errMsg)
	}

	response := &invoicepb.GetInvoiceResponse{
		Invoice: model.ConvertInvoiceToProto(invoice),
	}

	return response, nil
}

func (s *InvoiceServiceServer) ListInvoices(ctx context.Context, req *invoicepb.ListInvoicesRequest) (*invoicepb.ListInvoicesResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, service.ErrInvalidRequest.Error())
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	invoices, nextPageToken, err := s.service.ListInvoices(ctx, req.UserId, int(req.PageSize), req.PageToken)
	if err != nil {
		code, errMsg := mapToGRPCErrorCode(err), err.Error()
		return nil, status.Errorf(code, errMsg)
	}

	protoInvoices := make([]*invoicepb.Invoice, len(invoices))
	for i, invoice := range invoices {
		protoInvoices[i] = model.ConvertInvoiceToProto(invoice)
	}

	return &invoicepb.ListInvoicesResponse{
		Invoices:      protoInvoices,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *InvoiceServiceServer) UpdateInvoiceStatus(ctx context.Context, req *invoicepb.UpdateInvoiceStatusRequest) (*invoicepb.UpdateInvoiceStatusResponse, error) {
	if req == nil || req.InvoiceId == "" {
		return nil, status.Errorf(codes.InvalidArgument, service.ErrInvalidRequest.Error())
	}

	invoiceStatus := model.InvoiceStatus(req.Status)

	err := s.service.UpdateInvoiceStatus(ctx, req.InvoiceId, invoiceStatus)
	if err != nil {
		code, errMsg := mapToGRPCErrorCode(err), err.Error()
		return nil, status.Errorf(code, errMsg)
	}

	return &invoicepb.UpdateInvoiceStatusResponse{
		Message: "Invoice status updated successfully",
	}, nil
}
