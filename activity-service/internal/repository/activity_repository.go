package repository

import (
	"context"
	"database/sql"

	"github.com/emzola/numer/activity-service/internal/models"
)

type ActivityRepository struct {
	db *sql.DB
}

func NewActivityRepository(db *sql.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) LogActivity(ctx context.Context, activity *models.Activity) error {
	query := `INSERT INTO activities (invoice_id, user_id, action, description, timestamp) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, activity.InvoiceID, activity.UserID, activity.Action, activity.Description, activity.Timestamp)
	return err
}

func (r *ActivityRepository) GetRecentActivities(ctx context.Context, userID string, limit int) ([]*models.Activity, error) {
	query := `SELECT invoice_id, user_id, action, description, timestamp FROM activities 
              WHERE user_id = $1 ORDER BY timestamp DESC LIMIT $2`
	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []*models.Activity
	for rows.Next() {
		var activity models.Activity
		if err := rows.Scan(&activity.InvoiceID, &activity.UserID, &activity.Action, &activity.Description, &activity.Timestamp); err != nil {
			return nil, err
		}
		activities = append(activities, &activity)
	}
	return activities, nil
}

func (r *ActivityRepository) GetAllActivities(ctx context.Context, invoiceID string) ([]*models.Activity, error) {
	query := `SELECT invoice_id, user_id, action, description, timestamp FROM activities 
              WHERE invoice_id = $1 ORDER BY timestamp`
	rows, err := r.db.QueryContext(ctx, query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []*models.Activity
	for rows.Next() {
		var activity models.Activity
		if err := rows.Scan(&activity.InvoiceID, &activity.UserID, &activity.Action, &activity.Description, &activity.Timestamp); err != nil {
			return nil, err
		}
		activities = append(activities, &activity)
	}
	return activities, nil
}
