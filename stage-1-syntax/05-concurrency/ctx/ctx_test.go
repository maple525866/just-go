package ctx

import (
	"context"
	"testing"
	"time"
)

func TestCancellationDemo(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "worker exits on cancel", want: "cancelled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CancellationDemo(); got != tt.want {
				t.Fatalf("CancellationDemo() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTimeoutDemo(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		want    bool
	}{
		{name: "deadline exceeded", timeout: time.Millisecond, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeoutDemo(tt.timeout); got != tt.want {
				t.Fatalf("TimeoutDemo() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestWorkerWithContext(t *testing.T) {
	tests := []struct {
		name   string
		cancel bool
		job    string
		want   string
	}{
		{name: "processes job", job: "build", want: "processed:build"},
		{name: "returns cancellation", cancel: true, want: "context canceled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent, cancel := context.WithCancel(context.Background())
			defer cancel()
			jobs := make(chan string, 1)
			if tt.cancel {
				cancel()
			} else {
				jobs <- tt.job
			}
			if got := WorkerWithContext(parent, jobs); got != tt.want {
				t.Fatalf("WorkerWithContext() = %q, want %q", got, tt.want)
			}
		})
	}
}
