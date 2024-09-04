package model

// Invoice represents invoice data.
type Invoice struct {
	InvoiceID          string
	UserID             string
	CustomerID         string
	InvoiceNumber      string
	Status             InvoiceStatus
	IssueDate          string
	DueDate            string
	BillingCurrency    string
	Items              []InvoiceItem
	DiscountPercentage int64 // Represented as hundredths of a percent (e.g., 1000 = 10%)
	Subtotal           int64 // Represented in cents
	Discount           int64 // Represented in cents
	Total              int64 // Represented in cents
	PaymentInformation PaymentInformation
	Note               string
}

type InvoiceItem struct {
	Description string
	Quantity    int32
	UnitPrice   int64 // Represented in cents
}

type PaymentInformation struct {
	AccountName   string
	AccountNumber string
	BankName      string
	RoutingNumber string
}

// InvoiceStatus represents the status of an invoice.
type InvoiceStatus int32

// Enum values for InvoiceStatus.
const (
	PENDING   InvoiceStatus = 0
	PAID      InvoiceStatus = 1
	OVERDUE   InvoiceStatus = 2
	DRAFT     InvoiceStatus = 3
	CANCELLED InvoiceStatus = 4
)

// String provides a string representation of the InvoiceStatus.
func (s InvoiceStatus) String() string {
	switch s {
	case PENDING:
		return "PENDING"
	case PAID:
		return "PAID"
	case OVERDUE:
		return "OVERDUE"
	case DRAFT:
		return "DRAFT"
	case CANCELLED:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
}
