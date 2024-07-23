package domain

import "errors"

var (
	ErrInvalidDate = errors.New("invalid date")
	ErrInternal    = errors.New("internal error")
)
