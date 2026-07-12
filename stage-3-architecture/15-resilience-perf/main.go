package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := RunDemo(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "chapter 15 demo failed: %v\n", err)
		os.Exit(1)
	}
}
