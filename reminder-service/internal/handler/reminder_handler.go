package handler

import (
	"context"
	"log"

	"github.com/emzola/numer/reminder-service/internal/service"
	pb "github.com/emzola/numer/reminder-service/proto"
)

type ReminderHandler struct {
	service *service.ReminderService
	pb.UnimplementedReminderServiceServer
}

func NewReminderHandler(service *service.ReminderService) *ReminderHandler {
	return &ReminderHandler{service: service}
}

func (h *ReminderHandler) ScheduleReminder(ctx context.Context, req *pb.ScheduleReminderRequest) (*pb.ScheduleReminderResponse, error) {
	err := h.service.ScheduleReminder(req.InvoiceId, req.CustomerEmail, req.ReminderTime.AsTime(), req.Message)
	if err != nil {
		log.Printf("failed to schedule reminder: %v", err)
		return &pb.ScheduleReminderResponse{Status: "failed"}, err
	}
	return &pb.ScheduleReminderResponse{Status: "success"}, nil
}
