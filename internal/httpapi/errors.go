package httpapi

import "errors"

var (
	ErrInvalidJSON               = errors.New("invalid json")
	ErrStatusInternalServerError = errors.New("internal server error")
)
