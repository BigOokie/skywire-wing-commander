package main

import (
	"flag"

	"github.com/BigOokie/Skywire-Wing-Commander/wc"
	log "github.com/sirupsen/logrus"
)

/*

// startMonitor sends the message responce for the /start cmd
func startMonitor(m *tgbotapi.Message, monitorStopEvent <-chan bool) {
	if oldconfig.MonitorRunning {
		sendBotMsg(m, msgMonitorAlreadyStarted, false)
	} else {
		oldconfig.MonitorRunning = true
		sendBotMsg(m, msgMonitorStart, false)
		go watchFileLoop(m, selectClientFile(), monitorStopEvent)
		//go checkManagerNodesLoop(m, monitorStopEvent)
	}
}

// stopMonitor sends the message responce for the /start cmd
func stopMonitor(m *tgbotapi.Message, monitorStopEvent chan<- bool) {
	sendBotMsg(m, msgMonitorStop, false)
	oldconfig.MonitorRunning = false
	monitorStopEvent <- true
}

// parseFlags parses command line flags and populates the run-time applicaton configuration
func parseFlags() {
	flag.StringVar(&config.BotToken, "bottoken", "", "telegram bot token (provided by the @BotFather")
	flag.BoolVar(&config.BotDebug, "botdebug", false, "Bot API debugging")
	flag.Parse()

	log.Debugf("Parameter: bottoken = %s", config.BotToken)
	log.Debugf("Parameter: botdebug = %v", config.BotDebug)
}


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

func checkManagerNodesLoop(m *tgbotapi.Message, monitorStopEvent <-chan bool) {
	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-ticker.C:
			msgText := getGetAllNodes()
			sendBotMsg(m, msgText, false)
		case <-monitorStopEvent:
			return
		}
	}
}

*/

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	log.Infoln("Skywire Wing Commander Telegram Bot - Starting.")
	defer log.Infoln("Skywire Wing Commander Telegram Bot - Stopped.")

	configPath := flag.String("config", "config.toml", "Path to the Wing Commander config file")
	config, err := wingcommander.ReadConfig(*configPath)
	if err != nil {
		log.Panic(err)
	}

	bot, err := wingcommander.NewBot(config)
	if err != nil {
		log.Panic(err)
	}
	log.Infoln("Skywire Wing Commander Telegram Bot - Ready for duty.")
	bot.Start()
	log.Infoln("Skywire Wing Commander Telegram Bot - Signing off.")

	//log.Infoln(getGetAllNodes())

	//oldconfig.ClientFile = selectClientFile()

	/*
		bot, err = tgbotapi.NewBotAPI(config.Bot.Token)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"token": config.Bot.Token,
			}).Fatal("Could not connect to Telegram")
		}

		bot.Debug = config.Bot.Debug
		log.Infof("Skywire Wing Commander Telegram Bot connected and authorised on account %s", bot.Self.UserName)

		monitorStopEvent := make(chan bool)
		defer close(monitorStopEvent)

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, err := bot.GetUpdatesChan(u)

		for update := range updates {
			if update.Message == nil {
				continue
			}

			go handleBotMessage(update.Message, monitorStopEvent)
		}
	*/
}
