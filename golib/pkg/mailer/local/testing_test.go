package local

import (
	"testing"

	"github.com/vldem/go-code-example/golib/pkg/mailer"
)

type mailerFixture struct {
	boundaryId     string
	from           string
	to             string
	cc             string
	bcc            string
	subject        string
	replayTo       string
	organization   string
	unsubscribeUrl string
	mailText       string
	mailHTML       string
	customHeader   []mailer.MailCustomHeader
	expectedResult string
}

func setUp(t *testing.T) mailerFixture {
	var fixture mailerFixture
	fixture.boundaryId = "As1ls3ss02KF"
	fixture.from = "sender@dummy.com"
	fixture.to = "receiver@dummy.com"
	fixture.cc = "cc@dummy.com"
	fixture.bcc = "bcc@dummy.com"
	fixture.subject = "test mail"
	fixture.replayTo = "receiver@dummy.com"
	fixture.organization = "test_org"
	fixture.unsubscribeUrl = "https://dummy.com/unsubscribe/"
	fixture.mailText = "Text mail"
	fixture.mailHTML = "<html><body><p>HTML text in test mail</p></body></html>"
	fixture.customHeader = []mailer.MailCustomHeader{
		{
			Key:   "custom_header1",
			Value: "test custom header value 1",
		},
		{
			Key:   "custom_header2",
			Value: "test custom header value 2",
		},
	}
	fixture.expectedResult = `Reply-To: receiver@dummy.com
Organization: test_org
X-Priority: 3 (Normal)
From: sender@dummy.com
To: receiver@dummy.com
Subject: test mail
Cc: cc@dummy.com
Bcc: bcc@dummy.com
List-Unsubscribe: <https://dummy.com/unsubscribe/receiver@dummy.com>
custom_header1: test custom header value 1
custom_header2: test custom header value 2
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary="As1ls3ss02KF"

--As1ls3ss02KF
Content-Type: text/plain; charset=utf-8
Content-Transfer-Encoding: 7bit

Text mail
--As1ls3ss02KF
Content-Type: text/html; charset=utf-8
Content-Transfer-Encoding: 8bit

<html><body><p>HTML text in test mail</p></body></html>`

	return fixture

}
