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
