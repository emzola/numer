package service

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrInvalidRequest = errors.New("the request is invalid")
)
