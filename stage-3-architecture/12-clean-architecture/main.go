package main

import (
	"log"
	"net/http"
)

func main() {
	handler := initializeHandler()
	log.Println("clean architecture demo listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
