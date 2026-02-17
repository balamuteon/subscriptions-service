package domain

import "errors"

var ErrNotImplemented = errors.New("not implemented")

var (
	ErrSubscriptionNotFound  = errors.New("subscription not found")
	ErrMissingRequiredFields = errors.New("missing required fields")
	ErrInvalidID             = errors.New("invalid id")
	ErrInvalidServiceName    = errors.New("invalid service name")
	ErrInvalidUserID         = errors.New("invalid user id")
	ErrInvalidPrice          = errors.New("invalid price")
	ErrInvalidStartDate      = errors.New("invalid start date")
	ErrInvalidEndDate        = errors.New("invalid end date")
	ErrInvalidFromDate       = errors.New("invalid from date")
	ErrInvalidToDate         = errors.New("invalid to date")
	ErrInvalidPeriod         = errors.New("invalid period")
)

type ValidationError struct {
	Err error
}

func (v *ValidationError) Error() string {
	if v == nil || v.Err == nil {
		return "validation error"
	}
	return v.Err.Error()
}

func (v *ValidationError) Unwrap() error {
	return v.Err
}
