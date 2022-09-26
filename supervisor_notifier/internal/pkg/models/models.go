package models

import "bytes"

type Mail struct {
	From         string
	To           string
	Cc           string
	Bcc          string
	Subject      string
	ReplayTo     string
	ContentType  int
	HtmlBody     bytes.Buffer
	TextBody     bytes.Buffer
	CustomHeader []CustomHeader
}

type CustomHeader struct {
	Cmd   string
	Value string
}
