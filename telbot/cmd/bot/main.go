// This is Demidov Vladislav's telegram bot
package main

import (
	"log"

	"github.com/vldem/go-code-example/telbot/internal/commander"
	"github.com/vldem/go-code-example/telbot/internal/controller"
)

func main() {
	// log.Println(os.Getwd())
	log.Println("start main")
	cmd, err := commander.Init()
	if err != nil {
		log.Panic(err)
	}
	controller.AddControllers(cmd)

	if err := cmd.Run(); err != nil {
		log.Panic(err)
	}
}
