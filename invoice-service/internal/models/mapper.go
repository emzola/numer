package models

import (
	pb "github.com/emzola/numer/invoice-service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ConvertProtoItemsToInvoiceItems converts a slice of protobuf InvoiceItem messages to the corresponding Go model structs.
func ConvertProtoItemsToInvoiceItems(protoItems []*pb.InvoiceItem) []InvoiceItem {
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

// ConvertInvoiceToProto converts a Go model struct to protobuf Invoice message.
func ConvertInvoiceToProto(inv *Invoice) *pb.Invoice {
	protoInvoiceItems := make([]*pb.InvoiceItem, len(inv.Items))
	for i, item := range inv.Items {
		protoInvoiceItems[i] = &pb.InvoiceItem{
			Id:          item.ID,
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		}
	}

	return &pb.Invoice{
		Id:                 inv.ID,
		UserId:             inv.UserID,
		CustomerId:         inv.CustomerID,
		InvoiceNumber:      inv.InvoiceNumber,
		Status:             inv.Status,
		IssueDate:          timestamppb.New(inv.IssueDate),
		DueDate:            timestamppb.New(inv.DueDate),
		Currency:           inv.Currency,
		Items:              protoInvoiceItems,
		DiscountPercentage: inv.DiscountPercentage,
		Subtotal:           inv.Subtotal,
		DiscountAmount:     inv.DiscountAmount,
		Total:              inv.Total,
		AccountName:        inv.AccountName,
		AccountNumber:      inv.AccountName,
		BankName:           inv.BankName,
		RoutingNumber:      inv.RoutingNumber,
		Note:               inv.Note,
	}
}
