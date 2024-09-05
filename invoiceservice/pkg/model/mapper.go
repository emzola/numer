package model

import (
	invoicepb "github.com/emzola/numer/invoiceservice/genproto"
)

// ConvertProtoItemsToInvoiceItems converts a slice of protobuf InvoiceItem messages to the corresponding Go model structs.
func ConvertProtoItemsToInvoiceItems(protoItems []*invoicepb.InvoiceItem) []InvoiceItem {
	items := make([]InvoiceItem, len(protoItems))
	for i, protoItem := range protoItems {
		items[i] = InvoiceItem{
			Description: protoItem.Description,
			Quantity:    protoItem.Quantity,
			UnitPrice:   protoItem.UnitPrice,
		}
	}
	return items
}

// Helper function to convert Go model struct to protobuf Invoice message
func ConvertInvoiceToProto(inv *Invoice) *invoicepb.Invoice {
	protoItems := make([]*invoicepb.InvoiceItem, len(inv.Items))
	for i, item := range inv.Items {
		protoItems[i] = &invoicepb.InvoiceItem{
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		}
	}

	return &invoicepb.Invoice{
		InvoiceId:          inv.InvoiceID,
		UserId:             inv.UserID,
		CustomerId:         inv.CustomerID,
		InvoiceNumber:      inv.InvoiceNumber,
		Status:             invoicepb.InvoiceStatus(inv.Status),
		IssueDate:          inv.IssueDate,
		DueDate:            inv.DueDate,
		BillingCurrency:    inv.BillingCurrency,
		Items:              protoItems,
		DiscountPercentage: inv.DiscountPercentage,
		Subtotal:           inv.Subtotal,
		Discount:           inv.Discount,
		Total:              inv.Total,
		PaymentInformation: &invoicepb.PaymentInformation{
			AccountName:   inv.PaymentInformation.AccountName,
			AccountNumber: inv.PaymentInformation.AccountNumber,
			BankName:      inv.PaymentInformation.BankName,
			RoutingNumber: inv.PaymentInformation.RoutingNumber,
		},
		Note: inv.Note,
	}
}
