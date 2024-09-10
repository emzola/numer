package service

import (
	"context"

	invoicepb "github.com/emzola/numer/invoice-service/proto"
	"github.com/emzola/numer/stats-service/internal/models"
)

type statsClient interface {
	GetAllInvoices(ctx context.Context) ([]*invoicepb.Invoice, error)
}

type StatsService struct {
	invoiceStatsClient statsClient
}

func NewStatsService(invoiceStatsClient statsClient) *StatsService {
	return &StatsService{invoiceStatsClient: invoiceStatsClient}
}

func (s *StatsService) GetStats(ctx context.Context) (*models.Stats, error) {
	invoices, err := s.invoiceStatsClient.GetAllInvoices(ctx)
	if err != nil {
		return nil, err
	}

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
