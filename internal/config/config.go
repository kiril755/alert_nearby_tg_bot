package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	APIID int
	APIHash string
	Phone string
	BotToken string
	ChatIDs []int64
	Channels []string
	Keywords []string
	PollInterval time.Duration
	SessionFile string
}

func Load() (*Config, error) {
	cfg := &Config{
		SessionFile: getEnv("SESSION_FILE", "session.json"),
		PollInterval: 10 * time.Second,
	}

	var err error
	if cfg.APIID, err = strconv.Atoi(os.Getenv("API_ID")); err != nil {
		return nil, fmt.Errorf("API_ID must be a number: %w", err)
	}
	cfg.APIHash = os.Getenv("API_HASH")
	cfg.Phone = os.Getenv("PHONE")
	cfg.BotToken = os.Getenv("BOT_TOKEN")

	// Chat IDs: comma-separated list of numbers.
	for _, s := range splitList(os.Getenv("MY_CHAT_IDS")) {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("MY_CHAT_IDS: %q is not a number: %w", s, err)
		}
		cfg.ChatIDs = append(cfg.ChatIDs, id)
	}

	// Channels: comma-separated usernames (without @).
	cfg.Channels = splitList(os.Getenv("CHANNELS"))

	// Keywords: comma-separated, lowercased for case-insensitive matching.
	for _, kw := range splitList(os.Getenv("KEYWORDS")) {
		cfg.Keywords = append(cfg.Keywords, strings.ToLower(kw))
	}

	if s := os.Getenv("POLL_INTERVAL_SECONDS"); s != "" {
		secs, err := strconv.Atoi(s)
		if err != nil || secs <= 0 {
			return nil, fmt.Errorf("POLL_INTERVAL_SECONDS must be a positive number, got %q", s)
		}
		cfg.PollInterval = time.Duration(secs) * time.Second
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	switch {
	case c.APIHash == "":
		return fmt.Errorf("API_HASH is required")
	case c.Phone == "":
		return fmt.Errorf("PHONE is required")
	case c.BotToken == "":
		return fmt.Errorf("BOT_TOKEN is required")
	case len(c.ChatIDs) == 0:
		return fmt.Errorf("MY_CHAT_IDS is required (comma-separated)")
	case len(c.Channels) == 0:
		return fmt.Errorf("CHANNELS is required (comma-separated)")
	case len(c.Keywords) == 0:
		return fmt.Errorf("KEYWORDS is required (comma-separated)")
	}
	return nil
}

// splitList splits a comma-separated value, trimming spaces and dropping blanks.
func splitList(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
