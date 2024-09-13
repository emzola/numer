package service

import (
	"context"
	"fmt"
	"log"
	"time"

	notificationpb "github.com/emzola/numer/notification-service/proto"
	"github.com/go-co-op/gocron"
)

type ReminderService struct {
	scheduler          *gocron.Scheduler
	notificationClient notificationpb.NotificationServiceClient
}

func NewReminderService(scheduler *gocron.Scheduler, notificationClient notificationpb.NotificationServiceClient) *ReminderService {
	return &ReminderService{
		scheduler:          scheduler,
		notificationClient: notificationClient,
	}
}

func (s *ReminderService) ScheduleReminder(invoiceID int64, customerEmail string, reminderTime time.Time, message string) error {
	s.scheduler.Every(1).Day().At(reminderTime.Format("15:04")).Do(func() {
		err := sendEmail(s.notificationClient, customerEmail, message)
		if err != nil {
			log.Printf("failed to send reminder email: %v", err)
		}
	})
	s.scheduler.StartAsync()

	return nil
}

func sendEmail(client notificationpb.NotificationServiceClient, email string, message string) error {
	req := &notificationpb.SendNotificationRequest{
		Email:   email,
		Subject: "Invoice Payment Reminder",
		Message: message,
	}

	// Call the Notification Microservice's SendNotification method
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.SendNotification(ctx, req)
	if err != nil {
		log.Printf("failed to send email notification: %v", err)
		return err
	}

	if res.Status != "Success" {
		errMsg := fmt.Sprintf("failed to send email: %s", res.Status)
		log.Println(errMsg)
		return fmt.Errorf(errMsg)
	}

	log.Printf("email reminder sent successfully to %s", email)
	return nil
}
