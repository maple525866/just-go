package auth

import "testing"

func TestPasswordHashAndTokenRoundTrip(t *testing.T) {
	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == "secret" || !CheckPassword(hash, "secret") || CheckPassword(hash, "wrong") {
		t.Fatalf("password hash/check mismatch: %q", hash)
	}

	mgr := NewTokenManager([]byte("test-secret"))
	token, err := mgr.Sign(7, "alice")
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	claims, err := mgr.Verify(token)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if claims.UserID != 7 || claims.Username != "alice" {
		t.Fatalf("claims = %+v", claims)
	}
	if _, err := mgr.Verify(token + "tampered"); err == nil {
		t.Fatalf("tampered token should fail")
	}
}

func TestParseBearerRequiresBearerScheme(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{name: "valid", header: "Bearer abc.def", want: "abc.def"},
		{name: "lowercase scheme", header: "bearer abc.def", want: "abc.def"},
		{name: "extra spaces", header: "  Bearer   abc.def  ", want: "abc.def"},
		{name: "raw token", header: "abc.def", want: ""},
		{name: "empty", header: "", want: ""},
		{name: "empty bearer", header: "Bearer ", want: ""},
		{name: "wrong scheme", header: "Basic abc.def", want: ""},
		{name: "too many fields", header: "Bearer abc def", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseBearer(tt.header); got != tt.want {
				t.Fatalf("ParseBearer(%q) = %q, want %q", tt.header, got, tt.want)
			}
		})
	}
}
