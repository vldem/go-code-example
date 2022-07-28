// This is Demidov Vladislav's telegram bot
package main

import (
	"log"

	botPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/bot"
	cmdAddPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/bot/command/add"
	cmdDeletePkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/bot/command/delete"
	cmdHelpPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/bot/command/help"
	cmdListPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/bot/command/list"
	cmdUpdatePkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/bot/command/update"
	userPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user"
)

func main() {
	var user userPkg.Interface
	{
		user = userPkg.New()
	}

	var bot botPkg.Interface
	{
		bot = botPkg.MustNew()

		commandAdd := cmdAddPkg.New(user)
		bot.RegisterHandler(commandAdd)

		commandUpdate := cmdUpdatePkg.New(user)
		bot.RegisterHandler(commandUpdate)

		commandDelete := cmdDeletePkg.New(user)
		bot.RegisterHandler(commandDelete)

		commandList := cmdListPkg.New(user)
		bot.RegisterHandler(commandList)

		commandHelp := cmdHelpPkg.New(map[string]string{
			commandAdd.Name():    commandAdd.Description(),
			commandUpdate.Name(): commandUpdate.Description(),
			commandDelete.Name(): commandDelete.Description(),
			commandList.Name():   commandList.Description(),
		})
		bot.RegisterHandler(commandHelp)
	}
	go runBot(bot)
	go runREST()
	runGRPCServer(user)
}

func runBot(bot botPkg.Interface) {
	if err := bot.Run(); err != nil {
		log.Panic(err)
	}
}
