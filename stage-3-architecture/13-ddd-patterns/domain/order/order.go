package order

import (
	"fmt"
	"strings"
	"time"

	"just-go/stage-3-architecture/13-ddd-patterns/domain"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusConfirmed Status = "confirmed"
)

// Line is an entity whose identity is stable within an Order aggregate.
type Line struct {
	id        string
	productID string
	name      string
	unitPrice Money
	quantity  int
}

func NewLine(id, productID, name string, unitPrice Money, quantity int) (Line, error) {
	id = strings.TrimSpace(id)
	productID = strings.TrimSpace(productID)
	name = strings.TrimSpace(name)
	line := Line{id: id, productID: productID, name: name, unitPrice: unitPrice, quantity: quantity}
	if err := line.validate(); err != nil {
		return Line{}, err
	}
	return line, nil
}

func (l Line) ID() string        { return l.id }
func (l Line) ProductID() string { return l.productID }
func (l Line) Name() string      { return l.name }
func (l Line) UnitPrice() Money  { return l.unitPrice }
func (l Line) Quantity() int     { return l.quantity }

func (l Line) Total() (Money, error) { return l.unitPrice.Multiply(l.quantity) }

func (l Line) validate() error {
	if l.id == "" || l.productID == "" || l.name == "" || l.quantity <= 0 {
		return fmt.Errorf("%w: line id, product, name and positive quantity are required", domain.ErrValidation)
	}
	if _, ok := supportedCurrencies[l.unitPrice.currency]; !ok || l.unitPrice.minor < 0 {
		return fmt.Errorf("%w: line price is invalid", domain.ErrValidation)
	}
	return nil
}

// Order is the aggregate root. All mutation of lines goes through its methods.
type Order struct {
	id          string
	customerID  string
	address     Address
	lines       []Line
	status      Status
	version     uint64
	confirmedAt time.Time
	total       Money
	events      []domain.Event
}

func New(id, customerID string, address Address) (*Order, error) {
	id = strings.TrimSpace(id)
	customerID = strings.TrimSpace(customerID)
	if id == "" || customerID == "" {
		return nil, fmt.Errorf("%w: order and customer ids are required", domain.ErrValidation)
	}
	if address.recipient == "" || address.street == "" || address.city == "" || address.postal == "" || address.country == "" {
		return nil, fmt.Errorf("%w: shipping address is required", domain.ErrValidation)
	}
	return &Order{id: id, customerID: customerID, address: address, status: StatusDraft, version: 1}, nil
}

func (o *Order) ID() string             { return o.id }
func (o *Order) CustomerID() string     { return o.customerID }
func (o *Order) Address() Address       { return o.address }
func (o *Order) Status() Status         { return o.status }
func (o *Order) Version() uint64        { return o.version }
func (o *Order) ConfirmedAt() time.Time { return o.confirmedAt }
func (o *Order) Total() Money           { return o.total }

func (o *Order) Lines() []Line {
	return append([]Line(nil), o.lines...)
}

func (o *Order) AddLine(line Line) error {
	if o.status != StatusDraft {
		return fmt.Errorf("%w: only draft orders can change", domain.ErrInvalidState)
	}
	if err := line.validate(); err != nil {
		return err
	}
	for _, existing := range o.lines {
		if existing.id == line.id || existing.productID == line.productID {
			return domain.ErrDuplicateLine
		}
	}
	lineTotal, err := line.Total()
	if err != nil {
		return err
	}
	nextTotal := lineTotal
	if len(o.lines) > 0 {
		nextTotal, err = o.total.Add(lineTotal)
		if err != nil {
			return err
		}
	}
	o.lines = append(o.lines, line)
	o.total = nextTotal
	o.version++
	return nil
}

func (o *Order) RemoveLine(lineID string) error {
	if o.status != StatusDraft {
		return fmt.Errorf("%w: only draft orders can change", domain.ErrInvalidState)
	}
	for i, line := range o.lines {
		if line.id == lineID {
			remaining := append([]Line(nil), o.lines[:i]...)
			remaining = append(remaining, o.lines[i+1:]...)
			nextTotal, err := subtotal(remaining, line.unitPrice.currency)
			if err != nil {
				return err
			}
			o.lines = remaining
			o.total = nextTotal
			o.version++
			return nil
		}
	}
	return domain.ErrLineNotFound
}

func subtotal(lines []Line, emptyCurrency string) (Money, error) {
	total, err := NewMoney(0, emptyCurrency)
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
	return total, nil
}

func (o *Order) Confirm(total Money, at time.Time) error {
	if o.status != StatusDraft {
		return fmt.Errorf("%w: order is already confirmed", domain.ErrInvalidState)
	}
	if len(o.lines) == 0 {
		return fmt.Errorf("%w: an order needs at least one line", domain.ErrInvalidState)
	}
	if at.IsZero() {
		return fmt.Errorf("%w: confirmation time is required", domain.ErrValidation)
	}
	if total.Currency() != o.lines[0].unitPrice.Currency() {
		return domain.ErrCurrencyMismatch
	}

	o.status = StatusConfirmed
	o.confirmedAt = at
	o.total = total
	o.version++
	o.events = append(o.events, newOrderConfirmed(o, at))
	return nil
}

// PullEvents returns pending events and clears them from this aggregate instance.
func (o *Order) PullEvents() []domain.Event {
	events := append([]domain.Event(nil), o.events...)
	o.events = nil
	return events
}

// Clone returns an independent aggregate snapshot without transient events.
func (o *Order) Clone() *Order {
	if o == nil {
		return nil
	}
	clone := *o
	clone.lines = append([]Line(nil), o.lines...)
	clone.events = nil
	return &clone
}
