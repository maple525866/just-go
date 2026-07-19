package main

import (
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	blogv1 "just-go/stage-3-architecture/capstone-3-blog-ms/api/blog/v1"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/gateway"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/resilience"
	appruntime "just-go/stage-3-architecture/capstone-3-blog-ms/internal/runtime"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/tracekit"
)

func main() {
	usersConn := dial(env("USER_TARGET", "localhost:9001"))
	defer func() { _ = usersConn.Close() }()
	postsConn := dial(env("POST_TARGET", "localhost:9002"))
	defer func() { _ = postsConn.Close() }()
	commentsConn := dial(env("COMMENT_TARGET", "localhost:9003"))
	defer func() { _ = commentsConn.Close() }()
	handler, err := gateway.New(gateway.Options{
		Users: blogv1.NewUserServiceClient(usersConn), Posts: blogv1.NewPostServiceClient(postsConn), Comments: blogv1.NewCommentServiceClient(commentsConn),
		Limiter: resilience.NewLimiter(100, 50), CommentTimeout: 300 * time.Millisecond, RetryAttempts: 2,
	})
	if err != nil {
		log.Fatal(err)
	}
	address := env("GATEWAY_ADDR", ":8080")
	log.Printf("gateway listening on %s", address)
	log.Fatal(appruntime.ServeHTTP(address, handler))
}

func dial(target string) *grpc.ClientConn {
	connection, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(tracekit.UnaryClientInterceptor()))
	if err != nil {
		log.Fatal(err)
	}
	return connection
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
