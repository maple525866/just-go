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
