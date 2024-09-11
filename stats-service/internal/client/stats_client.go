package client

import (
	"context"

	invoicepb "github.com/emzola/numer/invoice-service/proto"
	"google.golang.org/grpc"
)

type StatsClient struct {
	invoiceClient invoicepb.InvoiceServiceClient
}

func NewStatsClient(cc *grpc.ClientConn) *StatsClient {
	return &StatsClient{invoiceClient: invoicepb.NewInvoiceServiceClient(cc)}
}

func (r *StatsClient) GetAllInvoices(ctx context.Context) ([]*invoicepb.Invoice, error) {
	req := &invoicepb.ListInvoicesRequest{UserId: 1}
	resp, err := r.invoiceClient.ListInvoices(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Invoices, nil
}
