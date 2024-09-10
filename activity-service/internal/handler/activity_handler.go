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

func (h *ActivityHandler) LogActivity(ctx context.Context, req *pb.LogActivityRequest) (*pb.LogActivityResponse, error) {
	err := h.service.LogActivity(ctx, req.InvoiceId, req.UserId, req.Action, req.Description)
	if err != nil {
		return nil, err
	}
	return &pb.LogActivityResponse{Status: "Success"}, nil
}

func (h *ActivityHandler) GetRecentActivities(ctx context.Context, req *pb.GetRecentActivitiesRequest) (*pb.GetRecentActivitiesResponse, error) {
	activities, err := h.service.GetRecentActivities(ctx, req.UserId, int(req.Limit))
	if err != nil {
		return nil, err
	}

	protoActivities := make([]*pb.Activity, len(activities))
	for i, activity := range activities {
		protoActivities[i] = models.ConvertActivityToProto(activity)
	}

	return &pb.GetRecentActivitiesResponse{Activities: protoActivities}, nil
}

func (h *ActivityHandler) GetAllActivities(ctx context.Context, req *pb.GetAllActivitiesRequest) (*pb.GetAllActivitiesResponse, error) {
	activities, err := h.service.GetAllActivities(ctx, req.InvoiceId)
	if err != nil {
		return nil, err
	}

	protoActivities := make([]*pb.Activity, len(activities))
	for i, activity := range activities {
		protoActivities[i] = models.ConvertActivityToProto(activity)
	}

	return &pb.GetAllActivitiesResponse{Activities: protoActivities}, nil
}
