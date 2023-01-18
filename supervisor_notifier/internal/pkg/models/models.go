package models

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
