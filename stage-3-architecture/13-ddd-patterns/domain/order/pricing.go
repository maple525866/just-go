package order

import (
	"fmt"

	"just-go/stage-3-architecture/13-ddd-patterns/domain"
)

// DiscountPolicy is a domain strategy used by PricingService.
type DiscountPolicy interface {
	Discount(subtotal Money) (Money, error)
}

type NoDiscount struct{}

func (NoDiscount) Discount(subtotal Money) (Money, error) {
	return NewMoney(0, subtotal.Currency())
}

// PercentageDiscount applies an integer percentage from 0 through 100.
type PercentageDiscount struct{ Percent int }

func (p PercentageDiscount) Discount(subtotal Money) (Money, error) {
	if p.Percent < 0 || p.Percent > 100 {
		return Money{}, fmt.Errorf("%w: discount percent must be between 0 and 100", domain.ErrValidation)
	}
	percent := int64(p.Percent)
	// Split the calculation to avoid overflowing before division. Both terms
	// are bounded by subtotal.Minor() when percent is between 0 and 100.
	discount := subtotal.Minor()/100*percent + subtotal.Minor()%100*percent/100
	return NewMoney(discount, subtotal.Currency())
}

// PricingService performs a domain calculation without persistence or dispatch.
type PricingService struct{}

func (PricingService) Calculate(lines []Line, policy DiscountPolicy) (Money, error) {
	if len(lines) == 0 {
		return Money{}, fmt.Errorf("%w: pricing requires at least one line", domain.ErrValidation)
	}
	if policy == nil {
		return Money{}, fmt.Errorf("%w: discount policy is required", domain.ErrValidation)
	}

	total, err := NewMoney(0, lines[0].unitPrice.Currency())
	if err != nil {
		return Money{}, err
	}
	for _, line := range lines {
		lineTotal, lineErr := line.Total()
		if lineErr != nil {
			return Money{}, lineErr
		}
		total, err = total.Add(lineTotal)
		if err != nil {
			return Money{}, err
		}
	}
	discount, err := policy.Discount(total)
	if err != nil {
		return Money{}, err
	}
	if discount.Currency() != total.Currency() {
		return Money{}, domain.ErrCurrencyMismatch
	}
	if discount.Minor() > total.Minor() {
		return Money{}, fmt.Errorf("%w: discount exceeds subtotal", domain.ErrValidation)
	}
	return NewMoney(total.Minor()-discount.Minor(), total.Currency())
}
