package main

import (
	"log"
	"os"

	"google.golang.org/grpc"

	blogv1 "just-go/stage-3-architecture/capstone-3-blog-ms/api/blog/v1"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/comments"
	appruntime "just-go/stage-3-architecture/capstone-3-blog-ms/internal/runtime"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/tracekit"
)

func main() {
	address := env("COMMENT_ADDR", ":9003")
	log.Printf("comment-svc listening on %s", address)
	log.Fatal(appruntime.ServeGRPC(address, func(server *grpc.Server) { blogv1.RegisterCommentServiceServer(server, comments.NewService()) }, grpc.UnaryInterceptor(tracekit.UnaryServerInterceptor("comment-svc", nil))))
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
