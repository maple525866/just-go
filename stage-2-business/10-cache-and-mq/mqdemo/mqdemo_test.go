package mqdemo

import (
	"testing"
	"time"
)

func TestPublishAndFetch(t *testing.T) {
	broker := NewBroker(time.Unix(100, 0), time.Minute)
	broker.Publish("article.created", "article-1")

	delivery, ok := broker.Fetch()
	if !ok {
		t.Fatal("Fetch returned no message")
	}
	if delivery.Topic != "article.created" || delivery.Body != "article-1" || delivery.DeliveryID == "" {
		t.Fatalf("delivery = %+v", delivery)
	}
}

func TestAckRemovesMessage(t *testing.T) {
	broker := NewBroker(time.Unix(100, 0), time.Minute)
	broker.Publish("article.created", "article-1")
	delivery, ok := broker.Fetch()
	if !ok {
		t.Fatal("Fetch returned no message")
	}

	if ok := broker.Ack(delivery.DeliveryID); !ok {
		t.Fatal("Ack returned false")
	}
	if next, ok := broker.Fetch(); ok {
		t.Fatalf("Fetch after Ack = %+v, want none", next)
	}
}

func TestUnackedMessageIsRedelivered(t *testing.T) {
	broker := NewBroker(time.Unix(100, 0), time.Second)
	broker.Publish("article.created", "article-1")
	first, ok := broker.Fetch()
	if !ok {
		t.Fatal("Fetch returned no message")
	}

	broker.Advance(2 * time.Second)
	broker.RequeueExpired()
	second, ok := broker.Fetch()
	if !ok {
		t.Fatal("Fetch after RequeueExpired returned no message")
	}
	if second.ID != first.ID || second.DeliveryID == first.DeliveryID {
		t.Fatalf("redelivery = %+v, first = %+v", second, first)
	}
}
