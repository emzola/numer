package models

import (
	pb "github.com/emzola/numer/user-service/proto"
)

// ConvertUserToProto converts a Go model struct to protobuf User message.
func ConvertUserToProto(user *User) *pb.User {
	return &pb.User{
		Id:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
}

// ConvertUserToProto converts a Go model struct to protobuf User message.
func ConvertCustomerToProto(customer *Customer) *pb.Customer {
	return &pb.Customer{
		Id:      customer.ID,
		UserId:  customer.UserID,
		Name:    customer.Name,
		Email:   customer.Email,
		Address: customer.Address,
	}
}
