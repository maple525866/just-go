package breaker

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitOpensAfterFailureThreshold(t *testing.T) {
	errDownstream := errors.New("downstream")
	circuit := New[string](Config{Name: "test", FailureThreshold: 2, Timeout: 50 * time.Millisecond, MaxRequests: 1})

	for i := 0; i < 2; i++ {
		_, err := circuit.Execute(func() (string, error) { return "", errDownstream })
		if !errors.Is(err, errDownstream) {
			t.Fatalf("failure %d err = %v", i, err)
		}
	}

	_, err := circuit.Execute(func() (string, error) { return "should-not-run", nil })
	if !errors.Is(err, ErrOpen) {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
	if got := circuit.State(); got != "open" {
		t.Fatalf("state = %s", got)
	}
}

func TestCircuitEventuallyAllowsProbeAfterTimeout(t *testing.T) {
	circuit := New[string](Config{Name: "probe", FailureThreshold: 1, Timeout: 10 * time.Millisecond, MaxRequests: 1})
	_, _ = circuit.Execute(func() (string, error) { return "", errors.New("boom") })
	if _, err := circuit.Execute(func() (string, error) { return "", nil }); !errors.Is(err, ErrOpen) {
		t.Fatalf("expected open rejection, got %v", err)
	}

	deadline := time.Now().Add(300 * time.Millisecond)
	for time.Now().Before(deadline) {
		got, err := circuit.Execute(func() (string, error) { return "recovered", nil })
		if err == nil && got == "recovered" {
			if state := circuit.State(); state != "closed" {
				t.Fatalf("state after successful probe = %s", state)
			}
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("breaker did not allow a successful recovery probe before deadline")
}
