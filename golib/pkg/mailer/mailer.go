package mailer

import "bytes"

const MailPriorityHight = "1 (Highest)"
const MailPriorityNormal = "3 (Normal)"
const MailPriorityLow = "5 (Lowest)"

type Mail struct {
	From           string
	To             string
	Cc             string
	Bcc            string
	Subject        string
	ReplayTo       string
	Organization   string
	Priority       string
	UnsubscribeUrl string
	ContentType    int
	HtmlBody       bytes.Buffer
	TextBody       bytes.Buffer
	CustomHeader   []MailCustomHeader
}

type MailCustomHeader struct {
	Key   string
	Value string
}

type MailerInterface interface {
	SendMail(mail Mail) ([]byte, error)
}
