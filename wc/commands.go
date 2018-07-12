package wingcommander

type Command struct {
	Admin       bool
	Command     string
	Handlerfunc CommandHandler
}

type Commands []Command

func (bot *Bot) setCommandHandlers() {
	for _, command := range commands {
		//bot.SetCommandHandler(command.Admin, command.Command, command.Handlerfunc)
		if command.Admin {
			bot.adminCommandHandlers[command.Command] = command.Handlerfunc
		} else {
			bot.commandHandlers[command.Command] = command.Handlerfunc
		}
	}

	bot.AddPrivateMessageHandler((*Bot).handleDirectMessageFallback)
	bot.AddGroupMessageHandler((*Bot).handleDirectMessageFallback)
}

var commands = Commands{
	Command{
		false,
		"help",
		(*Bot).handleCommandHelp,
	},
	Command{
		false,
		"about",
		(*Bot).handleCommandAbout,
	},
	Command{
		false,
		"start",
		(*Bot).handleCommandStart,
	},
	Command{
		true,
		"stop",
		(*Bot).handleCommandStop,
	},
	Command{
		true,
		"status",
		(*Bot).handleCommandStatus,
	},
	Command{
		true,
		"heartbeat",
		(*Bot).handleCommandHeartBeat,
	},
}
