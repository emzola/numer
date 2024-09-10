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
	userClient         *client.UserClient
	pb.UnimplementedRemindersServiceServer
}

func NewScheduler(invoiceClient *client.InvoiceClient, notificationClient *client.NotificationClient, userClient *client.UserClient) *Scheduler {
	return &Scheduler{
		invoiceClient:      invoiceClient,
		notificationClient: notificationClient,
		userClient:         userClient,
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

func (s *Scheduler) sendReminder(invoice *invoicepb.Invoice, numDays int) {
	// Get customer email associated with invoice
	customerEmail, err := s.userClient.GetCustomerEmail(context.TODO(), invoice)
	if err != nil {
		log.Println("Error getting customer: ", err)
	}

	// Call the notification service here to send an email to customer
	log.Printf("Sending reminder for Invoice ID: %d, Due in %d days", invoice.Id, numDays)
	s.notificationClient.SendReminder(context.TODO(), invoice.Id, customerEmail, numDays)
}
