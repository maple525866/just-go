package gateway

import "testing"

func TestAuthorizedBearerHeader(t *testing.T) {
	tests := []struct {
		name   string
		header string
		token  string
		want   bool
	}{
		{name: "valid", header: "Bearer teaching-token", token: "teaching-token", want: true},
		{name: "missing", token: "teaching-token"},
		{name: "wrong scheme", header: "Basic teaching-token", token: "teaching-token"},
		{name: "wrong token", header: "Bearer wrong", token: "teaching-token"},
		{name: "extra whitespace", header: "Bearer  teaching-token", token: "teaching-token"},
		{name: "blank expected token", header: "Bearer ", token: " "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := authorizedBearer(tt.header, tt.token); got != tt.want {
				t.Fatalf("authorizedBearer(%q, %q) = %v, want %v", tt.header, tt.token, got, tt.want)
			}
		})
	}
}
