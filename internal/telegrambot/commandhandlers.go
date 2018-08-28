// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package telegrambot

import (
	"context"
	"fmt"
	"time"

	"github.com/BigOokie/skywire-wing-commander/internal/utils"
	"github.com/BigOokie/skywire-wing-commander/internal/wcconst"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

func logSendError(from string, err error) {
	log.Errorf("%s - Error: %v", from, err)
}

// Handler for help command
func (bot *Bot) handleCommandHelp(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /help")
	err := bot.Send(ctx, "whisper", "markdown", fmt.Sprintf(wcconst.MsgHelp, bot.config.Telegram.Admin))
	if err != nil {
		logSendError("Bot.handleCommandHelp", err)
	}
	return err
}

// Handler for about command
func (bot *Bot) handleCommandAbout(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /about")
	err := bot.Send(ctx, "whisper", "markdown", wcconst.MsgAbout)
	if err != nil {
		logSendError("Bot.handleCommandAbout", err)
	}
	return err
}

// Handler for showconfig command
func (bot *Bot) handleCommandShowConfig(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /showconfig")
	err := bot.Send(ctx, "whisper", "markdown", fmt.Sprintf(wcconst.MsgShowConfig, bot.config.String()))
	if err != nil {
		logSendError("Bot.handleCommandShowConfig (Send):", err)
		log.Debug("Bot.handleCommandShowConfig: Attempting to resend as text.")
		err = bot.Send(ctx, "whisper", "text", fmt.Sprintf(wcconst.MsgShowConfig, bot.config.String()))
		if err != nil {
			logSendError("Bot.handleCommandShowConfig (Resend as Text):", err)
		}
	}
	return err
}

// Handler for start command
func (bot *Bot) handleCommandStart(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /start")

	if bot.skyMgrMonitor.IsRunning() {
		log.Debug(wcconst.MsgMonitorAlreadyStarted)
		err := bot.Send(ctx, "whisper", "markdown", wcconst.MsgMonitorAlreadyStarted)
		if err != nil {
			logSendError("Bot.handleCommandStart", err)
		}
		return err
	}

	log.Debug(wcconst.MsgMonitorStart)
	cancelContext, cancelFunc := context.WithCancel(context.Background())
	bot.skyMgrMonitor.SetCancelFunc(cancelFunc)
	monitorStatusMsgChan := make(chan string)

	// Start the Event Monitor - provide cancelContext
	go bot.monitorEventLoop(cancelContext, ctx, monitorStatusMsgChan)
	// Start monitoring the local Manager - provide cancelContext
	go bot.skyMgrMonitor.RunManagerMonitor(cancelContext, monitorStatusMsgChan, bot.config.Monitor.IntervalSec)
	// Start monitoring the local Manager - provide cancelContext
	//go bot.skyMgrMonitor.RunDiscoveryMonitor(cancelContext, monitorStatusMsgChan, bot.config.Monitor.DiscoveryMonitorIntMin)

	err := bot.Send(ctx, "whisper", "markdown", wcconst.MsgMonitorStart)
	if err != nil {
		logSendError("Bot.handleCommandStart", err)
	}
	return err
}

// Handler for stop command
func (bot *Bot) handleCommandStop(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /stop")

	if bot.skyMgrMonitor.IsRunning() {
		log.Debug(wcconst.MsgMonitorStop)
		bot.skyMgrMonitor.StopManagerMonitor()
		log.Debug(wcconst.MsgMonitorStopped)
		err := bot.Send(ctx, "whisper", "markdown", wcconst.MsgMonitorStop)
		if err != nil {
			logSendError("Bot.handleCommandStop", err)
		}
		return err
	}

	log.Debug(wcconst.MsgMonitorNotRunning)
	err := bot.Send(ctx, "whisper", "markdown", wcconst.MsgMonitorNotRunning)
	if err != nil {
		logSendError("Bot.handleCommandStop", err)
	}
	return err
}

// Handler for status command
func (bot *Bot) handleCommandStatus(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /status")

	if !bot.skyMgrMonitor.IsRunning() {
		// Monitor not running
		err := bot.Send(ctx, "whisper", "markdown", wcconst.MsgMonitorNotRunning)

		if err != nil {
			logSendError("Bot.handleCommandStatus", err)
		}

		return err
	}

	// Build Status Check Message
	msg := bot.skyMgrMonitor.BuildConnectionStatusMsg(wcconst.MsgStatus)
	err := bot.Send(ctx, "whisper", "markdown", msg)
	if err != nil {
		logSendError("Bot.handleCommandStatus", err)
	}
	return err
}

// Handler for help CheckUpdate
func (bot *Bot) handleCommandCheckUpdate(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /checkupdate")
	err := bot.Send(ctx, "whisper", "markdown", "Checking for updates...")
	if err != nil {
		logSendError("Bot.handleCommandCheckUpdate", err)
		// Return if an error has occurred
		return err
	}

	updateAvailable, updateMsg := utils.UpdateAvailable("BigOokie", "skywire-wing-commander", wcconst.BotVersion)
	if updateAvailable {
		err = bot.Send(ctx, "whisper", "markdown", fmt.Sprintf("*Update available:* %s", updateMsg))
	} else {
		err = bot.Send(ctx, "whisper", "markdown", fmt.Sprintf("*Up to date:* %s", updateMsg))
	}

	if err != nil {
		logSendError("Bot.handleCommandCheckUpdate", err)
	}
	return err
}

// Handler for nodes command
func (bot *Bot) handleCommandListNodes(ctx *BotContext, command, args string) error {
	log.Debug("Handle command /listnodes")

	if bot.skyMgrMonitor.GetConnectedNodeCount() == 0 {
		log.Debug("Bot.handleCommandListNodes: No connected Nodes.")
		err := bot.Send(ctx, "whisper", "markdown", "No connected Nodes.")
		if err != nil {
			logSendError("Bot.handleCommandListNodes", err)
		}
		return err
	}

	var replyKeyboard tgbotapi.InlineKeyboardMarkup
	var keyboard [][]tgbotapi.InlineKeyboardButton
	var btnrow []tgbotapi.InlineKeyboardButton
	var btn tgbotapi.InlineKeyboardButton

	// Iterate the connectedNodes and build a keyboard with one button
	// containing the Node Key per row
	nodeMap := bot.skyMgrMonitor.GetConnectedNodes()
	log.Debugf("Bot.handleCommandListNodes: Node Count %v", len(nodeMap))
	for _, v := range nodeMap {
		log.Debugf("Bot.handleCommandListNodes: Creating button for Node: %s", v.Key)
		btn = tgbotapi.NewInlineKeyboardButtonData(v.Key, v.Key)
		btnrow = tgbotapi.NewInlineKeyboardRow(btn)
		keyboard = append(keyboard, btnrow)
	}

	replyKeyboard = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}

	/*

		var nodeListKB = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("1"),
				tgbotapi.NewKeyboardButton("2"),
				tgbotapi.NewKeyboardButton("3"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("4"),
				tgbotapi.NewKeyboardButton("5"),
				tgbotapi.NewKeyboardButton("6"),
			),
		)
	*/

	// Mark the keyboard as one time use. The keyboard will be hidden
	// once a button is selected
	//nodeListKB.OneTimeKeyboard = true

	err := bot.SendReplyInlineKeyboard(ctx, replyKeyboard)
	if err != nil {
		log.Error(err)
	}

	return err
}

func (bot *Bot) handleDirectMessageFallback(ctx *BotContext, text string) (bool, error) {
	errmsg := fmt.Sprintf("Sorry, I only take commands. '%s' is not a command.\n\n%s", text, wcconst.MsgHelpShort)
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
			err := bot.Send(botctx, "whisper", "markdown", msg)
			if err != nil {
				logSendError("Bot.monitorEventLoop", err)
			}

		// Heartbeat ticker event
		case <-tickerHB.C:
			log.Debug("Bot.monitorEventLoop - Heartbeat event")
			// Build Heartbeat Status Message
			msg := bot.skyMgrMonitor.BuildConnectionStatusMsg(wcconst.MsgHeartbeat)
			log.Debug(msg)
			err := bot.Send(botctx, "whisper", "markdown", msg)
			if err != nil {
				logSendError("Bot.handleCommandStatus", err)
			}

		// Context has been cancelled. Shutdown
		case <-runctx.Done():
			log.Debugln("Bot.monitorEventLoop - Done event.")
			return
		}
	}
}
