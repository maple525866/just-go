package channel

import (
	"reflect"
	"testing"
	"time"
)

func TestPingPong(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "message crosses unbuffered channel", in: "ping", want: "ping"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PingPong(tt.in); got != tt.want {
				t.Fatalf("PingPong() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBufferedQueue(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{name: "keeps order", in: []int{1, 2, 3}, want: []int{1, 2, 3}},
		{name: "empty", in: []int{}, want: []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BufferedQueue(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("BufferedQueue() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestCloseAndRange(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{name: "reads until closed", in: []string{"a", "b"}, want: []string{"a", "b"}},
		{name: "empty", in: []string{}, want: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CloseAndRange(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("CloseAndRange() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestReceiveWithTimeout(t *testing.T) {
	tests := []struct {
		name    string
		delay   time.Duration
		timeout time.Duration
		want    string
	}{
		{name: "receives before timeout", delay: 0, timeout: 50 * time.Millisecond, want: "ready"},
		{name: "times out", delay: 50 * time.Millisecond, timeout: time.Millisecond, want: "timeout"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReceiveWithTimeout(tt.delay, tt.timeout); got != tt.want {
				t.Fatalf("ReceiveWithTimeout() = %q, want %q", got, tt.want)
			}
		})
	}
}
