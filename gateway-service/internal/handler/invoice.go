package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/emzola/numer/gateway-service/internal/grpcutil"
	invoicepb "github.com/emzola/numer/invoice-service/proto"
)

func (h *Handler) GetInvoices(w http.ResponseWriter, r *http.Request) {
	userID, err := h.readIDParam(r)
	if err != nil {
		http.Error(w, "id not found", http.StatusNotFound)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Create gRPC connection to invoice service
	conn, err := grpcutil.ServiceConnection(ctx, "invoice-service", h.registry)
	if err != nil {
		http.Error(w, "could not connect to Invoice Service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := invoicepb.NewInvoiceServiceClient(conn)
	res, err := client.ListInvoices(context.Background(), &invoicepb.ListInvoicesRequest{UserId: userID})
	if err != nil {
		http.Error(w, "error fetching invoices", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}
