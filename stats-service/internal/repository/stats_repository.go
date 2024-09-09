package repository

import (
	"context"

	invoicepb "github.com/emzola/numer/invoice-service/proto"
	"google.golang.org/grpc"
)

type StatsRepository struct {
	invoiceClient invoicepb.InvoiceServiceClient
}

func NewStatsRepository(cc *grpc.ClientConn) *StatsRepository {
	return &StatsRepository{
		invoiceClient: invoicepb.NewInvoiceServiceClient(cc),
	}
}

func (r *StatsRepository) GetAllInvoices(ctx context.Context) ([]*invoicepb.Invoice, error) {
	req := &invoicepb.ListInvoicesRequest{}
	resp, err := r.invoiceClient.ListInvoices(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Invoices, nil
}
