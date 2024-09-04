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
