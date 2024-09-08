package models

import "time"

type Invoice struct {
	ID                 int64
	UserID             int64
	CustomerID         int64
	InvoiceNumber      string
	Status             string
	IssueDate          time.Time
	DueDate            time.Time
	Currency           string
	Items              []*InvoiceItem
	DiscountPercentage int64 // Represented as hundredths of a percent (e.g., 1000 = 10%)
	Subtotal           int64 // Represented in cents
	DiscountAmount     int64 // Represented in cents
	Total              int64 // Represented in cents
	AccountName        string
	AccountNumber      string
	BankName           string
	RoutingNumber      string
	Note               string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type InvoiceItem struct {
	ID          int64
	Description string
	Quantity    int32
	UnitPrice   int64 // Represented in cents
}
