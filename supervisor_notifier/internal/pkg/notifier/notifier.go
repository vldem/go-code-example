package notifier

import (
	"ur-services/spv-notif/internal/config"
	"ur-services/spv-notif/internal/pkg/bot"
	"ur-services/spv-notif/internal/pkg/mailsender"
	"ur-services/spv-notif/internal/pkg/models"
)

func Notify(bot bot.TelBot, header string, payload []byte) error {

	if config.AppConfig.Email.Prog.Cmd != "" {
		mail := &models.Mail{
			From:    config.AppConfig.Email.From,
			To:      config.AppConfig.Email.To,
			Subject: config.AppConfig.Email.Subject,
		}

		mail.TextBody.Write([]byte(header))
		mail.TextBody.Write(payload)
		if err := mailsender.SendMail(*mail); err != nil {
			return err
		}
	}

	if bot != nil {
		if err := bot.SendMessage(header + string(payload)); err != nil {
			return err
		}
	}

	return nil
}
