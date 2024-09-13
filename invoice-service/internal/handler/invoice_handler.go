package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/emzola/numer/invoice-service/internal/models"
	"github.com/emzola/numer/invoice-service/internal/service"
	"github.com/emzola/numer/invoice-service/internal/service/rabbitmq"
	pb "github.com/emzola/numer/invoice-service/proto"
	notificationpb "github.com/emzola/numer/notification-service/proto"
	reminderpb "github.com/emzola/numer/reminder-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	maxRetries        = 5
	initialBackoffMs  = 100
	backoffMultiplier = 2
)

type InvoiceHandler struct {
	service            *service.InvoiceService
	publisher          *rabbitmq.Publisher
	reminderClient     reminderpb.ReminderServiceClient
	notificationClient notificationpb.NotificationServiceClient
	pb.UnimplementedInvoiceServiceServer
}

func NewInvoiceHandler(service *service.InvoiceService, publisher *rabbitmq.Publisher, reminderClient reminderpb.ReminderServiceClient, notificationClient notificationpb.NotificationServiceClient) *InvoiceHandler {
	return &InvoiceHandler{
		service:            service,
		publisher:          publisher,
		reminderClient:     reminderClient,
		notificationClient: notificationClient,
	}
}

func (h *InvoiceHandler) CreateInvoice(ctx context.Context, req *pb.CreateInvoiceRequest) (*pb.CreateInvoiceResponse, error) {
	invoice := &models.Invoice{
		UserID:             req.UserId,
		CustomerID:         req.CustomerId,
		IssueDate:          req.IssueDate.AsTime(),
		DueDate:            req.DueDate.AsTime(),
		Currency:           req.Currency,
		DiscountPercentage: req.DiscountPercentage,
		AccountName:        req.AccountName,
		AccountNumber:      req.AccountNumber,
		BankName:           req.BankName,
		RoutingNumber:      req.RoutingNumber,
		Note:               req.Note,
	}

	// Add invoice items
	for _, item := range req.Items {
		invoice.Items = append(invoice.Items, &models.InvoiceItem{
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		})
	}

	invoice, err := h.service.CreateInvoice(ctx, invoice)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Publish activity to rabbitMQ
	activity := map[string]interface{}{
		"invoice_id":  invoice.ID,
		"user_id":     invoice.UserID,
		"action":      "Invoice creation",
		"description": fmt.Sprintf("Created invoice %s", invoice.InvoiceNumber),
	}

	h.publisher.Publish(activity)

	return &pb.CreateInvoiceResponse{InvoiceId: invoice.ID}, nil
}

func (h *InvoiceHandler) GetInvoice(ctx context.Context, req *pb.GetInvoiceRequest) (*pb.GetInvoiceResponse, error) {
	invoice, err := h.service.GetInvoice(ctx, req.InvoiceId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "invoice not found")
	}

	return &pb.GetInvoiceResponse{Invoice: models.ConvertInvoiceToProto(invoice)}, nil
}

func (h *InvoiceHandler) UpdateInvoice(ctx context.Context, req *pb.UpdateInvoiceRequest) (*pb.UpdateInvoiceResponse, error) {
	invoice, err := h.service.GetInvoice(ctx, req.InvoiceId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "invoice not found")
	}

	invoice.Status = req.Status
	invoice.IssueDate = req.IssueDate.AsTime()
	invoice.DueDate = req.DueDate.AsTime()
	invoice.Currency = req.Currency
	invoice.DiscountPercentage = req.DiscountPercentage
	invoice.AccountName = req.AccountName
	invoice.AccountName = req.AccountNumber
	invoice.BankName = req.BankName
	invoice.RoutingNumber = req.RoutingNumber
	invoice.Note = req.Note

	// Update invoice items
	invoice.Items = []*models.InvoiceItem{}
	for _, itemReq := range req.Items {
		invoice.Items = append(invoice.Items, &models.InvoiceItem{
			Description: itemReq.Description,
			Quantity:    itemReq.Quantity,
			UnitPrice:   itemReq.UnitPrice,
		})
	}

	// Call service to update invoice
	err = h.service.UpdateInvoice(ctx, invoice)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateInvoiceResponse{
		InvoiceId: invoice.ID,
		Message:   "invoice successfully updated",
	}, nil
}

func (s *InvoiceHandler) ListInvoices(ctx context.Context, req *pb.ListInvoicesRequest) (*pb.ListInvoicesResponse, error) {
	invoices, nextPageToken, err := s.service.ListInvoicesByUserID(ctx, req.UserId, int(req.PageSize), req.PageToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	protoInvoices := make([]*pb.Invoice, len(invoices))
	for i, invoice := range invoices {
		protoInvoices[i] = models.ConvertInvoiceToProto(invoice)
	}

	return &pb.ListInvoicesResponse{
		Invoices:      protoInvoices,
		NextPageToken: nextPageToken,
	}, nil
}

func (h *InvoiceHandler) ScheduleInvoiceReminder(ctx context.Context, req *pb.ScheduleInvoiceReminderRequest) (*pb.ScheduleInvoiceReminderResponse, error) {
	invoice, err := h.service.GetInvoice(ctx, req.InvoiceId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "invoice not found")
	}

	// Calculate the reminder time based on the reminder_type (14, 7, 3, or 1 day before due date)
	var reminderTime time.Time
	switch req.ReminderType {
	case 14:
		reminderTime = invoice.DueDate.AddDate(0, 0, -14)
	case 7:
		reminderTime = invoice.DueDate.AddDate(0, 0, -7)
	case 3:
		reminderTime = invoice.DueDate.AddDate(0, 0, -3)
	case 1:
		reminderTime = invoice.DueDate.AddDate(0, 0, -1)
	default:
		return nil, fmt.Errorf("invalid reminder type")
	}

	// Prepare the email message
	message := fmt.Sprintf("Reminder: Your invoice is due on %s. Please ensure payment of $%d is made.",
		invoice.DueDate.Format("2006-01-02"), invoice.Total)

	// Send the reminder request to the ReminderService
	_, err = h.reminderClient.ScheduleReminder(ctx, &reminderpb.ScheduleReminderRequest{
		InvoiceId:     req.InvoiceId,
		CustomerEmail: req.CustomerEmail,
		ReminderTime:  reminderTime.Format(time.RFC3339),
		Message:       message,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to schedule reminder: %v", err)
	}

	return &pb.ScheduleInvoiceReminderResponse{
		Status: "Reminder scheduled successfully",
	}, nil
}

func (h *InvoiceHandler) SendInvoiceEmail(ctx context.Context, req *pb.SendInvoiceRequest) (*pb.SendInvoiceResponse, error) {
	invoice, err := h.service.GetInvoice(ctx, req.InvoiceId)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %v", err)
	}

	// Prepare the email message body
	message := fmt.Sprintf("Dear %d, \n\nPlease find your invoice for $%d due on %s. \n\n%s",
		invoice.ID, invoice.Total, invoice.DueDate, invoice.Note)

	// Retry sending email
	err = h.retrySendEmail(ctx, "test@example.com", "Your Invoice", message)
	if err != nil {
		return nil, err
	}

	// Publish activity to rabbitMQ
	activity := map[string]interface{}{
		"invoice_id":  invoice.ID,
		"user_id":     invoice.UserID,
		"action":      "Invoice sent",
		"description": fmt.Sprintf("Sent invoice %s to user %d", invoice.InvoiceNumber, invoice.UserID),
	}

	h.publisher.Publish(activity)

	return &pb.SendInvoiceResponse{
		Status: "email sent successfully",
	}, nil
}

func (h *InvoiceHandler) retrySendEmail(ctx context.Context, email, subject, message string) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		_, err = h.notificationClient.SendNotification(ctx, &notificationpb.SendNotificationRequest{
			Email:   email,
			Subject: subject,
			Message: message,
		})
		if err == nil {
			return nil
		}

		// Log the retry attempt
		fmt.Printf("SendInvoiceEmail attempt %d failed: %v\n", i+1, err)

		// Exponential backoff
		backoffDuration := time.Duration(initialBackoffMs*(1<<i)) * time.Millisecond
		time.Sleep(backoffDuration)
	}
	return fmt.Errorf("failed to send email after %d attempts: %w", maxRetries, err)
}
