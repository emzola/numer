package handler

import "github.com/emzola/numer/gateway-service/pkg/discovery/consul"

type Handler struct {
	registry *consul.Registry
}

func NewGatewayHandler(registry *consul.Registry) *Handler {
	return &Handler{registry: registry}
}
