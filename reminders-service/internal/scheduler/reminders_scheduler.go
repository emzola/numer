package scheduler

import (
	"context"
	"log"
	"time"

	invoicepb "github.com/emzola/numer/invoice-service/proto"
	"github.com/emzola/numer/reminders-service/internal/client"
	pb "github.com/emzola/numer/reminders-service/proto"
	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	invoiceClient      *client.InvoiceClient
	notificationClient *client.NotificationClient
	pb.UnimplementedRemindersServiceServer
}

func NewScheduler(invoiceClient *client.InvoiceClient, notificationClient *client.NotificationClient) *Scheduler {
	return &Scheduler{
		invoiceClient:      invoiceClient,
		notificationClient: notificationClient,
	}
}

func (s *Scheduler) Start() {
	cron := gocron.NewScheduler(time.UTC)

	// Run every 24 hours to check for upcoming invoice due dates
	cron.Every(1).Day().Do(s.CheckDueDates)

	cron.StartAsync()
}

func (s *Scheduler) CheckDueDates() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	invoices, err := s.invoiceClient.GetInvoicesDueSoon(ctx)
	if err != nil {
		log.Println("Error fetching invoices: ", err)
		return
	}

	for _, invoice := range invoices {
		// Schedule reminders for 14, 7, 3 days and 24 hours before due date
		s.scheduleReminder(invoice)
	}
}

func (s *Scheduler) scheduleReminder(invoice *invoicepb.Invoice) {
	// Calculate the time difference
	dueDate, _ := time.Parse(time.RFC3339, invoice.DueDate.AsTime().String())

	// 14 days before
	if time.Until(dueDate).Hours() > 14*24 {
		go s.sendReminder(invoice, 14)
	}
	// 7 days before
	if time.Until(dueDate).Hours() > 7*24 {
		go s.sendReminder(invoice, 7)
	}
	// 3 days before
	if time.Until(dueDate).Hours() > 3*24 {
		go s.sendReminder(invoice, 3)
	}
	// 24 hours before
	if time.Until(dueDate).Hours() > 24 {
		go s.sendReminder(invoice, 1)
	}
}

func (s *Scheduler) sendReminder(invoice *invoicepb.Invoice, daysBefore int) {
	// Call the Notification Microservice here to send an email or other notification
	log.Printf("Sending reminder for Invoice ID: %s, Due in %d days", invoice.Id, daysBefore)
	s.notificationClient.SendReminder(context.TODO(), invoice.Id, "TODO", daysBefore)
}
