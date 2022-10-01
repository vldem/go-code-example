package notifier

import (
	"ur-services/spv-notif/internal/config"
	botPkg "ur-services/spv-notif/internal/pkg/bot"
	"ur-services/spv-notif/internal/pkg/mailer"
	lmPkg "ur-services/spv-notif/internal/pkg/mailer/local"
	"ur-services/spv-notif/internal/pkg/models"
)

type Notifier struct {
	bot    botPkg.TelBot
	mailer mailer.Mailer
}

func New() *Notifier {
	var bot botPkg.TelBot
	var mailer mailer.Mailer

	if config.AppConfig.Telegram.ChatId != 0 {
		bot = botPkg.MustNew()
	}

	if config.AppConfig.Email.Prog.Cmd != "" {
		mailer = lmPkg.New()
	}

	return &Notifier{
		bot:    bot,
		mailer: mailer,
	}

}

func (n *Notifier) Notify(header string, payload string) error {

	if n.mailer != nil {
		mail := &models.Mail{
			From:    config.AppConfig.Email.From,
			To:      config.AppConfig.Email.To,
			Subject: config.AppConfig.Email.Subject,
		}

		mail.TextBody.Write([]byte(header))
		mail.TextBody.Write([]byte(payload))
		if err := n.mailer.SendMail(*mail); err != nil {
			return err
		}
	}

	if n.bot != nil {
		if err := n.bot.SendMessage(header + payload); err != nil {
			return err
		}
	}

	return nil
}
