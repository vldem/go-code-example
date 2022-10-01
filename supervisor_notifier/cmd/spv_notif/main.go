package main

import (
	"log"
	"ur-services/spv-notif/internal/config"
	listenerPgk "ur-services/spv-notif/internal/pkg/listener"
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
