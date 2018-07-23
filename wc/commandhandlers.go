package wingcommander

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Handler for help command
func (bot *Bot) handleCommandHelp(ctx *Context, command, args string) error {
	log.Debug("Handle command /help")
	return bot.Send(ctx, "whisper", "markdown", msgHelp)
}

// Handler for about command
func (bot *Bot) handleCommandAbout(ctx *Context, command, args string) error {
	log.Debug("Handle command /about")
	return bot.Send(ctx, "whisper", "markdown", msgAbout)
}

// Handler for start command
func (bot *Bot) handleCommandStart(ctx *Context, command, args string) error {
	log.Debug("Handle command /start")

	return nil
}

// Handler for stop command
func (bot *Bot) handleCommandStop(ctx *Context, command, args string) error {
	log.Debug("Handle command /stop")
	return nil
}

// Handler for status command
func (bot *Bot) handleCommandStatus(ctx *Context, command, args string) error {
	log.Debug("Handle command /status")
	return bot.Send(ctx, "whisper", "markdown", msgStatus)
}

// Handler for heartbeat command
func (bot *Bot) handleCommandHeartBeat(ctx *Context, command, args string) error {
	log.Debug("Handle command /heartbeat")
	go botHeartbeatLoop(bot, ctx)
	return nil
}

func (bot *Bot) handleDirectMessageFallback(ctx *Context, text string) (bool, error) {
	log.Debug("Ignoring unknown command: %s", text)
	return true, bot.Reply(ctx, "Unknown command.")
}

func (bot *Bot) AddPrivateMessageHandler(handler MessageHandler) {
	bot.privateMessageHandlers = append(bot.privateMessageHandlers, handler)
}

func (bot *Bot) AddGroupMessageHandler(handler MessageHandler) {
	bot.groupMessageHandlers = append(bot.groupMessageHandlers, handler)
}

// Bot Heartbeat Loop
func botHeartbeatLoop(bot *Bot, ctx *Context) {
	ticker := time.NewTicker(time.Hour * 1)

	for {
		select {
		case <-ticker.C:
			bot.Send(ctx, "whisper", "markdown", msgStatus)
		}
	}
}
