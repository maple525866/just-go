package order

import (
	"errors"
	"math"
	"testing"

	"just-go/stage-3-architecture/13-ddd-patterns/domain"
)

func TestMoneyValidationAndArithmetic(t *testing.T) {
	t.Parallel()

	if _, err := NewMoney(-1, "CNY"); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("NewMoney negative error = %v", err)
	}
	if _, err := NewMoney(10, "GBP"); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("NewMoney currency error = %v", err)
	}

	left := mustMoney(t, 1250, "cny")
	right := mustMoney(t, 250, "CNY")
	sum, err := left.Add(right)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if sum.Minor() != 1500 || sum.Currency() != "CNY" {
		t.Fatalf("Add() = %d %s", sum.Minor(), sum.Currency())
	}
	if left.Minor() != 1250 || right.Minor() != 250 {
		t.Fatal("Add mutated an operand")
	}

	product, err := right.Multiply(3)
	if err != nil || !product.Equal(mustMoney(t, 750, "CNY")) {
		t.Fatalf("Multiply() = %#v, %v", product, err)
	}
	if _, err = left.Add(mustMoney(t, 1, "USD")); !errors.Is(err, domain.ErrCurrencyMismatch) {
		t.Fatalf("Add currency error = %v", err)
	}
}

func TestPercentageDiscountDoesNotOverflow(t *testing.T) {
	t.Parallel()

	subtotal := mustMoney(t, math.MaxInt64, "CNY")
	discount, err := (PercentageDiscount{Percent: 2}).Discount(subtotal)
	if err != nil {
		t.Fatalf("Discount() error = %v", err)
	}
	want := int64(184467440737095516)
	if discount.Minor() != want {
		t.Fatalf("Discount() = %d, want %d", discount.Minor(), want)
	}
}

func TestAddressValidationNormalizationAndEquality(t *testing.T) {
	t.Parallel()

	if _, err := NewAddress("", "road", "city", "100000", "CN"); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("NewAddress error = %v", err)
	}
	one, err := NewAddress(" Alice ", " 1 Go Road ", " Beijing ", " 100000 ", " cn ")
	if err != nil {
		t.Fatalf("NewAddress() error = %v", err)
	}
	two, err := NewAddress("Alice", "1 Go Road", "Beijing", "100000", "CN")
	if err != nil {
		t.Fatalf("NewAddress() error = %v", err)
	}
	if !one.Equal(two) || one.Country() != "CN" || one.Recipient() != "Alice" {
		t.Fatalf("normalized address = %#v", one)
	}
}

func mustMoney(t *testing.T, minor int64, currency string) Money {
	t.Helper()
	money, err := NewMoney(minor, currency)
	if err != nil {
		t.Fatalf("NewMoney() error = %v", err)
	}
	return money
}

func mustAddress(t *testing.T) Address {
	t.Helper()
	address, err := NewAddress("Alice", "1 Go Road", "Beijing", "100000", "CN")
	if err != nil {
		t.Fatalf("NewAddress() error = %v", err)
	}
	return address
}
