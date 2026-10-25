package watcher

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

// Message is a new channel post surfaced to the caller.
type Message struct {
	Channel string
	ChannelTitle string
	Text string
	Date time.Time
}

type Config struct {
	APIID int
	APIHash string
	Phone string
	Channels []string
	SessionFile string
	PollInterval time.Duration
}

type Watcher struct {
	cfg Config
}

func New(cfg Config) *Watcher {
	return &Watcher{cfg: cfg}
}

// watched holds the resolved peer and last-seen message ID for one channel.
type watched struct {
	name string
	title string
	peer *tg.InputPeerChannel
	lastID int
}

// Run authenticates, resolves every channel, and polls until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context, onNew func(Message)) error {
	client := telegram.NewClient(w.cfg.APIID, w.cfg.APIHash, telegram.Options{
		SessionStorage: &telegram.FileSessionStorage{Path: w.cfg.SessionFile},
	})

	return client.Run(ctx, func(ctx context.Context) error {
		if err := w.authenticate(ctx, client); err != nil {
			return fmt.Errorf("authenticate: %w", err)
		}

		api := client.API()

		var channels []*watched
		for _, name := range w.cfg.Channels {
			peer, title, err := resolveChannel(ctx, api, name)
			if err != nil {
				return fmt.Errorf("resolve channel %q: %w", name, err)
			}
			channels = append(channels, &watched{name: name, title: title, peer: peer})
		}

		return w.poll(ctx, api, channels, onNew)
	})
}

func (w *Watcher) authenticate(ctx context.Context, client *telegram.Client) error {
	flow := auth.NewFlow(
		auth.Constant(w.cfg.Phone, "", auth.CodeAuthenticatorFunc(
			func(ctx context.Context, _ *tg.AuthSentCode) (string, error) {
				fmt.Print("Enter the code Telegram sent you: ")
				code, err := bufio.NewReader(os.Stdin).ReadString('\n')
				if err != nil {
					return "", fmt.Errorf("read code: %w", err)
				}
				return strings.TrimSpace(code), nil
			},
		)),
		auth.SendCodeOptions{},
	)
	if err := client.Auth().IfNecessary(ctx, flow); err != nil {
		return err
	}
	self, err := client.Self(ctx)
	if err != nil {
		return err
	}
	log.Printf("logged in as %s", self.FirstName)
	return nil
}

func resolveChannel(ctx context.Context, api *tg.Client, username string) (*tg.InputPeerChannel, string, error) {
	resolved, err := api.ContactsResolveUsername(ctx, &tg.ContactsResolveUsernameRequest{
		Username: username,
	})
	if err != nil {
		return nil, "", err
	}
	for _, chat := range resolved.Chats {
		if ch, ok := chat.(*tg.Channel); ok {
			peer := &tg.InputPeerChannel{ChannelID: ch.ID, AccessHash: ch.AccessHash}
			return peer, ch.Title, nil
		}
	}
	return nil, "", fmt.Errorf("not found")
}

func (w *Watcher) poll(ctx context.Context, api *tg.Client, channels []*watched, onNew func(Message)) error {
	// Baseline: record the newest existing message per channel so we don't
	// replay old history on startup.
	for _, c := range channels {
		msgs, err := fetchHistory(ctx, api, c.peer)
		if err != nil {
			return err
		}
		for _, m := range msgs {
			if m.ID > c.lastID {
				c.lastID = m.ID
			}
		}
	}
	log.Printf("watching %d channel(s) for new messages...", len(channels))

	ticker := time.NewTicker(w.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			for _, c := range channels {
				msgs, err := fetchHistory(ctx, api, c.peer)
				if err != nil {
					log.Printf("fetch error (%s): %v", c.name, err)
					continue
				}
				newMax := c.lastID
				for _, m := range msgs {
					if m.ID > c.lastID && m.Message != "" {
						onNew(Message{
							Channel:      c.name,
							ChannelTitle: c.title,
							Text:         m.Message,
							Date:         time.Unix(int64(m.Date), 0),
						})
					}
					if m.ID > newMax {
						newMax = m.ID
					}
				}
				c.lastID = newMax
			}
		}
	}
}

func fetchHistory(ctx context.Context, api *tg.Client, peer *tg.InputPeerChannel) ([]*tg.Message, error) {
	history, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
		Peer:  peer,
		Limit: 20,
	})
	if err != nil {
		return nil, err
	}
	channelMsgs, ok := history.(*tg.MessagesChannelMessages)
	if !ok {
		return nil, nil
	}
	var out []*tg.Message
	for _, m := range channelMsgs.Messages {
		if msg, ok := m.(*tg.Message); ok {
			out = append(out, msg)
		}
	}
	return out, nil
}
