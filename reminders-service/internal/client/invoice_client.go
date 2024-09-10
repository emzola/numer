package client

import (
	"context"

	invoicepb "github.com/emzola/numer/invoice-service/proto"
	"google.golang.org/grpc"
)

type InvoiceClient struct {
	client invoicepb.InvoiceServiceClient
}

func NewInvoiceClient(conn *grpc.ClientConn) *InvoiceClient {
	return &InvoiceClient{client: invoicepb.NewInvoiceServiceClient(conn)}
}

func (c *InvoiceClient) GetInvoicesDueSoon(ctx context.Context) ([]*invoicepb.Invoice, error) {
	req := &invoicepb.GetDueInvoicesRequest{
		DaysBeforeDue: 14, // Fetch invoices due in 14 days or less
	}
	res, err := c.client.GetDueInvoices(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.Invoices, nil
}
