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

type EventHeader struct {
	Ver        string
	Server     string
	Pool       string
	EventName  string
	Serial     uint
	PoolSerial uint
	Len        int
}

type EventPayload struct {
	ProcessName string
	GroupName   string
	FromState   string
	Expected    uint8
	Pid         uint
}
