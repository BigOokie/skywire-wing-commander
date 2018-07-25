package main

import (
	"flag"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Setup OS Notification for Interupt or Kill signal - to cleanly terminate the app
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, os.Kill)

	// Setup Log Formatter
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	log.Infoln("Skywire Wing Commander Telegram Bot - Starting.")
	defer log.Infoln("Skywire Wing Commander Telegram Bot - Stopped.")

	// Load configuration
	configPath := flag.String("config", "config.toml", "Path to the Wing Commander config file")
	config, err := ReadConfig(*configPath)
	if err != nil {
		log.Panic(err)
	}

	// Initiate a new Bot instance
	log.Infoln("Initiating Bot instance.")
	bot, err := NewBot(config)
	if err != nil {
		log.Panic(err)
	}

	log.Infoln("Starting Bot instance.")
	go bot.Start()

	// Start the Bot Running (in the background)
	log.Infoln("Skywire Wing Commander Telegram Bot - Ready for duty.")
	defer log.Infoln("Skywire Wing Commander Telegram Bot - Signing off.")

	// Wait for the app to be signaled to terminate
	select {
	case signal := <-osSignal:
		if signal == os.Interrupt {
			log.Debugln(msgOSInteruptSig)
		} else if signal == os.Kill {
			log.Debugln(msgOSKillSig)
		}
	}
}
