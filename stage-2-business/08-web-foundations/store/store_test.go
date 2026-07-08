package store

import "testing"

func TestMemoryStoreListGetCreate(t *testing.T) {
	s := NewSeededMemoryStore()

	initial := s.List()
	if len(initial) != 2 {
		t.Fatalf("seeded store length = %d, want 2", len(initial))
	}

	tests := []struct {
		name string
		id   string
		want string
	}{
		{name: "first seed", id: "1", want: "HTTP handlers are functions"},
		{name: "second seed", id: "2", want: "Middleware wraps handlers"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article, ok := s.Get(tt.id)
			if !ok {
				t.Fatalf("Get(%q) missing article", tt.id)
			}
			if article.Title != tt.want {
				t.Fatalf("title = %q, want %q", article.Title, tt.want)
			}
		})
	}

	created := s.Create("JSON belongs at the boundary", "Handlers decode requests and encode responses.", []string{"json", "rest"})
	if created.ID != "3" {
		t.Fatalf("created ID = %q, want 3", created.ID)
	}
	if got := len(s.List()); got != 3 {
		t.Fatalf("store length after create = %d, want 3", got)
	}
}

func TestMemoryStoreReturnsCopies(t *testing.T) {
	s := NewSeededMemoryStore()
	article, ok := s.Get("1")
	if !ok {
		t.Fatal("seed article missing")
	}
	article.Tags[0] = "mutated"

	again, ok := s.Get("1")
	if !ok {
		t.Fatal("seed article missing on second get")
	}
	if again.Tags[0] == "mutated" {
		t.Fatal("Get returned shared tag slice")
	}
}
