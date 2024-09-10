package client

import (
	"context"

	invoicepb "github.com/emzola/numer/invoice-service/proto"
	userpb "github.com/emzola/numer/user-service/proto"
	"google.golang.org/grpc"
)

type UserClient struct {
	client userpb.UserServiceClient
}

func NewUserClient(conn *grpc.ClientConn) *UserClient {
	return &UserClient{client: userpb.NewUserServiceClient(conn)}
}

func (c *UserClient) GetCustomerEmail(ctx context.Context, invoice *invoicepb.Invoice) (string, error) {
	req := &userpb.GetCustomerRequest{
		CustomerId: invoice.Id,
	}
	res, err := c.client.GetCustomer(ctx, req)
	if err != nil {
		return "", err
	}
	return res.Email, nil
}
