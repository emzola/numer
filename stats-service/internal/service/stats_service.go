package service

import (
	"context"

	invoicepb "github.com/emzola/numer/invoice-service/proto"
	"github.com/emzola/numer/stats-service/internal/models"
	"google.golang.org/grpc"
)

// type statsClient interface {
// 	GetAllInvoices(ctx context.Context) ([]*invoicepb.Invoice, error)
// }

type StatsService struct {
	conn *grpc.ClientConn
}

func NewStatsService(conn *grpc.ClientConn) *StatsService {
	return &StatsService{conn: conn}
}

func (s *StatsService) GetStats(ctx context.Context, userId int64) (*models.Stats, error) {
	grpcReq := &invoicepb.ListInvoicesRequest{
		UserId: userId,
	}
	client := invoicepb.NewInvoiceServiceClient(s.conn)
	response, err := client.ListInvoices(ctx, grpcReq)
	if err != nil {
		return nil, err
	}

	invoices := response.Invoices

	var totalPaid, totalOverdue, totalDraft, totalUnpaid int64
	var totalAmountPaid, totalAmountOverdue, totalAmountDraft, totalAmountUnpaid int64

	for _, invoice := range invoices {
		switch invoice.Status {
		case "paid":
			totalPaid++
			totalAmountPaid += invoice.Total
		case "overdue":
			totalOverdue++
			totalAmountOverdue += invoice.Total
		case "draft":
			totalDraft++
			totalAmountDraft += invoice.Total
		case "unpaid":
			totalUnpaid++
			totalAmountUnpaid += invoice.Total
		}
	}

	stats := &models.Stats{
		TotalInvoices:        int64(len(invoices)),
		TotalPaidInvoices:    totalPaid,
		TotalOverdueInvoices: totalOverdue,
		TotalDraftInvoices:   totalDraft,
		TotalUnpaidInvoices:  totalUnpaid,

		TotalAmountPaid:    totalAmountPaid,
		TotalAmountOverdue: totalAmountOverdue,
		TotalAmountDraft:   totalAmountDraft,
		TotalAmountUnpaid:  totalAmountUnpaid,
	}

	return stats, nil
}
