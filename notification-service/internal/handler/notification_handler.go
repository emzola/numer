package handler

import (
	"context"
	"log"

	"github.com/emzola/numer/notification-service/internal/email"
	pb "github.com/emzola/numer/notification-service/proto"
)

type NotificationHandler struct {
	emailSender *email.EmailSender
	pb.UnimplementedNotificationServiceServer
}

func NewNotificationHandler(emailSender *email.EmailSender) *NotificationHandler {
	return &NotificationHandler{emailSender: emailSender}
}

// SendNotification implements the gRPC method for sending notifications
func (s *NotificationHandler) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	log.Printf("Received request to send notification to: %s", req.Email)

	// Send the email
	err := s.emailSender.SendEmail(req.Email, req.Subject, req.Message)
	if err != nil {
		return nil, err
	}

	// Return success response
	return &pb.SendNotificationResponse{
		Status: "Success",
	}, nil
}
