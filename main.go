// Telegram notification bot
package main

import (
	"flag"
	"fmt"

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
	config BotConfig
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

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	log.Infoln("Starting Telegram Notification Bot App.")
	parseFlags()

	// Create new telegram bot using the bot token passed on the cmd line
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Panic(err)
	}

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
}
