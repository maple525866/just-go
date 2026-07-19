package users

import (
	"context"
	"testing"

	blogv1 "just-go/stage-3-architecture/capstone-3-blog-ms/api/blog/v1"
)

func TestRegisterLoginAndValidate(t *testing.T) {
	service := NewService([]byte("test-secret"))
	registered, err := service.Register(context.Background(), &blogv1.RegisterRequest{Username: "alice", Password: "secret1"})
	if err != nil {
		t.Fatal(err)
	}
	loggedIn, err := service.Login(context.Background(), &blogv1.LoginRequest{Username: "alice", Password: "secret1"})
	if err != nil || loggedIn.GetToken() == "" {
		t.Fatalf("login=%#v err=%v", loggedIn, err)
	}
	user, err := service.ValidateToken(context.Background(), &blogv1.ValidateTokenRequest{Token: registered.GetToken()})
	if err != nil || user.GetUsername() != "alice" {
		t.Fatalf("user=%#v err=%v", user, err)
	}
}
