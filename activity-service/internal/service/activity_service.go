package service

import (
	"context"
	"time"

	"github.com/emzola/numer/activity-service/internal/models"
)

type activityRepository interface {
	LogActivity(ctx context.Context, activity *models.Activity) error
	GetRecentActivities(ctx context.Context, userID string, limit int) ([]*models.Activity, error)
	GetAllActivities(ctx context.Context, invoiceID string) ([]*models.Activity, error)
}

type ActivityService struct {
	repo activityRepository
}

func NewActivityService(repo activityRepository) *ActivityService {
	return &ActivityService{repo: repo}
}

func (s *ActivityService) LogActivity(ctx context.Context, invoiceID, userID, action, description string) error {
	activity := &models.Activity{
		InvoiceID:   invoiceID,
		UserID:      userID,
		Action:      action,
		Description: description,
		Timestamp:   time.Now(),
	}
	return s.repo.LogActivity(ctx, activity)
}

func (s *ActivityService) GetRecentActivities(ctx context.Context, userID string, limit int) ([]*models.Activity, error) {
	return s.repo.GetRecentActivities(ctx, userID, limit)
}

func (s *ActivityService) GetAllActivities(ctx context.Context, invoiceID string) ([]*models.Activity, error) {
	return s.repo.GetAllActivities(ctx, invoiceID)
}
