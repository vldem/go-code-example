package notifier

import (
	"github.com/vldem/go-code-example/golib/pkg/mailer"
	lmPkg "github.com/vldem/go-code-example/golib/pkg/mailer/local"
	"github.com/vldem/go-code-example/supervisor_notifier/internal/config"
	botPkg "github.com/vldem/go-code-example/supervisor_notifier/internal/pkg/bot"
)

type Notifier struct {
	bot    botPkg.TelBot
	mailer mailer.MailerInterface
}

func New() *Notifier {
	var bot botPkg.TelBot
	var mailer mailer.MailerInterface

	if config.AppConfig.Telegram.ChatId != 0 {
		bot = botPkg.MustNew()
	}

	if config.AppConfig.Email.Prog.Cmd != "" {
		mailer = lmPkg.New(config.AppConfig.Email.Prog.Cmd, config.AppConfig.Email.Prog.Args)
	}

	return &Notifier{
		bot:    bot,
		mailer: mailer,
	}

}

func (n *Notifier) Notify(header string, payload string) error {

	if n.mailer != nil {
		mail := &mailer.Mail{
			From:    config.AppConfig.Email.From,
			To:      config.AppConfig.Email.To,
			Subject: config.AppConfig.Email.Subject,
		}

		mail.TextBody.Write([]byte(header))
		mail.TextBody.Write([]byte(payload))
		if _, err := n.mailer.SendMail(*mail); err != nil {
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
