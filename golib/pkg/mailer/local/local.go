package local

import (
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/vldem/go-code-example/golib/pkg/mailer"
	"github.com/vldem/go-code-example/golib/pkg/utils"

	"github.com/pkg/errors"
)

type localMailerImplementation struct {
	Commander commander
	Utils     utils.UtilsInterface
}

func New(command string, args []string) mailer.MailerInterface {
	cmd := exec.Command(command, args...)
	stdin, err := cmd.StdinPipe()
	utils := utils.New()
	if err != nil {
		log.Fatal(err.Error())
	}

	return &localMailerImplementation{
		Commander: commander{
			Cmd:   cmd,
			stdin: stdin,
		},
		Utils: utils,
	}
}

type CommanderInterface interface {
	CombinedOutput(string, ...string) ([]byte, error)
}

type commander struct {
	Cmd   *exec.Cmd
	stdin io.WriteCloser
}

func (c commander) CombinedOutput(command string, args ...string) ([]byte, error) {
	return c.Cmd.CombinedOutput()
}

func (m *localMailerImplementation) SendMail(mail mailer.Mail) ([]byte, error) {
	sender := m.Commander.Cmd
	stdin := m.Commander.stdin

	if mail.ReplayTo != "" {
		io.WriteString(stdin, fmt.Sprintf("Reply-To: %s\n", mail.ReplayTo))
	}

	if mail.Organization != "" {
		io.WriteString(stdin, fmt.Sprintf("Organization: %s\n", mail.Organization))
	}

	if mail.Priority == "" {
		mail.Priority = mailer.MailPriorityNormal
	}

	io.WriteString(stdin, fmt.Sprintf("X-Priority: %s\n", mail.Priority))
	io.WriteString(stdin, fmt.Sprintf("From: %s\n", mail.From))
	io.WriteString(stdin, fmt.Sprintf("To: %s\n", mail.To))
	io.WriteString(stdin, fmt.Sprintf("Subject: %s\n", mail.Subject))

	if mail.Cc != "" {
		io.WriteString(stdin, fmt.Sprintf("Cc: %s\n", mail.Cc))
	}

	if mail.Bcc != "" {
		io.WriteString(stdin, fmt.Sprintf("Bcc: %s\n", mail.Bcc))
	}

	if mail.UnsubscribeUrl != "" {
		io.WriteString(stdin, fmt.Sprintf("List-Unsubscribe: <%s%s>\n", mail.UnsubscribeUrl, mail.To))
	}

	for _, cHeader := range mail.CustomHeader {
		io.WriteString(stdin, fmt.Sprintf("%s: %s\n", cHeader.Key, cHeader.Value))
	}

	if mail.TextBody.Len() > 0 {
		boundaryId := m.Utils.GetRandIdString(12)
		io.WriteString(stdin, "MIME-Version: 1.0\n")
		io.WriteString(stdin, fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\n\n", boundaryId))
		io.WriteString(stdin, fmt.Sprintf("--%s\n", boundaryId))
		io.WriteString(stdin, "Content-Type: text/plain; charset=utf-8\n")
		io.WriteString(stdin, "Content-Transfer-Encoding: 7bit\n\n")
		io.WriteString(stdin, fmt.Sprintf("%s\n", mail.TextBody.String()))
		io.WriteString(stdin, fmt.Sprintf("--%s\n", boundaryId))
	}

	if mail.HtmlBody.Len() > 0 {
		io.WriteString(stdin, "Content-Type: text/html; charset=utf-8\n")
		io.WriteString(stdin, "Content-Transfer-Encoding: 8bit\n\n")
		io.WriteString(stdin, mail.HtmlBody.String())
	}

	err := stdin.Close()
	if err != nil {
		return nil, errors.Wrap(err, "[local mailer] closing stdin of mail prog")
	}

	out, err := sender.CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, "[local mailer] executing mail prog")
	}

	return out, nil
}
