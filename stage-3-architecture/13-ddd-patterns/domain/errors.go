// Package domain contains primitives shared by the Chapter 13 domain model.
package domain

import "errors"

var (
	ErrValidation       = errors.New("domain validation failed")
	ErrCurrencyMismatch = errors.New("currency mismatch")
	ErrInvalidState     = errors.New("invalid aggregate state")
	ErrDuplicateLine    = errors.New("duplicate order line")
	ErrLineNotFound     = errors.New("order line not found")
)
