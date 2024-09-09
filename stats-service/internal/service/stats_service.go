package service

import (
	"context"

	invoicepb "github.com/emzola/numer/invoice-service/proto"
)

type statsRepository interface {
	GetAllInvoices(ctx context.Context) ([]*invoicepb.Invoice, error)
}

type Stats struct {
	TotalInvoices        int64
	TotalPaidInvoices    int64
	TotalOverdueInvoices int64
	TotalDraftInvoices   int64
	TotalUnpaidInvoices  int64

	TotalAmountPaid    int64
	TotalAmountOverdue int64
	TotalAmountDraft   int64
	TotalAmountUnpaid  int64
}

type StatsService struct {
	invoiceRepo statsRepository
}

func NewStatsService(repo statsRepository) *StatsService {
	return &StatsService{invoiceRepo: repo}
}

func (s *StatsService) GetStats(ctx context.Context) (*Stats, error) {
	invoices, err := s.invoiceRepo.GetAllInvoices(ctx)
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

	stats := &Stats{
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
