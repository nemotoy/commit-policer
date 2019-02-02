package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nemotoy/commit-policer/handler"

	"github.com/google/go-github/github"
)

var (
	userName             = "nemotoy"
	jst                  = time.FixedZone("Asia/Tokyo", 9*60*60)
	dayHour              = 24 * time.Hour
	commitcount          int
	warningRateRemaining = 10
	interval             = 5 * time.Minute
)

type client struct {
	GitCli  *github.Client
	HttpCli *http.Client
}

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

	httpCli := &http.Client{}

	client := &client{
		GitCli:  github.NewClient(nil),
		HttpCli: httpCli,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go client.Watcher(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, fmt.Sprintf("pong"))
	})

	webhook := handler.NewWebhookHandler()
	mux.Handle("/webhook", webhook)

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) Watcher(ctx context.Context) {

	t := time.NewTicker(interval)

	for {
		select {
		case <-t.C:
			err := c.CommitSend()
			if err != nil {
				log.Fatal(err)
				return
			}
		case <-ctx.Done():
			return
		}
	}

}

func (c *client) CommitSend() error {

	now := time.Now()

	events, resp, err := c.GitCli.Activity.ListEventsPerformedByUser(context.Background(), userName, true, nil)
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
		// TODO: This calculation is not daily strictly. `within 24 hours`.
		// log.Printf("Now: %v, Event: %v, Dur: %v, Daily: %v\n", now, eventTime, dur, dayHour-dur)

		if dayHour-dur >= 0 {
			commitcount++
		}
	}

	fmt.Printf("%v commits within 24 hours\n", commitcount)

	body, err := json.Marshal(events)
	if err != nil {
		fmt.Errorf("Failed to marshal. error: %v", err)
	}

	req, err := http.NewRequest("POST", "https://hooks.slack.com/services/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", bytes.NewBuffer(body))
	if err != nil {
		fmt.Errorf("Failed to create request. error: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	httpResp, err := c.HttpCli.Do(req)
	if err != nil {
		fmt.Errorf("Failed to request. error: %v", err)
	}

	fmt.Printf("Success to request.", httpResp)

	if commitcount == 0 {
		fmt.Println("TODO: Remind!!!")
		// TODO: send to LINE
	}

	return nil
}
