package clock

import (
	"testing"
	"time"
)

func TestFormatSummary(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		d    time.Duration
		want Summary
	}{
		{
			name: "formats UTC and duration",
			time: time.Date(2026, 6, 27, 10, 30, 0, 0, time.FixedZone("CST", 8*3600)),
			d:    90 * time.Minute,
			want: Summary{Formatted: "2026-06-27T02:30:00Z", Minutes: 90},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatSummary(tt.time, tt.d); got != tt.want {
				t.Fatalf("FormatSummary() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestCountTicks(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{name: "finite ticks", n: 2, want: 2},
		{name: "no ticks", n: 0, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountTicks(tt.n, time.Millisecond); got != tt.want {
				t.Fatalf("CountTicks() = %d, want %d", got, tt.want)
			}
		})
	}
}
