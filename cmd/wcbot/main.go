// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/signal"

	"github.com/BigOokie/skywire-wing-commander/internal/telegrambot"
	"github.com/BigOokie/skywire-wing-commander/internal/utils"
	"github.com/BigOokie/skywire-wing-commander/internal/wcconst"
	log "github.com/sirupsen/logrus"
)

var wc wcBotApp

func main() {
	// Setup and initalise application logging
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
	defer appInstance.TryUnlock()

	// Setup OS Notification for Interupt or Kill signal - to cleanly terminate the app
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, os.Kill)

	log.Infoln("Skywire Wing Commander Telegram Bot - Starting.")
	defer log.Infoln("Skywire Wing Commander Telegram Bot - Stopped.")

	// Initiate a new Bot instance
	log.Infoln("Initiating Bot instance.")
	bot, err := telegrambot.NewBot(wc.config)
	if err != nil {
		log.Error(err)
		return
	}

	log.Infoln("Starting Bot instance.")
	go bot.Start()

	// Wait for the app to be signaled to terminate
	select {
	case signal := <-osSignal:
		if signal == os.Interrupt {
			log.Debugln(wcconst.MsgOSInteruptSig)
		} else if signal == os.Kill {
			log.Debugln(wcconst.MsgOSKillSig)
		}
	}
}
