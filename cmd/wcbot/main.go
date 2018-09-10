// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/BigOokie/skywire-wing-commander/internal/telegrambot"
	"github.com/BigOokie/skywire-wing-commander/internal/utils"
	"github.com/BigOokie/skywire-wing-commander/internal/wcconst"
	log "github.com/sirupsen/logrus"
)

var wc wcBotApp

func main() {
	// Setup and initialise application logging
	wc.initLogging()

	// Parse and handle known command line flags
	wc.cmdFlags.parseCmdLineFlags()
	wc.cmdFlags.handleCmdLineFlags()

	// Load configuration
	wc.loadConfig()
	wc.config.PrintConfig()
	if wc.cmdFlags.dumpconfig {
		os.Exit(0)
	}

	// Check and setup application instance control. Only allow a single instance to run
	appInstance := utils.InitAppInstance(wcconst.AppInstanceID)
	defer utils.ReleaseAppInstance(appInstance)

	// Setup OS Notification for Interrupt or Kill signal - to cleanly terminate the app
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	log.Infoln("Skywire Wing Commander Telegram Bot - Starting.")
	defer log.Infoln("Skywire Wing Commander Telegram Bot - Stopped.")

	// Initiate a new Bot instance
	log.Infoln("Initiating Bot instance.")
	bot, err := telegrambot.NewBot(wc.config)
	if err != nil {
		log.Error(err)
		return
	}

	var startmsg string
	// Check to see if we are starting because of an upgrade.
	if wc.cmdFlags.upgradecompleted {
		startmsg = fmt.Sprintf("*Successfully restarted after upgrade to %s*", wcconst.BotVersion)
	} else {
		startmsg = fmt.Sprintf("*Started: %s*", wcconst.BotAppVersion)
	}
	log.Debug(startmsg)
	err = bot.SendNewMessage("markdown", startmsg)
	if err != nil {
		log.Fatalf("Failed to send startup message to Telegram: %v", err)
	}

	err = bot.SendMainMenuMessage(nil)
	if err != nil {
		log.Fatalf("Failed to Send Main Menu: %v", err)
	}

	log.Infoln("Starting Bot instance.")
	go bot.Start()

	// Wait for the app to be signaled to terminate
	signal := <-osSignal
	bot.Stop()
	log.Debugln(wcconst.MsgOSInteruptSig, signal)
}
