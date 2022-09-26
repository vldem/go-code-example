package main

import (
	"log"
	"ur-services/spv-notif/internal/config"
	"ur-services/spv-notif/internal/pkg/listener"
)

func main() {

	listener.Listen()

}

func init() {
	if err := config.ReadConfig(); err != nil {
		log.Fatal(err)
	}
}
