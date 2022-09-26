package listener

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"ur-services/spv-notif/internal/config"
	botPkg "ur-services/spv-notif/internal/pkg/bot"
	"ur-services/spv-notif/internal/pkg/notifier"
)

func Listen() {
	var bot botPkg.TelBot

	if config.AppConfig.Telegram.ChatId != 0 {
		bot = botPkg.MustNew()
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		ready()

		header, err := readHeader(reader)
		if err != nil {
			failure(err)
			continue
		}

		payload, err := readPayload(reader, getPayloadLen(header))
		if err != nil {
			failure(err)
			continue
		}

		if err := notifier.Notify(bot, header, payload); err != nil {
			failure(err)
			continue
		}

		success()
	}
}

func ready() {
	fmt.Fprint(os.Stdout, "READY\n")
}

func success() {
	fmt.Fprint(os.Stdout, "RESULT 2\nOK")
}

func failure(err error) {
	fmt.Fprintln(os.Stderr, err)
	fmt.Fprint(os.Stdout, "RESULT 4\nFAIL")
}

func readHeader(reader *bufio.Reader) (string, error) {
	data, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return data, nil
}

func readPayload(reader *bufio.Reader, payloadLen int) ([]byte, error) {
	buf := make([]byte, payloadLen)
	_, err := reader.Read(buf)
	if err != nil {
		return buf, err
	}

	return buf, nil
}

func getPayloadLen(header string) int {
	header = strings.TrimSpace(header)

	items := strings.Split(header, " ")
	if len(items) == 0 {
		return 0
	}

	length := strings.Split(items[len(items)-1], ":")
	if len(length) < 2 {
		return 0
	}

	l, _ := strconv.Atoi(length[1])
	return l
}
