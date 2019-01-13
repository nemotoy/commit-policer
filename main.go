package main

import (
	"fmt"
	"net/http"
	"os"
)

const (
	headerRateLimit     = "X-RateLimit-Limit"
	headerRateRemaining = "X-RateLimit-Remaining"
	headerRateReset     = "X-RateLimit-Reset"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {

	client := &http.Client{}

	res, err := client.Get("https://api.github.com/users/nemotoy/events")
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", res.Header.Get(headerRateRemaining))

	return nil
}
