package client

import (
	"context"
	"fmt"

	notificationpb "github.com/emzola/numer/notification-service/proto"
	"google.golang.org/grpc"
)

type NotificationClient struct {
	client notificationpb.NotificationServiceClient
}

func NewNotificationClient(conn *grpc.ClientConn) *NotificationClient {
	return &NotificationClient{client: notificationpb.NewNotificationServiceClient(conn)}
}

func (c *NotificationClient) SendReminder(ctx context.Context, invoiceId int64, email string, numDays int) error {
	req := &notificationpb.SendNotificationRequest{
		Email:   email,
		Message: fmt.Sprintf("Reminder: Invoice %d is due in %d days", invoiceId, numDays),
		Subject: "Invoice Reminder",
	}
	_, err := c.client.SendNotification(ctx, req)
	return err
}
