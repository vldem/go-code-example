package bot

import (
	"log"

	"ur-services/spv-notif/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type TelBot interface {
	SendMessage(message string) error
}

func MustNew() TelBot {
	bot, err := tgbotapi.NewBotAPI(config.AppConfig.Telegram.BotKey)
	if err != nil {
		log.Panic(errors.Wrap(err, "init tgbot"))
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &commander{
		bot: bot,
	}
}

type commander struct {
	bot *tgbotapi.BotAPI
}

func (c *commander) SendMessage(message string) error {
	msg := tgbotapi.NewMessage(config.AppConfig.Telegram.ChatId, "")
	msg.Text = message
	_, err := c.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "send tg message")
	}
	return nil
}
