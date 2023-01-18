package main

import (
	"log"

	"github.com/vldem/go-code-example/supervisor_notifier/internal/config"
	listenerPgk "github.com/vldem/go-code-example/supervisor_notifier/internal/pkg/listener"
)

func main() {
	listener := listenerPgk.New()
	listener.Listen()

}

func init() {
	if err := config.ReadConfig(); err != nil {
		log.Fatal(err)
	}
}
