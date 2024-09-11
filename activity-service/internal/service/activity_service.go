package service

import (
	"context"
	"time"

	"github.com/emzola/numer/activity-service/internal/models"
)

type activityRepository interface {
	LogActivity(ctx context.Context, activity *models.Activity)
	GetUserActivities(ctx context.Context, userID int64, limit int) ([]*models.Activity, error)
	GetInvoiceActivities(ctx context.Context, invoiceID int64) ([]*models.Activity, error)
}

type ActivityService struct {
	repo activityRepository
}

func NewActivityService(repo activityRepository) *ActivityService {
	return &ActivityService{repo: repo}
}

func (s *ActivityService) LogActivity(ctx context.Context, invoiceID, userID int64, action, description string) {
	activity := &models.Activity{
		InvoiceID:   invoiceID,
		UserID:      userID,
		Action:      action,
		Description: description,
		Timestamp:   time.Now(),
	}
	s.repo.LogActivity(ctx, activity)
}

func (s *ActivityService) GetUserActivities(ctx context.Context, userID int64, limit int) ([]*models.Activity, error) {
	return s.repo.GetUserActivities(ctx, userID, limit)
}

func (s *ActivityService) GetInvoiceActivities(ctx context.Context, invoiceID int64) ([]*models.Activity, error) {
	return s.repo.GetInvoiceActivities(ctx, invoiceID)
}
