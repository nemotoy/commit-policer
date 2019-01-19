package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/github"
)

const (
	headerRateLimit     = "X-RateLimit-Limit"
	headerRateRemaining = "X-RateLimit-Remaining"
	headerRateReset     = "X-RateLimit-Reset"
)

var (
	userName                           = "nemotoy"
	jst                                = time.FixedZone("Asia/Tokyo", 9*60*60)
	dayHour              time.Duration = 24
	commitcount          int
	warningRateRemaining = 10
	port                 = ":9000"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {

	// TODO: new goroutine
	now := time.Now()

	client := github.NewClient(nil)

	events, resp, err := client.Activity.ListEventsPerformedByUser(context.Background(), userName, true, nil)
	if err != nil {
		return err
	}

	fmt.Println(events, resp.Rate)

	if resp.Rate.Remaining <= warningRateRemaining {
		fmt.Printf("Rate remaining is warn %v", resp.Rate.Remaining)
	}

	for _, event := range events {

		eTime := event.CreatedAt.In(jst)
		dur := now.Sub(eTime)
		if dur < dayHour {
			commitcount++
		}
	}

	if commitcount == 0 {
		fmt.Println("TODO: Remind!!!")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, fmt.Sprintf("pong"))
	})
	http.ListenAndServe(port, nil)

	return nil
}
