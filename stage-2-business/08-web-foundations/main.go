package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"just-go/stage-2-business/08-web-foundations/server"
	"just-go/stage-2-business/08-web-foundations/store"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	handler := server.NewRouter(store.NewSeededMemoryStore())
	log.Printf("chapter 08 web foundations listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
