package wingcommander

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

// Handler for help command
func (bot *Bot) handleCommandHelp(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /help")
	return bot.Send(ctx, "whisper", "markdown", msgHelp)
}

// Handler for about command
func (bot *Bot) handleCommandAbout(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /about")
	return bot.Send(ctx, "whisper", "markdown", msgAbout)
}

// Handler for start command
func (bot *Bot) handleCommandStart(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /start")

	if bot.skyMgrMonitor.IsRunning() {
		log.Debug(msgMonitorAlreadyStarted)
		return bot.Send(ctx, "whisper", "markdown", msgMonitorAlreadyStarted)
	} else {
		log.Debug(msgMonitorStart)
		cancelContext, cancelFunc := context.WithCancel(context.Background())
		bot.skyMgrMonitor.CancelFunc = cancelFunc
		go bot.skyMgrMonitor.Run(cancelContext, bot.config.Monitor.IntervalSec)
		return bot.Send(ctx, "whisper", "markdown", msgMonitorStart)
	}
}

// Handler for stop command
func (bot *Bot) handleCommandStop(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /stop")

	if bot.skyMgrMonitor.IsRunning() {
		log.Debug(msgMonitorStop)
		bot.skyMgrMonitor.CancelFunc()
		bot.skyMgrMonitor.CancelFunc = nil
		log.Debug(msgMonitorStopped)
		return bot.Send(ctx, "whisper", "markdown", msgMonitorStop)
	} else {
		log.Debug(msgMonitorNotRunning)
		return bot.Send(ctx, "whisper", "markdown", msgMonitorNotRunning)
	}
}

// Handler for status command
func (bot *Bot) handleCommandStatus(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /status")
	return bot.Send(ctx, "whisper", "markdown", msgStatus)
}

// Handler for heartbeat command
func (bot *Bot) handleCommandHeartBeat(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /heartbeat")
	go botHeartbeatLoop(bot, ctx)
	return nil
}

func (bot *Bot) handleDirectMessageFallback(ctx *BotContext, text string) (bool, error) {
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
func botHeartbeatLoop(bot *Bot, ctx *BotContext) {
	ticker := time.NewTicker(time.Minute * bot.config.Monitor.HeartbeatIntMin)

	for {
		select {
		case <-ticker.C:
			bot.Send(ctx, "whisper", "markdown", msgStatus)
		}
	}
}
