package main

import (
	"fmt"
	"time"

	"just-go/stage-2-business/capstone-2-blog-api/auth"
	"just-go/stage-2-business/capstone-2-blog-api/cache"
	"just-go/stage-2-business/capstone-2-blog-api/server"
	"just-go/stage-2-business/capstone-2-blog-api/store"
)

func main() {
	_ = server.NewAPI(store.NewMemoryStore(), cache.NewArticleCache(time.Minute), auth.NewTokenManager([]byte("dev-secret")))
	fmt.Println("Capstone 2 Blog API")
	fmt.Println("routes: POST /api/register, POST /api/login, GET/POST /api/articles, POST /api/articles/{id}/comments")
	fmt.Println("observability: GET /livez, GET /readyz, GET /metrics")
}
