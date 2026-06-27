package syncx

import "testing"

func TestCountWithMutex(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{name: "ten goroutines", n: 10, want: 10},
		{name: "none", n: 0, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountWithMutex(tt.n); got != tt.want {
				t.Fatalf("CountWithMutex() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestScoreBoard(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		score     int
		wantScore int
		wantFound bool
	}{
		{name: "found", key: "Ada", score: 95, wantScore: 95, wantFound: true},
		{name: "zero score still found", key: "Bob", score: 0, wantScore: 0, wantFound: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := NewScoreBoard()
			board.Set(tt.key, tt.score)
			got, found := board.Get(tt.key)
			if found != tt.wantFound || got != tt.wantScore {
				t.Fatalf("Get() = (%d, %t), want (%d, %t)", got, found, tt.wantScore, tt.wantFound)
			}
		})
	}
}

func TestInitOnceConcurrently(t *testing.T) {
	tests := []struct {
		name  string
		times int
		want  int
	}{
		{name: "many callers one init", times: 20, want: 1},
		{name: "no callers", times: 0, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitOnceConcurrently(tt.times); got != tt.want {
				t.Fatalf("InitOnceConcurrently() = %d, want %d", got, tt.want)
			}
		})
	}
}
