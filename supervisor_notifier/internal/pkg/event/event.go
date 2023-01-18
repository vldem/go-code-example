package event

import (
	"strconv"
	"strings"

	"github.com/vldem/go-code-example/supervisor_notifier/internal/pkg/models"
)

func GetHeader(rawData string) *models.EventHeader {
	rawData = strings.TrimSpace(rawData)
	header := &models.EventHeader{}

	items := strings.Split(rawData, " ")
	if len(items) == 0 {
		return header
	}

	for _, item := range items {
		params := strings.Split(item, ":")
		if len(params) < 2 {
			continue
		}
		switch params[0] {
		case "ver":
			header.Ver = params[1]
		case "server":
			header.Server = params[1]
		case "serial":
			val, _ := strconv.ParseUint(params[1], 10, 64)
			header.Serial = uint(val)
		case "pool":
			header.Pool = params[1]
		case "poolserial":
			val, _ := strconv.ParseUint(params[1], 10, 64)
			header.PoolSerial = uint(val)
		case "eventname":
			header.EventName = params[1]
		case "len":
			val, _ := strconv.ParseUint(params[1], 10, 64)
			header.Len = int(val)
		}
	}
	return header
}

func GetPayload(rawData string) *models.EventPayload {
	rawData = strings.TrimSpace(rawData)
	payload := &models.EventPayload{}

	items := strings.Split(rawData, " ")
	if len(items) == 0 {
		return payload
	}

	for _, item := range items {
		params := strings.Split(item, ":")
		if len(params) < 2 {
			continue
		}
		switch params[0] {
		case "processname":
			payload.ProcessName = params[1]
		case "groupname":
			payload.GroupName = params[1]
		case "from_state":
			payload.FromState = params[1]
		case "expected":
			val, _ := strconv.ParseUint(params[1], 10, 64)
			payload.Expected = uint8(val)
		case "pid":
			val, _ := strconv.ParseUint(params[1], 10, 64)
			payload.Pid = uint(val)
		}
	}
	return payload
}
