package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hc100/tweet-bot/internal/config"
	"github.com/hc100/tweet-bot/internal/jobs"
	"github.com/hc100/tweet-bot/internal/scheduler"
	"github.com/hc100/tweet-bot/internal/xclient"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	client := xclient.NewClient(cfg.Credentials, 15*time.Second)
	loc := cfg.Location

	s := scheduler.New(
		client,
		loc,
		[]scheduler.Job{
			jobs.NewMorningMonyounJob(loc),
			jobs.NewCountdownTo2027Job(loc),
		},
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("tweet-bot started in timezone=%s", loc.String())
	if err := s.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("scheduler stopped: %v", err)
	}
}
