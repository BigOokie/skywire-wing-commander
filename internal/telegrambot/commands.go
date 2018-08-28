// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package telegrambot

// Command struct is used to define a Telegram Bot command, including
// if its an Admin only command, the string command (i.e. `/start`) and
// the function that will handle the command
type Command struct {
	Admin       bool
	Command     string
	Handlerfunc CommandHandler
}

// Commands provides an array (slice) of Command structs
type Commands []Command

func (bot *Bot) setCommandHandlers() {
	for _, command := range commands {
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
		false,
		"stop",
		(*Bot).handleCommandStop,
	},
	Command{
		false,
		"status",
		(*Bot).handleCommandStatus,
	},
	Command{
		false,
		"showconfig",
		(*Bot).handleCommandShowConfig,
	},
	Command{
		false,
		"checkupdate",
		(*Bot).handleCommandCheckUpdate,
	},
	Command{
		false,
		"listnodes",
		(*Bot).handleCommandListNodes,
	},
}
