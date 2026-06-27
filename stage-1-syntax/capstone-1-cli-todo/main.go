package main

import (
	"fmt"
	"os"

	"just-go/stage-1-syntax/capstone-1-cli-todo/app"
)

func main() {
	if err := app.Run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
