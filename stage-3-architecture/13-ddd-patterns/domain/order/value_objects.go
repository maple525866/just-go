// Package order models the ordering aggregate and its domain behavior.
package order

import (
	"fmt"
	"math"
	"strings"

	"just-go/stage-3-architecture/13-ddd-patterns/domain"
)

var supportedCurrencies = map[string]struct{}{
	"CNY": {},
	"EUR": {},
	"USD": {},
}

// Money is an immutable monetary value stored in minor units.
type Money struct {
	minor    int64
	currency string
}

func NewMoney(minor int64, currency string) (Money, error) {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if minor < 0 {
		return Money{}, fmt.Errorf("%w: money amount must not be negative", domain.ErrValidation)
	}
	if _, ok := supportedCurrencies[currency]; !ok {
		return Money{}, fmt.Errorf("%w: unsupported currency %q", domain.ErrValidation, currency)
	}
	return Money{minor: minor, currency: currency}, nil
}

func (m Money) Minor() int64     { return m.minor }
func (m Money) Currency() string { return m.currency }

func (m Money) Equal(other Money) bool {
	return m == other
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, domain.ErrCurrencyMismatch
	}
	if other.minor > math.MaxInt64-m.minor {
		return Money{}, fmt.Errorf("%w: money addition overflows", domain.ErrValidation)
	}
	return Money{minor: m.minor + other.minor, currency: m.currency}, nil
}

func (m Money) Multiply(quantity int) (Money, error) {
	if quantity <= 0 {
		return Money{}, fmt.Errorf("%w: quantity must be positive", domain.ErrValidation)
	}
	if m.minor > math.MaxInt64/int64(quantity) {
		return Money{}, fmt.Errorf("%w: money multiplication overflows", domain.ErrValidation)
	}
	return Money{minor: m.minor * int64(quantity), currency: m.currency}, nil
}

// Address is an immutable shipping address.
type Address struct {
	recipient string
	street    string
	city      string
	postal    string
	country   string
}

func NewAddress(recipient, street, city, postal, country string) (Address, error) {
	values := []string{recipient, street, city, postal, country}
	for i := range values {
		values[i] = strings.TrimSpace(values[i])
		if values[i] == "" {
			return Address{}, fmt.Errorf("%w: address fields are required", domain.ErrValidation)
		}
	}
	return Address{
		recipient: values[0],
		street:    values[1],
		city:      values[2],
		postal:    values[3],
		country:   strings.ToUpper(values[4]),
	}, nil
}

func (a Address) Recipient() string { return a.recipient }
func (a Address) Street() string    { return a.street }
func (a Address) City() string      { return a.city }
func (a Address) Postal() string    { return a.postal }
func (a Address) Country() string   { return a.country }

func (a Address) Equal(other Address) bool { return a == other }
