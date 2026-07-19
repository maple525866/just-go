package users

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"

	blogv1 "just-go/stage-3-architecture/capstone-3-blog-ms/api/blog/v1"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/rpcerr"
)

type record struct {
	id       int64
	username string
	hash     string
}

type Service struct {
	blogv1.UnimplementedUserServiceServer
	mu       sync.RWMutex
	nextID   int64
	byID     map[int64]record
	byName   map[string]int64
	tokenKey []byte
}

func NewService(tokenKey []byte) *Service {
	return &Service{nextID: 1, byID: map[int64]record{}, byName: map[string]int64{}, tokenKey: append([]byte(nil), tokenKey...)}
}

func (s *Service) Register(_ context.Context, req *blogv1.RegisterRequest) (*blogv1.AuthResponse, error) {
	username := strings.TrimSpace(req.GetUsername())
	if username == "" || len(req.GetPassword()) < 6 {
		return nil, rpcerr.ToStatus(rpcerr.ErrInvalid)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.byName[username]; exists {
		return nil, rpcerr.ToStatus(rpcerr.ErrAlreadyExists)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, rpcerr.ToStatus(err)
	}
	item := record{id: s.nextID, username: username, hash: string(hash)}
	s.nextID++
	s.byID[item.id], s.byName[item.username] = item, item.id
	return &blogv1.AuthResponse{User: toProto(item), Token: s.sign(item)}, nil
}

func (s *Service) Login(_ context.Context, req *blogv1.LoginRequest) (*blogv1.AuthResponse, error) {
	s.mu.RLock()
	id, ok := s.byName[strings.TrimSpace(req.GetUsername())]
	item := s.byID[id]
	s.mu.RUnlock()
	if !ok || bcrypt.CompareHashAndPassword([]byte(item.hash), []byte(req.GetPassword())) != nil {
		return nil, rpcerr.ToStatus(rpcerr.ErrUnauthenticated)
	}
	return &blogv1.AuthResponse{User: toProto(item), Token: s.sign(item)}, nil
}

func (s *Service) ValidateToken(_ context.Context, req *blogv1.ValidateTokenRequest) (*blogv1.User, error) {
	parts := strings.Split(req.GetToken(), ".")
	if len(parts) != 2 {
		return nil, rpcerr.ToStatus(rpcerr.ErrUnauthenticated)
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, rpcerr.ToStatus(rpcerr.ErrUnauthenticated)
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, rpcerr.ToStatus(rpcerr.ErrUnauthenticated)
	}
	mac := hmac.New(sha256.New, s.tokenKey)
	_, _ = mac.Write(payload)
	if !hmac.Equal(signature, mac.Sum(nil)) {
		return nil, rpcerr.ToStatus(rpcerr.ErrUnauthenticated)
	}
	id, err := strconv.ParseInt(string(payload), 10, 64)
	if err != nil {
		return nil, rpcerr.ToStatus(rpcerr.ErrUnauthenticated)
	}
	return s.GetUser(context.Background(), &blogv1.GetUserRequest{Id: id})
}

func (s *Service) GetUser(_ context.Context, req *blogv1.GetUserRequest) (*blogv1.User, error) {
	s.mu.RLock()
	item, ok := s.byID[req.GetId()]
	s.mu.RUnlock()
	if !ok {
		return nil, rpcerr.ToStatus(rpcerr.ErrNotFound)
	}
	return toProto(item), nil
}

func (s *Service) sign(item record) string {
	payload := []byte(strconv.FormatInt(item.id, 10))
	mac := hmac.New(sha256.New, s.tokenKey)
	_, _ = mac.Write(payload)
	return base64.RawURLEncoding.EncodeToString(payload) + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func toProto(item record) *blogv1.User {
	return &blogv1.User{Id: item.id, Username: item.username}
}
