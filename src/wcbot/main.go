// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/signal"
	"path/filepath"

	"github.com/BigOokie/skywire-wing-commander/src/utils"
	"github.com/BigOokie/skywire-wing-commander/src/wcconfig"
	"github.com/BigOokie/skywire-wing-commander/src/wcconst"

	log "github.com/sirupsen/logrus"
)

// loadConfig manages the configuration load specifics
// offloading the detail from the `main()` funct
func loadConfig() (config wcconfig.Config, err error) {
	log.Debugln("loadConfig: Start")
	// Load configuration
	config, err = wcconfig.LoadConfigParameters("config", filepath.Join(utils.UserHome(), ".wingcommander"), map[string]interface{}{
		"telegram.debug":                 false,
		"monitor.intervalsec":            10,
		"monitor.heartbeatintmin":        120,
		"monitor.discoverymonitorintmin": 120,
		"skymanager.address":             "127.0.0.1:8000",
		"skymanager.discoveryaddress":    "discovery.skycoin.net:8001",
	})
	log.Debugln("loadConfig: Complete")
	return
}

func main() {
	// Setup OS Notification for Interupt or Kill signal - to cleanly terminate the app
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, os.Kill)

	// Setup Log Formatter
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	log.Infoln("Skywire Wing Commander Telegram Bot - Starting.")
	defer log.Infoln("Skywire Wing Commander Telegram Bot - Stopped.")

	config, err := loadConfig()

	if err != nil {
		log.Error(err)
		return
	}

	// Initiate a new Bot instance
	log.Infoln("Initiating Bot instance.")
	bot, err := NewBot(config)
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
