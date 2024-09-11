package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/emzola/numer/gateway-service/internal/grpcutil"
	invoicepb "github.com/emzola/numer/invoice-service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user := h.contextGetUser(r)

	// Step 1: Decode the JSON body into HTTP Request struct
	var httpReq CreateInvoiceHTTPReq
	err := h.decodeJSON(w, r, &httpReq)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	// Step 2: Convert time.Time to protobuf Timestamp
	issueDateProto := timestamppb.New(httpReq.IssueDate)
	dueDateProto := timestamppb.New(httpReq.DueDate)

	// Step 3: Convert the HTTP request into the gRPC CreateInvoiceRequest
	grpcReq := &invoicepb.CreateInvoiceRequest{
		UserId:             user.Id,
		CustomerId:         httpReq.CustomerID,
		IssueDate:          issueDateProto,
		DueDate:            dueDateProto,
		Currency:           httpReq.Currency,
		DiscountPercentage: httpReq.DiscountPercentage,
		AccountName:        httpReq.AccountName,
		AccountNumber:      httpReq.AccountNumber,
		BankName:           httpReq.BankName,
		RoutingNumber:      httpReq.RoutingNumber,
		Note:               httpReq.Note,
	}

	// Step 4: Map Invoice items from HTTP request to gRPC request with []*InvoiceItem
	for _, item := range httpReq.Items {
		grpcReq.Items = append(grpcReq.Items, &invoicepb.InvoiceItem{
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		})
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Create gRPC connection to invoice service
	conn, err := grpcutil.ServiceConnection(ctx, "invoice-service", h.registry)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	defer conn.Close()

	client := invoicepb.NewInvoiceServiceClient(conn)
	invoice, err := client.CreateInvoice(ctx, grpcReq)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	err = h.encodeJSON(w, http.StatusCreated, envelope{"invoice": invoice}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	// Extract ID param
	invoiceId, err := h.readIDParam(r)
	if err != nil {
		h.notFoundResponse(w, r)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Create gRPC connection to invoice service
	conn, err := grpcutil.ServiceConnection(ctx, "invoice-service", h.registry)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	defer conn.Close()

	client := invoicepb.NewInvoiceServiceClient(conn)
	invoice, err := client.GetInvoice(context.Background(), &invoicepb.GetInvoiceRequest{InvoiceId: invoiceId})
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	err = h.encodeJSON(w, http.StatusOK, envelope{"invoice": invoice}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) GetInvoices(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user := h.contextGetUser(r)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Create gRPC connection to invoice service
	conn, err := grpcutil.ServiceConnection(ctx, "invoice-service", h.registry)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	defer conn.Close()

	client := invoicepb.NewInvoiceServiceClient(conn)
	invoices, err := client.ListInvoices(context.Background(), &invoicepb.ListInvoicesRequest{UserId: user.Id})
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	err = h.encodeJSON(w, http.StatusOK, envelope{"invoices": invoices}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// Struct to capture the HTTP request JSON data (matching the gRPC CreateInvoiceRequest)
type CreateInvoiceHTTPReq struct {
	CustomerID         int64         `json:"customer_id"`
	IssueDate          time.Time     `json:"issue_date"`
	DueDate            time.Time     `json:"due_date"`
	Currency           string        `json:"currency"`
	Items              []InvoiceItem `json:"items"`
	DiscountPercentage int64         `json:"discount_percentage"`
	AccountName        string        `json:"account_name"`
	AccountNumber      string        `json:"account_number"`
	BankName           string        `json:"bank_name"`
	RoutingNumber      string        `json:"routing_number"`
	Note               string        `json:"note"`
}

// Struct for invoice items
type InvoiceItem struct {
	Description string `json:"description"`
	Quantity    int32  `json:"quantity"`
	UnitPrice   int64  `json:"price"`
}
