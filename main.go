// Telegram notification bot
package main

import (
	"flag"
	"fmt"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

// The BotConfig struct is used to store run-time configuration
// information for the bot application.
type BotConfig struct {
	BotToken string `json:"bot_token"`
	ChatID   int64  `json:"chat_id"`
	Locked   bool   `json:"locked"`
	BotDebug bool   `json:"botdebug"`
}

var (
	config      BotConfig
	telegramBot tgbotapi.BotAPI
)

// parseFlags parses command line flags and populates the run-time applicaton configuration
func parseFlags() {
	flag.StringVar(&config.BotToken, "bottoken", "", "telegram bot token (provided by the @BotFather")
	flag.BoolVar(&config.BotDebug, "botdebug", false, "Bot API debugging")
	flag.Parse()

	log.Debugf("Parameter: bottoken = %s", config.BotToken)
	log.Debugf("Parameter: botdebug = %v", config.BotDebug)
}

func chatSetup() bool {
	return (config.ChatID != -1)
}

// watchFile will watch the file specified by filename
func watchFile(eventMsg chan<- string, filename string) {
	log.Infof("Seting up file watch on file: %s", filename)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panic(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		log.Infof("Now watching file: %s", filename)
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					msgText := fmt.Sprintf("File watcher(%s) handling event  [%s]", event.Name, event.Op)
					log.Info(msgText)
					//eventMsg <- msgText
				} else {
					log.Debugf("File watcher(%s) ignorning event [%s]", event.Name, event.Op)
				}
			case err := <-watcher.Errors:
				log.Errorln("File watcher error:", err)
			}
		}
	}()

	err = watcher.Add(filename)
	if err != nil {
		log.Panic(err)
	}
	<-done
}

// Create new telegram bot using the bot token passed on the cmd line
func startTelegramBot() {
	telegramBot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Panic(err)
	}
	telegramBot.Debug = config.BotDebug
	log.Infof("Telegram Bot connected and authorised on account %s", telegramBot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := telegramBot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			log.Warn("Ignoring empty message.")
			continue
		}
		config.ChatID = update.Message.Chat.ID
		log.Debugf("Message recieved from ChatID: %v", config.ChatID)
		msg := tgbotapi.NewMessage(config.ChatID, "")
		msg.Text = "Hello. I'm up and running. Further updates will be provided in this chat session."
		telegramBot.Send(msg)
		return
	}
}

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	log.Infoln("Starting Skywire Telegram Notification Bot App.")
	parseFlags()

	msgChannel := make(chan string, 1)

	// Start the telegram bot
	startTelegramBot()

	// Start watching the Skywire Monitors clients.json file
	watchFile(msgChannel, "./test.json")

	for {
		log.Debugf("Event: %v", <-msgChannel)
	}

	/*
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates, err := bot.GetUpdatesChan(u)
		for update := range updates {
			if update.Message == nil {
				log.Warn("Ignoring empty message.")
				continue
			}

			log.Infof("Message recieved from ChatID: %v", update.Message.Chat.ID)

			if allowedToRespondToChat(update.Message.Chat.ID) {
				// Allowed to respond to this Chat ID
				if update.Message.IsCommand() {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					log.Debugf("Command: %s", update.Message.Command())
					switch update.Message.Command() {
					case "help":
						msg.Text = "type /sayhi or /status or /chatid or /lock or /unlock."
					case "sayhi":
						msg.Text = "Hi :)"
					case "status":
						msg.Text = "I'm ok."
					case "chatid":
						msg.Text = fmt.Sprintf("ChatID: %v", update.Message.Chat.ID)
					case "lock":
						if config.ChatID == -1 {
							config.ChatID = update.Message.Chat.ID
							msg.Text = fmt.Sprintf("Locked to current ChatID [%v]", update.Message.Chat.ID)
						} else {
							msg.Text = "Already locked."
						}
					case "unlock":
						if config.ChatID == -1 {
							msg.Text = "Already unlocked."
						} else {
							config.ChatID = -1
							msg.Text = "Unlocked."
						}
					default:
						msg.Text = "Sorry. I don't know that command."
					}
					bot.Send(msg)
				}
			} else {
				log.Warnf("Chat is locked to ID [%v]. Not allowed to respond to Chat ID [%v]", config.ChatID, update.Message.Chat.ID)
			}
		}
	*/
}
