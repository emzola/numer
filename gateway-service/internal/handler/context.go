package handler

import (
	"context"
	"net/http"

	userpb "github.com/emzola/numer/user-service/proto"
)

type contextKey string

const UserContextKey = contextKey("user")

func (h *Handler) contextSetUser(r *http.Request, user *userpb.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func (h *Handler) contextGetUser(r *http.Request) *userpb.User {
	user, ok := r.Context().Value(UserContextKey).(*userpb.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
