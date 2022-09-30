package local

import (
	"io"
	"os/exec"
	"ur-services/spv-notif/internal/config"
	"ur-services/spv-notif/internal/pkg/mailer"
	"ur-services/spv-notif/internal/pkg/models"

	"github.com/pkg/errors"
)

type localMailerImplementation struct {
	Cmd  string
	Args []string
}

func New() mailer.Mailer {
	return &localMailerImplementation{
		Cmd:  config.AppConfig.Email.Prog.Cmd,
		Args: config.AppConfig.Email.Prog.Args,
	}
}

func (m *localMailerImplementation) SendMail(mail models.Mail) error {
	sender := exec.Command(m.Cmd, m.Args...)
	stdin, err := sender.StdinPipe()
	if err != nil {
		return errors.Wrap(err, "opening stdin for mail prog")
	}

	io.WriteString(stdin, "From: "+mail.From+"\n")
	io.WriteString(stdin, "To: "+mail.To+"\n")
	io.WriteString(stdin, "Subject: "+mail.Subject+"\n")
	for _, cHeader := range mail.CustomHeader {
		io.WriteString(stdin, cHeader.Cmd+": "+cHeader.Value+"\n")
	}

	if mail.TextBody.Len() > 0 {
		io.WriteString(stdin, "MIME-Version: 1.0\n")
		io.WriteString(stdin, "Content-Type: text/plain; charset=utf-8\n")
		io.WriteString(stdin, "Content-Transfer-Encoding: 7bit\n\n")
		io.WriteString(stdin, mail.TextBody.String()+"\n")
	}

	err = stdin.Close()
	if err != nil {
		return errors.Wrap(err, "closing stdin of mail prog")
	}

	err = sender.Run()
	if err != nil {
		return errors.Wrap(err, "executing mail prog")
	}

	return nil
}
