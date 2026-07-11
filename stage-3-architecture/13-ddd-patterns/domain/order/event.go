package order

import "time"

const OrderConfirmedName = "order.confirmed"

type ConfirmedLine struct {
	productID string
	quantity  int
}

func (l ConfirmedLine) ProductID() string { return l.productID }
func (l ConfirmedLine) Quantity() int     { return l.quantity }

// OrderConfirmed is an immutable fact produced by the Order aggregate.
type OrderConfirmed struct {
	orderID    string
	lines      []ConfirmedLine
	occurredAt time.Time
}

func newOrderConfirmed(aggregate *Order, occurredAt time.Time) OrderConfirmed {
	lines := make([]ConfirmedLine, len(aggregate.lines))
	for i, line := range aggregate.lines {
		lines[i] = ConfirmedLine{productID: line.productID, quantity: line.quantity}
	}
	return OrderConfirmed{orderID: aggregate.id, lines: lines, occurredAt: occurredAt}
}

func (e OrderConfirmed) Name() string           { return OrderConfirmedName }
func (e OrderConfirmed) OrderID() string        { return e.orderID }
func (e OrderConfirmed) OccurredAt() time.Time  { return e.occurredAt }
func (e OrderConfirmed) Lines() []ConfirmedLine { return append([]ConfirmedLine(nil), e.lines...) }
