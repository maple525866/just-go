package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"just-go/stage-2-business/capstone-2-blog-api/auth"
	"just-go/stage-2-business/capstone-2-blog-api/cache"
	"just-go/stage-2-business/capstone-2-blog-api/server"
	"just-go/stage-2-business/capstone-2-blog-api/store"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	fmt.Println("Capstone 2 Blog API")
	fmt.Println("routes: POST /api/register, POST /api/login, GET/POST /api/articles, GET/PUT/DELETE /api/articles/{id}, POST /api/articles/{id}/comments")
	fmt.Println("observability: GET /livez, GET /readyz, GET /metrics")
	fmt.Printf("listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, buildHandler()))
}

func buildHandler() http.Handler {
	return server.NewAPI(store.NewMemoryStore(), cache.NewArticleCache(time.Minute), auth.NewTokenManager([]byte("dev-secret")))
}
