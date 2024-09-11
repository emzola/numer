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

func (h *StatsHandler) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	stats, err := h.service.GetStats(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.GetStatsResponse{
		TotalInvoices:        stats.TotalInvoices,
		TotalPaidInvoices:    stats.TotalPaidInvoices,
		TotalOverdueInvoices: stats.TotalOverdueInvoices,
		TotalDraftInvoices:   stats.TotalDraftInvoices,
		TotalUnpaidInvoices:  stats.TotalUnpaidInvoices,
		TotalAmountPaid:      stats.TotalAmountPaid,
		TotalAmountOverdue:   stats.TotalAmountOverdue,
		TotalAmountDraft:     stats.TotalAmountDraft,
		TotalAmountUnpaid:    stats.TotalAmountUnpaid,
	}, nil
}
