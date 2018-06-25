// Telegram notification bot
package main

import (
	"flag"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

// The BotConfig struct is used to store run-time configuration
// information for the bot application.
type botConfig struct {
	BotToken       string `json:"bot_token"`
	ChatID         int64  `json:"chat_id"`
	BotDebug       bool   `json:"botdebug"`
	ClientFile     string `json:"clientfile"`
	MonitorRunning bool   `json:"monitorrunning"`
}

var (
	bot    *tgbotapi.BotAPI
	config botConfig
)

// sendBotHelpMessage sends the message responce for the /help cmd
func sendBotHelpMessage(m *tgbotapi.Message) {
	sendBotMsg(m, msgHelp, false)
}

// sendBotAboutMessage sends the message responce for the /about cmd
func sendBotAboutMessage(m *tgbotapi.Message) {
	sendBotMsg(m, msgAbout, false)
}

// sendBotStatusMessage sends the message responce for the /status cmd
func sendBotStatusMessage(m *tgbotapi.Message) {
	sendBotMsg(m, msgStatus, false)
}

// startMonitor sends the message responce for the /start cmd
func startMonitor(m *tgbotapi.Message, monitorStopEvent <-chan bool) {
	if config.MonitorRunning {
		sendBotMsg(m, msgMonitorAlreadyStarted, false)
	} else {
		config.MonitorRunning = true
		sendBotMsg(m, msgMonitorStart, false)
		go watchFileLoop(selectClientFile(), monitorStopEvent)
	}
}

// stopMonitor sends the message responce for the /start cmd
func stopMonitor(m *tgbotapi.Message, monitorStopEvent chan<- bool) {
	sendBotMsg(m, msgMonitorStop, false)
	config.MonitorRunning = false
	monitorStopEvent <- true
}

// sendBotMsg sends messages from the Monitoring Bot
func sendBotMsg(m *tgbotapi.Message, msgText string, reply bool) {
	if reply {
		sendBotReplyMsgToChatID(m.Chat.ID, msgText, m.MessageID)
	} else {
		sendBotMsgToChatID(m.Chat.ID, msgText)
	}
}

// sendBotMsgToChatID sends messages from the Monitoring Bot
func sendBotMsgToChatID(chatid int64, msgText string) {
	msg := tgbotapi.NewMessage(chatid, msgText)
	msg.ParseMode = tgbotapi.ModeMarkdown
	log.Debugf("[sendBotMsgToChatID] %s", msgText)
	bot.Send(msg)
}

// sendBotReplyMsgToChatID sends messages from the Monitoring Bot
func sendBotReplyMsgToChatID(chatid int64, msgText string, replyMsgID int) {
	msg := tgbotapi.NewMessage(chatid, msgText)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyToMessageID = replyMsgID
	log.Debugf("[sendBotReplyMsgToChatID] %s", msgText)
	bot.Send(msg)
}

// handleBotMessages processes Telegram Bot commands and responds
func handleBotMessage(m *tgbotapi.Message, monitorStopEvent chan bool) {
	if !m.IsCommand() {
		log.Debugf("[handleBotMessage] Message is not a command: %s", m.Text)
		sendBotHelpMessage(m)
		return
	}

	botcmd := m.Command()
	switch botcmd {
	case "start":
		log.Debugln("[handleBotMessage] Handling /start command")
		startMonitor(m, monitorStopEvent)
		break

	case "stop":
		log.Debugln("[handleBotMessage] Handling /stop command")
		stopMonitor(m, monitorStopEvent)
		break

	case "help":
		log.Debugln("[handleBotMessage] Handling /help command")
		sendBotHelpMessage(m)
		break

	case "about":
		log.Debugln("[handleBotMessage] Handling /about command")
		sendBotAboutMessage(m)
		break

	case "status":
		log.Debugln("[handleBotMessage] Handling /status command")
		sendBotStatusMessage(m)
		break

	default:
		log.Debugln("[handleBotMessage] Unhandled command recieved.")
	}
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
func watchFileLoop(filename string, monitorStopEvent <-chan bool) {
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
				sendBotMsgToChatID(config.ChatID, msgText)
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

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	log.Infoln("Starting Skywire Telegram Notification Bot App.")
	defer log.Infoln("Stopping Skywire Telegram Notification Bot App. Bye.")
	parseFlags()

	config.ClientFile = selectClientFile()

	var err error
	bot, err = tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"token": config.BotToken,
		}).Fatal("Could not connect to telegram")
	}

	bot.Debug = config.BotDebug
	log.Infof("Telegram Bot connected and authorised on account %s", bot.Self.UserName)

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
}
