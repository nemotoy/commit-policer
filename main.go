package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
)

const (
	headerRateLimit     = "X-RateLimit-Limit"
	headerRateRemaining = "X-RateLimit-Remaining"
	headerRateReset     = "X-RateLimit-Reset"
)

var (
	userName = "nemotoy"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {

	// use package 'go-github'
	client := github.NewClient(nil)

	event, res, err := client.Activity.ListEventsPerformedByUser(context.Background(), userName, true, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Event: %v\n", event[0])
	fmt.Printf("%v\n", res.Header.Get(headerRateRemaining))

	return nil
}
