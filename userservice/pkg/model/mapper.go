package model

import (
	userpb "github.com/emzola/numer/userservice/genproto"
)

// ConvertUserToProto converts Go model struct to protobuf User message
func ConvertUserToProto(u *User) *userpb.User {
	return &userpb.User{
		Id:     u.ID,
		Name:   u.Name,
		Email:  u.Email,
		Role:   u.Role,
		Active: u.Active,
	}
}
