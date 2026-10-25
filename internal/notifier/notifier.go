package notifier

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Notifier sends alert messages to one or more chats through a Telegram bot.
type Notifier struct {
	bot     *tgbotapi.BotAPI
	chatIDs []int64
}

// New creates a Notifier from a bot token and target chat IDs.
func New(token string, chatIDs []int64) (*Notifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("invalid BOT_TOKEN - Telegram rejected it (%w)", err)
	}
	return &Notifier{bot: bot, chatIDs: chatIDs}, nil
}

// Send delivers a text message to every configured chat. It attempts all chats
// and returns the first error encountered, if any.
func (n *Notifier) Send(text string) error {
	var firstErr error
	for _, id := range n.chatIDs {
		if _, err := n.bot.Send(tgbotapi.NewMessage(id, text)); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}