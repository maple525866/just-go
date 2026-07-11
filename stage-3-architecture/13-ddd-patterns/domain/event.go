package domain

import "time"

// Event is a fact that occurred inside the domain.
type Event interface {
	Name() string
	OccurredAt() time.Time
}
