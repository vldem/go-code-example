package listener

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"ur-services/spv-notif/internal/pkg/event"
	notifPkg "ur-services/spv-notif/internal/pkg/notifier"
)

type SupervisorListener struct {
	programs map[string]time.Time
}

func New() *SupervisorListener {

	return &SupervisorListener{
		programs: make(map[string]time.Time),
	}

}

func (l *SupervisorListener) Listen() {

	notifier := notifPkg.New()
	reader := bufio.NewReader(os.Stdin)

	for {
		l.ready()

		headerString, err := l.readHeader(reader)
		if err != nil {
			l.failure(err)
			continue
		}
		header := event.GetHeader(headerString)

		payloadString, err := l.readPayload(reader, header.Len)
		if err != nil {
			l.failure(err)
			continue
		}
		payload := event.GetPayload(payloadString)

		// send notification once per 5 minutes in order to avoid huge amount of the same notifications
		if time.Since(l.programs[payload.ProcessName]) > time.Minute*5 {
			l.programs[payload.ProcessName] = time.Now()

			if err := notifier.Notify(headerString, payloadString); err != nil {
				l.failure(err)
				continue
			}
		}

		l.success()
	}
}

func (l SupervisorListener) ready() {
	fmt.Fprint(os.Stdout, "READY\n")
}

func (l SupervisorListener) success() {
	fmt.Fprint(os.Stdout, "RESULT 2\nOK")
}

func (l SupervisorListener) failure(err error) {
	fmt.Fprintln(os.Stderr, err)
	fmt.Fprint(os.Stdout, "RESULT 4\nFAIL")
}

func (l SupervisorListener) readHeader(reader *bufio.Reader) (string, error) {
	data, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return data, nil
}

func (l SupervisorListener) readPayload(reader *bufio.Reader, payloadLen int) (string, error) {
	buf := make([]byte, payloadLen)
	_, err := reader.Read(buf)
	if err != nil {
		return string(buf), err
	}

	return string(buf), nil
}
