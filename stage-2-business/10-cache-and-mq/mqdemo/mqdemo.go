package mqdemo

import (
	"fmt"
	"sync"
	"time"
)

// Message is the stable queue payload.
type Message struct {
	ID    string
	Topic string
	Body  string
}

// Delivery is a leased message that must be acked.
type Delivery struct {
	ID         string
	DeliveryID string
	Topic      string
	Body       string
}

type leasedMessage struct {
	message    Message
	deliveryID string
	deadline   time.Time
}

// Broker is an in-memory teaching broker with at-least-once semantics.
type Broker struct {
	mu                sync.Mutex
	now               time.Time
	manual            bool
	visibilityTimeout time.Duration
	nextMessageID     int
	nextDeliveryID    int
	ready             []Message
	inflight          map[string]leasedMessage
}

// NewBroker creates a broker with a controllable clock.
func NewBroker(now time.Time, visibilityTimeout time.Duration) *Broker {
	manual := time.Since(now) > time.Second || time.Until(now) > time.Second
	return &Broker{now: now, manual: manual, visibilityTimeout: visibilityTimeout, inflight: map[string]leasedMessage{}}
}

// Advance moves the broker clock forward.
func (b *Broker) Advance(d time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.manual = true
	b.now = b.now.Add(d)
}

func (b *Broker) currentTime() time.Time {
	if b.manual {
		return b.now
	}
	return time.Now()
}

// Publish enqueues one message.
func (b *Broker) Publish(topic, body string) Message {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.nextMessageID++
	message := Message{ID: fmt.Sprintf("msg-%d", b.nextMessageID), Topic: topic, Body: body}
	b.ready = append(b.ready, message)
	return message
}

// Fetch leases the next available message.
func (b *Broker) Fetch() (Delivery, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.ready) == 0 {
		return Delivery{}, false
	}
	message := b.ready[0]
	b.ready = b.ready[1:]
	b.nextDeliveryID++
	deliveryID := fmt.Sprintf("delivery-%d", b.nextDeliveryID)
	b.inflight[deliveryID] = leasedMessage{message: message, deliveryID: deliveryID, deadline: b.currentTime().Add(b.visibilityTimeout)}
	return Delivery{ID: message.ID, DeliveryID: deliveryID, Topic: message.Topic, Body: message.Body}, true
}

// Ack removes an in-flight message.
func (b *Broker) Ack(deliveryID string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.inflight[deliveryID]; !ok {
		return false
	}
	delete(b.inflight, deliveryID)
	return true
}

// RequeueExpired makes timed-out in-flight messages ready again.
func (b *Broker) RequeueExpired() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for deliveryID, leased := range b.inflight {
		if !b.currentTime().Before(leased.deadline) {
			delete(b.inflight, deliveryID)
			b.ready = append(b.ready, leased.message)
		}
	}
}
