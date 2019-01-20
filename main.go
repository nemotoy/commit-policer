package main

import (
	"context"
	"fmt"
	"log"
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
	userName             = "nemotoy"
	jst                  = time.FixedZone("Asia/Tokyo", 9*60*60)
	dayHour              = 24 * time.Hour
	commitcount          int
	warningRateRemaining = 10
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// // TODO: new goroutine
	now := time.Now()

	client := github.NewClient(nil)

	events, resp, err := client.Activity.ListEventsPerformedByUser(context.Background(), userName, true, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println(resp.Rate)

	if resp.Rate.Remaining <= warningRateRemaining {
		fmt.Printf("Rate remaining is warn %v", resp.Rate.Remaining)
	}

	for _, event := range events {

		eventTime := event.CreatedAt.In(jst)
		dur := now.Sub(eventTime)
		// TODO: This calculation is not daily strictly. `within 24 hours`
		if dayHour-dur >= 0 {
			fmt.Printf("Event: %v, Dur: %v, Daily: %v\n", eventTime, dur, dayHour-dur)
			commitcount++
		}
	}

	fmt.Printf("%v commits within 24 hours\n", commitcount)

	if commitcount == 0 {
		fmt.Println("TODO: Remind!!!")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, fmt.Sprintf("pong"))
	})

	webhook := NewWebhookHandler()
	mux.Handle("/webhook", webhook)

	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		return err
	}

	return nil
}

// WebhookHandler ...
type WebhookHandler struct {
}

// NewWebhookHandler ...
func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, fmt.Sprintf("webhook"))
}
