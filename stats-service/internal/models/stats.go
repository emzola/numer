package models

type Stats struct {
	TotalInvoices        int64
	TotalPaidInvoices    int64
	TotalOverdueInvoices int64
	TotalDraftInvoices   int64
	TotalUnpaidInvoices  int64

	TotalAmountPaid    int64
	TotalAmountOverdue int64
	TotalAmountDraft   int64
	TotalAmountUnpaid  int64
}
