package capstone3_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	blogv1 "just-go/stage-3-architecture/capstone-3-blog-ms/api/blog/v1"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/comments"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/gateway"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/posts"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/resilience"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/tracekit"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/users"
)

func TestBlogMicroservicesEndToEnd(t *testing.T) {
	exporter := &tracekit.MemoryExporter{}
	userConn := startService(t, "user-svc", exporter, func(server *grpc.Server) {
		blogv1.RegisterUserServiceServer(server, users.NewService([]byte("test-secret")))
	})
	postConn := startService(t, "post-svc", exporter, func(server *grpc.Server) {
		blogv1.RegisterPostServiceServer(server, posts.NewService())
	})
	commentConn := startService(t, "comment-svc", exporter, func(server *grpc.Server) {
		blogv1.RegisterCommentServiceServer(server, comments.NewService())
	})
	handler, err := gateway.New(gateway.Options{
		Users: blogv1.NewUserServiceClient(userConn), Posts: blogv1.NewPostServiceClient(postConn), Comments: blogv1.NewCommentServiceClient(commentConn),
		Exporter: exporter, Limiter: resilience.NewLimiter(20, 20), CommentTimeout: time.Second, RetryAttempts: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()

	client := server.Client()
	var auth blogv1.AuthResponse
	traceparent := "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01"
	doJSON(t, client, http.MethodPost, server.URL+"/api/register", traceparent, "", `{"username":"alice","password":"secret1"}`, http.StatusCreated, &auth)
	doJSON(t, client, http.MethodPost, server.URL+"/api/login", traceparent, "", `{"username":"alice","password":"secret1"}`, http.StatusOK, &auth)
	var post blogv1.Post
	doJSON(t, client, http.MethodPost, server.URL+"/api/posts", traceparent, auth.Token, `{"title":"Go 微服务","body":"边界与治理","tags":["go","ddd"]}`, http.StatusCreated, &post)
	var root blogv1.Comment
	doJSON(t, client, http.MethodPost, server.URL+"/api/posts/1/comments", traceparent, auth.Token, `{"body":"first"}`, http.StatusCreated, &root)
	doJSON(t, client, http.MethodPost, server.URL+"/api/posts/1/comments", traceparent, auth.Token, `{"parent_id":1,"body":"reply"}`, http.StatusCreated, &blogv1.Comment{})
	var detail struct {
		Post             *blogv1.Post      `json:"post"`
		Author           *blogv1.User      `json:"author"`
		Comments         []*blogv1.Comment `json:"comments"`
		CommentsDegraded bool              `json:"comments_degraded"`
	}
	doJSON(t, client, http.MethodGet, server.URL+"/api/posts/1", traceparent, "", "", http.StatusOK, &detail)
	if detail.Post == nil || detail.Author == nil || detail.Post.Title != "Go 微服务" || detail.Author.Username != "alice" || detail.CommentsDegraded || len(detail.Comments) != 1 || len(detail.Comments[0].Replies) != 1 {
		t.Fatalf("unexpected detail: post=%v author=%v comments=%d degraded=%v", detail.Post, detail.Author, len(detail.Comments), detail.CommentsDegraded)
	}
	spans := exporter.Spans()
	services := map[string]bool{}
	for _, span := range spans {
		if span.TraceID == "0123456789abcdef0123456789abcdef" {
			services[span.Service] = true
		}
	}
	for _, service := range []string{"gateway", "user-svc", "post-svc", "comment-svc"} {
		if !services[service] {
			t.Fatalf("trace missing %s: %#v", service, spans)
		}
	}
}

func startService(t *testing.T, name string, exporter tracekit.Exporter, register func(*grpc.Server)) *grpc.ClientConn {
	t.Helper()
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer(grpc.UnaryInterceptor(tracekit.UnaryServerInterceptor(name, exporter)))
	register(server)
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)
	connection, err := grpc.NewClient("passthrough:///bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(tracekit.UnaryClientInterceptor()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = connection.Close() })
	return connection
}

func doJSON(t *testing.T, client *http.Client, method, url, traceparent, token, body string, wantStatus int, target any) {
	t.Helper()
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("traceparent", traceparent)
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var buffer bytes.Buffer
	_, _ = buffer.ReadFrom(response.Body)
	if response.StatusCode != wantStatus {
		t.Fatalf("%s %s status=%d body=%s", method, url, response.StatusCode, buffer.String())
	}
	if target != nil && buffer.Len() > 0 {
		if err := json.Unmarshal(buffer.Bytes(), target); err != nil {
			t.Fatalf("decode %s: %v", buffer.String(), err)
		}
	}
}
