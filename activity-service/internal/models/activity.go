package models

import "time"

type Activity struct {
	InvoiceID   int64
	UserID      int64
	Action      string
	Description string
	Timestamp   time.Time
}
