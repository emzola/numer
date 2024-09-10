package models

import "time"

type Activity struct {
	InvoiceID   string
	UserID      string
	Action      string
	Description string
	Timestamp   time.Time
}
