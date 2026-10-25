package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kiril755/alert-nearby-bot/internal/config"
	"github.com/kiril755/alert-nearby-bot/internal/matcher"
	"github.com/kiril755/alert-nearby-bot/internal/notifier"
	"github.com/kiril755/alert-nearby-bot/internal/watcher"
)

func main() {
	log.Println("starting alert bot")

	if err := godotenv.Load(); err != nil {
		log.Fatalf("could not load .env file: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	notify, err := notifier.New(cfg.BotToken, cfg.ChatIDs)
	if err != nil {
		log.Fatalf("notifier: %v", err)
	}

	match := matcher.New(cfg.Keywords)

	w := watcher.New(watcher.Config{
		APIID: cfg.APIID,
		APIHash: cfg.APIHash,
		Phone: cfg.Phone,
		Channels: cfg.Channels,
		SessionFile: cfg.SessionFile,
		PollInterval: cfg.PollInterval,
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	handler := func(m watcher.Message) {
		if !match.Match(m.Text) {
			return
		}
		text := fmt.Sprintf("[%s] %s (%s)\n%s", m.Date.Format("02.01.2006 15:04"), m.ChannelTitle, m.Channel, m.Text)
		if err := notify.Send(text); err != nil {
			log.Printf("send error: %v", err)
		}
	}

	if err := w.Run(ctx, handler); err != nil && ctx.Err() == nil {
		log.Fatalf("watcher: %v", err)
	}
	log.Println("shutting down")
}
