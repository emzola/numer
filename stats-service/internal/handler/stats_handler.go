package handler

import (
	"context"

	"github.com/emzola/numer/stats-service/internal/service"
	pb "github.com/emzola/numer/stats-service/proto"
)

type StatsHandler struct {
	service *service.StatsService
	pb.UnimplementedStatsServiceServer
}

func NewStatsHandler(service *service.StatsService) *StatsHandler {
	return &StatsHandler{service: service}
}

func (h *StatsHandler) GetMetrics(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	metrics, err := h.service.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetStatsResponse{
		TotalInvoices:        metrics.TotalInvoices,
		TotalPaidInvoices:    metrics.TotalPaidInvoices,
		TotalOverdueInvoices: metrics.TotalOverdueInvoices,
		TotalDraftInvoices:   metrics.TotalDraftInvoices,
		TotalUnpaidInvoices:  metrics.TotalUnpaidInvoices,
		TotalAmountPaid:      metrics.TotalAmountPaid,
		TotalAmountOverdue:   metrics.TotalAmountOverdue,
		TotalAmountDraft:     metrics.TotalAmountDraft,
		TotalAmountUnpaid:    metrics.TotalAmountUnpaid,
	}, nil
}
