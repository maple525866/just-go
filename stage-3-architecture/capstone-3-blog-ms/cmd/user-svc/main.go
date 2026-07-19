package main

import (
	"log"
	"os"

	"google.golang.org/grpc"

	blogv1 "just-go/stage-3-architecture/capstone-3-blog-ms/api/blog/v1"
	appruntime "just-go/stage-3-architecture/capstone-3-blog-ms/internal/runtime"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/tracekit"
	"just-go/stage-3-architecture/capstone-3-blog-ms/internal/users"
)

func main() {
	address := env("USER_ADDR", ":9001")
	service := users.NewService([]byte(env("TOKEN_SECRET", "capstone-3-secret")))
	log.Printf("user-svc listening on %s", address)
	log.Fatal(appruntime.ServeGRPC(address, func(server *grpc.Server) { blogv1.RegisterUserServiceServer(server, service) }, grpc.UnaryInterceptor(tracekit.UnaryServerInterceptor("user-svc", nil))))
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
