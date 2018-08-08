package main

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// Handler for help command
func (bot *Bot) handleCommandHelp(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /help")
	return bot.Send(ctx, "whisper", "markdown", fmt.Sprintf(msgHelp, bot.config.Telegram.Admin))
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
	}

	log.Debug(msgMonitorStart)
	cancelContext, cancelFunc := context.WithCancel(context.Background())
	bot.skyMgrMonitor.CancelFunc = cancelFunc
	bot.skyMgrMonitor.monitorStatusMsgChan = make(chan string)

	// Start the Event Monitor - provide cancelContext
	go bot.monitorEventLoop(cancelContext, ctx, bot.skyMgrMonitor.monitorStatusMsgChan)
	// Start the monitor - provide cancelContext
	go bot.skyMgrMonitor.Run(cancelContext, bot.skyMgrMonitor.monitorStatusMsgChan, bot.config.Monitor.IntervalSec)

	return bot.Send(ctx, "whisper", "markdown", msgMonitorStart)

}

// Handler for stop command
func (bot *Bot) handleCommandStop(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /stop")

	if bot.skyMgrMonitor.IsRunning() {
		log.Debug(msgMonitorStop)
		bot.skyMgrMonitor.CancelFunc()
		bot.skyMgrMonitor.CancelFunc = nil
		close(bot.skyMgrMonitor.monitorStatusMsgChan)
		bot.skyMgrMonitor.monitorStatusMsgChan = nil
		log.Debug(msgMonitorStopped)
		return bot.Send(ctx, "whisper", "markdown", msgMonitorStop)
	}

	log.Debug(msgMonitorNotRunning)
	return bot.Send(ctx, "whisper", "markdown", msgMonitorNotRunning)

}

// Handler for status command
func (bot *Bot) handleCommandStatus(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /status")
	return bot.Send(ctx, "whisper", "markdown", fmt.Sprintf(msgStatus, bot.skyMgrMonitor.GetConnectedNodeCount()))
}

func (bot *Bot) handleDirectMessageFallback(ctx *BotContext, text string) (bool, error) {
	errmsg := fmt.Sprintf("Sorry, I only take commands. '%s' is not a command.\n\n%s", text, msgHelpShort)
	log.Debugf(errmsg)
	return true, bot.Reply(ctx, "markdown", errmsg)
}

// AddPrivateMessageHandler adds a private MessageHandler to the Bot
func (bot *Bot) AddPrivateMessageHandler(handler MessageHandler) {
	bot.privateMessageHandlers = append(bot.privateMessageHandlers, handler)
}

// AddGroupMessageHandler adds a group MessageHandler to the Bot
func (bot *Bot) AddGroupMessageHandler(handler MessageHandler) {
	bot.groupMessageHandlers = append(bot.groupMessageHandlers, handler)
}

// monitorEventLoop monitors for event messages from the SkyMgrMonitor (when running).
// Its also responsible for managing the Heartbeat (if configured)
func (bot *Bot) monitorEventLoop(runctx context.Context, botctx *BotContext, statusMsgChan <-chan string) {
	tickerHB := time.NewTicker(bot.config.Monitor.HeartbeatIntMin)

	for {
		select {
		// Monitor Status Message
		case msg := <-statusMsgChan:
			log.Debugf("Bot.monitorEventLoop: Status event: %s", msg)
			bot.Send(botctx, "whisper", "markdown", msg)

		// Heartbeat ticker event
		case <-tickerHB.C:
			log.Debug("Bot.monitorEventLoop - Heartbeat event")
			bot.Send(botctx, "whisper", "markdown", fmt.Sprintf(msgHeartbeat, bot.skyMgrMonitor.GetConnectedNodeCount()))

		// Context has been cancelled. Shutdown
		case <-runctx.Done():
			log.Debugln("Bot.monitorEventLoop - Done event.")
			return
		}
	}
}
