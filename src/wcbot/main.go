package main

import (
	"flag"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

/*



// watchFile will watch the file specified by filename
func watchFileLoop(m *tgbotapi.Message, filename string, monitorStopEvent <-chan bool) {
	log.Debugf("[WFL] Seting up monitoring on file: %s", filename)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panic(err)
	}
	defer watcher.Close()
	defer log.Infof("[WFL] Stop watching file: %s", filename)

	err = watcher.Add(filename)
	if err != nil {
		log.Panic(err)
	}

	log.Infof("[WFL] Start watching file: %s", filename)
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Debugf("[WFL] (%s) handling event  [%s]\n", event.Name, event.Op)
				msgText := getClientConnectionCountString()
				sendBotMsg(m, msgText, false)
				//sendBotMsgToChatID(config.ChatID, msgText)
			} else {
				log.Debugf("[WFL] (%s) ignorning event [%s]\n", event.Name, event.Op)
			}

		case err := <-watcher.Errors:
			log.Errorln("[WFL] Error:", err)

		case stop := <-monitorStopEvent:
			if stop {
				log.Debugln("[WFL] Stop event recieved.")
				return
			}

		default:
			continue
		}
	}
}

*/

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
	bot, err := NewBot(config)
	if err != nil {
		log.Panic(err)
	}
	// Start the Bot Running (in the background)
	log.Infoln("Skywire Wing Commander Telegram Bot - Ready for duty.")
	go bot.Start()
	log.Infoln("Skywire Wing Commander Telegram Bot - Signing off.")

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
