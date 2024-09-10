package models

import (
	pb "github.com/emzola/numer/activity-service/proto"
)

// ConvertActivityToProto converts a Go model struct to protobuf Activity message.
func ConvertActivityToProto(a *Activity) *pb.Activity {
	return &pb.Activity{
		InvoiceId:   a.InvoiceID,
		UserId:      a.UserID,
		Action:      a.Action,
		Description: a.Description,
		Timestamp:   a.Timestamp.String(),
	}
}
