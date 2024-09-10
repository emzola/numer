package handler

import (
	"context"

	"github.com/emzola/numer/activity-service/internal/models"
	"github.com/emzola/numer/activity-service/internal/service"
	pb "github.com/emzola/numer/activity-service/proto"
)

type ActivityHandler struct {
	service *service.ActivityService
	pb.UnimplementedActivityServiceServer
}

func NewActivityHandler(service *service.ActivityService) *ActivityHandler {
	return &ActivityHandler{service: service}
}

func (h *ActivityHandler) GetUserActivities(ctx context.Context, req *pb.GetUserActivitiesRequest) (*pb.GetUserActivitiesResponse, error) {
	activities, err := h.service.GetUserActivities(ctx, req.UserId, int(req.Limit))
	if err != nil {
		return nil, err
	}

	protoActivities := make([]*pb.Activity, len(activities))
	for i, activity := range activities {
		protoActivities[i] = models.ConvertActivityToProto(activity)
	}

	return &pb.GetUserActivitiesResponse{Activities: protoActivities}, nil
}

func (h *ActivityHandler) GetInvoiceActivities(ctx context.Context, req *pb.GetInvoiceActivitiesRequest) (*pb.GetInvoiceActivitiesResponse, error) {
	activities, err := h.service.GetInvoiceActivities(ctx, req.InvoiceId)
	if err != nil {
		return nil, err
	}

	protoActivities := make([]*pb.Activity, len(activities))
	for i, activity := range activities {
		protoActivities[i] = models.ConvertActivityToProto(activity)
	}

	return &pb.GetInvoiceActivitiesResponse{Activities: protoActivities}, nil
}
