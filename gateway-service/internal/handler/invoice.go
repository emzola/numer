package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/emzola/numer/gateway-service/internal/grpcutil"
	invoicepb "github.com/emzola/numer/invoice-service/proto"
	userpb "github.com/emzola/numer/user-service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) CreateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user := h.contextGetUser(r)

	// Decode the JSON body into HTTP Request struct
	var httpReq CreateInvoiceHTTPReq
	err := h.decodeJSON(w, r, &httpReq)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	// Convert time.Time to protobuf Timestamp
	issueDateProto := timestamppb.New(httpReq.IssueDate)
	dueDateProto := timestamppb.New(httpReq.DueDate)

	// Convert the HTTP request into the gRPC CreateInvoiceRequest
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

	// Map Invoice items from HTTP request to gRPC request with []*InvoiceItem
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
	grpcRes, err := client.CreateInvoice(ctx, grpcReq)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Map the gRPC CreateInvoiceResponse back to the HTTP response
	invoiceResp := CreateInvoiceHTTPResp{
		InvoiceID: grpcRes.InvoiceId,
	}

	err = h.encodeJSON(w, http.StatusCreated, envelope{"invoice": invoiceResp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) GetInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID param
	invoiceId, err := h.readIDParam(r)
	if err != nil {
		h.notFoundResponse(w, r)
		return
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
	grpcRes, err := client.GetInvoice(context.Background(), &invoicepb.GetInvoiceRequest{InvoiceId: invoiceId})
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Step 5: Map the gRPC GetInvoiceResponse back to the HTTP response
	invoiceResp := GetInvoiceHTTPResp{
		InvoiceID:          grpcRes.Invoice.Id,
		UserID:             grpcRes.Invoice.UserId,
		CustomerID:         grpcRes.Invoice.CustomerId,
		IssueDate:          grpcRes.Invoice.IssueDate.AsTime(),
		DueDate:            grpcRes.Invoice.DueDate.AsTime(),
		Currency:           grpcRes.Invoice.Currency,
		Items:              convertInvoiceItems(grpcRes.Invoice.Items),
		DiscountPercentage: grpcRes.Invoice.DiscountPercentage,
		AccountName:        grpcRes.Invoice.AccountName,
		AccountNumber:      grpcRes.Invoice.AccountNumber,
		BankName:           grpcRes.Invoice.BankName,
		RoutingNumber:      grpcRes.Invoice.RoutingNumber,
		Note:               grpcRes.Invoice.Note,
	}

	err = h.encodeJSON(w, http.StatusOK, envelope{"invoice": invoiceResp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) UpdateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	// Extract invoice ID param
	invoiceId, err := h.readIDParam(r)
	if err != nil {
		h.notFoundResponse(w, r)
		return
	}

	// Decode the JSON body into the HTTP request struct
	var httpReq UpdateInvoiceHTTPReq
	err = h.decodeJSON(w, r, &httpReq)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	// Convert time.Time to protobuf Timestamp
	issueDateProto := timestamppb.New(httpReq.IssueDate)
	dueDateProto := timestamppb.New(httpReq.DueDate)

	// Convert the HTTP request into the gRPC UpdateInvoiceRequest
	grpcReq := &invoicepb.UpdateInvoiceRequest{
		InvoiceId:          invoiceId,
		Status:             httpReq.Status,
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

	// Map Invoice items from HTTP request to gRPC request with []*InvoiceItem
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
	grpcRes, err := client.UpdateInvoice(ctx, grpcReq)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Map the gRPC UpdateInvoiceResponse back to the HTTP response
	updateInvResp := UpdateInvoiceHTTPResp{
		InvoiceID: grpcRes.InvoiceId,
		Message:   grpcRes.Message,
	}

	err = h.encodeJSON(w, http.StatusOK, envelope{"invoice": updateInvResp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) GetInvoicesHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user := h.contextGetUser(r)

	// Read url query params
	qs := r.URL.Query()
	pageSize := h.ReadInt(qs, "page_size", 10)
	pageToken := h.ReadString(qs, "page_token", "")

	// Convert the HTTP request into the gRPC ListInvoicesRequest
	grpcReq := &invoicepb.ListInvoicesRequest{
		UserId:    user.Id,
		PageSize:  int32(pageSize),
		PageToken: pageToken,
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
	grpcRes, err := client.ListInvoices(context.Background(), grpcReq)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Map the gRPC ListInvoicesResponse back to the HTTP response
	invoiceResp := GetInvoicesHTTPResp{
		Invoices:      convertInvoices(grpcRes.Invoices),
		NextPageToken: grpcRes.NextPageToken,
	}

	err = h.encodeJSON(w, http.StatusOK, envelope{"invoices": invoiceResp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) ScheduleInvoiceReminderHandler(w http.ResponseWriter, r *http.Request) {
	// Extract invoice ID param
	invoiceId, err := h.readIDParam(r)
	if err != nil {
		h.notFoundResponse(w, r)
		return
	}

	// Decode the JSON body into the HTTP request struct
	var httpReq ScheduleInvoiceReminderHTTPReq
	err = h.decodeJSON(w, r, &httpReq)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Create gRPC connection to user service
	userConn, err := grpcutil.ServiceConnection(ctx, "user-service", h.registry)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	defer userConn.Close()

	userClient := userpb.NewUserServiceClient(userConn)

	// Create gRPC connection to invoice service
	invConn, err := grpcutil.ServiceConnection(ctx, "invoice-service", h.registry)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	defer invConn.Close()

	invClient := invoicepb.NewInvoiceServiceClient(invConn)

	// Fetch customer email associated with invoice
	invoiceResp, err := invClient.GetInvoice(ctx, &invoicepb.GetInvoiceRequest{InvoiceId: invoiceId})
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	customerResp, err := userClient.GetCustomer(ctx, &userpb.GetCustomerRequest{CustomerId: invoiceResp.Invoice.CustomerId})
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Prepare the gRPC SendInvoiceRequest
	grpcReq := &invoicepb.ScheduleInvoiceReminderRequest{
		InvoiceId:     invoiceResp.Invoice.Id,
		CustomerEmail: customerResp.Customer.Email,
		ReminderType:  httpReq.ReminderType,
	}

	// Call the ScheduleInvoiceReminder gRPC method
	grpcRes, err := invClient.ScheduleInvoiceReminder(ctx, grpcReq)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Map the gRPC ScheduleInvoiceReminderResponse to the HTTP response
	resp := ScheduleInvoiceReminderHTTPResp{
		Status: grpcRes.Status,
	}

	err = h.encodeJSON(w, http.StatusOK, envelope{"message": resp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) SendInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	// Extract invoice ID param
	invoiceId, err := h.readIDParam(r)
	if err != nil {
		h.notFoundResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Create gRPC connection to user service
	userConn, err := grpcutil.ServiceConnection(ctx, "user-service", h.registry)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	defer userConn.Close()

	userClient := userpb.NewUserServiceClient(userConn)

	// Create gRPC connection to invoice service
	invConn, err := grpcutil.ServiceConnection(ctx, "invoice-service", h.registry)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	defer invConn.Close()

	invClient := invoicepb.NewInvoiceServiceClient(invConn)

	// Fetch customer email associated with invoice
	invoiceResp, err := invClient.GetInvoice(ctx, &invoicepb.GetInvoiceRequest{InvoiceId: invoiceId})
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	customerResp, err := userClient.GetCustomer(ctx, &userpb.GetCustomerRequest{CustomerId: invoiceResp.Invoice.CustomerId})
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Prepare the gRPC SendInvoiceRequest
	grpcReq := &invoicepb.SendInvoiceRequest{
		InvoiceId:     invoiceResp.Invoice.Id,
		CustomerEmail: customerResp.Customer.Email,
	}

	// Call the SendInvoice gRPC method
	grpcRes, err := invClient.SendInvoice(ctx, grpcReq)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Map the gRPC SendInvoiceResponse to the HTTP response
	resp := SendInvoiceHTTPResp{
		Status: grpcRes.Status,
	}

	err = h.encodeJSON(w, http.StatusOK, envelope{"message": resp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// Struct to capture the HTTP request JSON data
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

// Struct to capture the HTTP response
type CreateInvoiceHTTPResp struct {
	InvoiceID int64 `json:"invoice_id"`
}

// Struct to capture the HTTP response
type GetInvoiceHTTPResp struct {
	InvoiceID          int64         `json:"invoice_id"`
	UserID             int64         `json:"user_id"`
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

// Struct to capture the HTTP request JSON data
type UpdateInvoiceHTTPReq struct {
	Status             string        `json:"status"`
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

// Struct to capture the HTTP response
type UpdateInvoiceHTTPResp struct {
	InvoiceID int64  `json:"invoice_id"`
	Message   string `json:"message"`
}

// Struct to capture the HTTP response
type GetInvoicesHTTPResp struct {
	Invoices      []InvoiceHTTP `json:"invoices"`
	NextPageToken string        `json:"next_page_token"`
}

// Struct to represent an Invoice in the HTTP response
type InvoiceHTTP struct {
	InvoiceID          int64         `json:"invoice_id"`
	UserID             int64         `json:"user_id"`
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

// Struct to capture the HTTP response
type SendInvoiceHTTPResp struct {
	Status string `json:"status"`
}

type ScheduleInvoiceReminderHTTPReq struct {
	ReminderType int32
}

// Struct to capture the HTTP response
type ScheduleInvoiceReminderHTTPResp struct {
	Status string `json:"status"`
}
