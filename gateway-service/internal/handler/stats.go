package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/emzola/numer/gateway-service/internal/grpcutil"
	statspb "github.com/emzola/numer/stats-service/proto"
)

func (h *Handler) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user := h.contextGetUser(r)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Create gRPC connection to stats service
	conn, err := grpcutil.ServiceConnection(ctx, "stats-service", h.registry)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	defer conn.Close()

	client := statspb.NewStatsServiceClient(conn)
	grpcRes, err := client.GetStats(context.Background(), &statspb.GetStatsRequest{UserId: user.Id})
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Map the gRPC GetInvoiceResponse back to the HTTP response
	statsResp := GetStatsHTTPResp{
		TotalInvoices:        grpcRes.TotalInvoices,
		TotalPaidInvoices:    grpcRes.TotalPaidInvoices,
		TotalOverdueInvoices: grpcRes.TotalOverdueInvoices,
		TotalDraftInvoices:   grpcRes.TotalDraftInvoices,
		TotalUnpaidInvoices:  grpcRes.TotalUnpaidInvoices,
		TotalAmountPaid:      grpcRes.TotalAmountPaid,
		TotalAmountOverdue:   grpcRes.TotalAmountOverdue,
		TotalAmountDraft:     grpcRes.TotalAmountDraft,
		TotalAmountUnpaid:    grpcRes.TotalAmountUnpaid,
	}

	err = h.encodeJSON(w, http.StatusOK, envelope{"stats": statsResp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// Struct to capture the HTTP response
type GetStatsHTTPResp struct {
	TotalInvoices        int64 `json:"total_invoices"`
	TotalPaidInvoices    int64 `json:"total_paid_invoices"`
	TotalOverdueInvoices int64 `json:"total_overdue_invoices"`
	TotalDraftInvoices   int64 `json:"total_draft_invoices"`
	TotalUnpaidInvoices  int64 `json:"total_unpaid_invoices"`

	TotalAmountPaid    int64 `json:"total_amount_paid"`
	TotalAmountOverdue int64 `json:"total_amount_overdue"`
	TotalAmountDraft   int64 `json:"total_amount_draft"`
	TotalAmountUnpaid  int64 `json:"total_amount_unpaid"`
}
