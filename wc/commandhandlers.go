package wingcommander

// Handler for help command
func (bot *Bot) handleCommandHelp(ctx *Context, command, args string) error {
	return bot.Reply(ctx, msgHelp)
}

// Handler for about command
func (bot *Bot) handleCommandAbout(ctx *Context, command, args string) error {
	return bot.Reply(ctx, msgAbout)
}

// Handler for start command
func (bot *Bot) handleCommandStart(ctx *Context, command, args string) error {
	return nil
}

// Handler for stop command
func (bot *Bot) handleCommandStop(ctx *Context, command, args string) error {
	return nil
}

// Handler for status command
func (bot *Bot) handleCommandStatus(ctx *Context, command, args string) error {
	return nil
}

// Handler for heartbeat command
func (bot *Bot) handleCommandHeartBeat(ctx *Context, command, args string) error {
	return nil
}

func (bot *Bot) handleDirectMessageFallback(ctx *Context, text string) (bool, error) {
	return true, bot.Reply(ctx, "unknown command.")
}

func (bot *Bot) AddPrivateMessageHandler(handler MessageHandler) {
	bot.privateMessageHandlers = append(bot.privateMessageHandlers, handler)
}

func (bot *Bot) AddGroupMessageHandler(handler MessageHandler) {
	bot.groupMessageHandlers = append(bot.groupMessageHandlers, handler)
}
