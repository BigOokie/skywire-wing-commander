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
	flag.Int64Var(&config.ChatID, "chatid", -1, "lock telegram bot to specific chatid")
	flag.BoolVar(&config.BotDebug, "botdebug", false, "Bot API debugging")
	flag.Parse()

	log.Debugf("Parameter: bottoken = %s", config.BotToken)
	log.Debugf("Parameter: chatid =  %v", config.ChatID)
	log.Debugf("Parameter: botdebug = %v", config.BotDebug)
	if config.ChatID == -1 {
		config.Locked = false
	} else {
		config.Locked = true
		log.Infof("Bot locked to ChatID: ", config.ChatID)
	}
}

// allowedToRespondToChat determines if the bot has been locked to a specific
// chat ID or no. It it has it then compares this to the chatid parameter provided
// to determine if it is allowed to respond or not.
func allowedToRespondToChat(chatid int64) bool {
	if config.Locked {
		// Chat is locked, so check the provided chatid to see if we are allowed to respond.
		return chatid == config.ChatID
	}

	// In all other cases, the chat is not locked, so we are allowed to respond to anyone.
	return true
}

// watchFile will watch the file specified by filename
func watchFile(filename string) {
	log.Debugf("Seting up file watch on file: %s", filename)
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
					log.Infof("File watcher(%s) handling event  [%s]", event.Name, event.Op)
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
}

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	log.Infoln("Starting Telegram Notification Bot App.")
	parseFlags()

	// Start the telegram bot
	startTelegramBot()

	// Start watching the Skywire Monitors clients.json file
	watchFile("test.json")

	/*
		// Setup bot interface debugging
		bot.Debug = config.BotDebug
		log.Infof("Authorised on account %s", bot.Self.UserName)

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
