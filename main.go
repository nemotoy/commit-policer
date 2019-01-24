package main

import (
	"context"
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
	interval             = 1 * time.Hour
)

type client struct {
	GitCli *github.Client
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

	client := &client{
		GitCli: github.NewClient(nil),
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
		log.Printf("Now: %v, Event: %v, Dur: %v, Daily: %v\n", now, eventTime, dur, dayHour-dur)

		if dayHour-dur >= 0 {
			fmt.Printf("Event: %v, Dur: %v, Daily: %v\n", eventTime, dur, dayHour-dur)
			commitcount++
		}
	}

	fmt.Printf("%v commits within 24 hours\n", commitcount)

	if commitcount == 0 {
		fmt.Println("TODO: Remind!!!")
		// TODO: send to LINE
	}

	return nil
}
