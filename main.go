// Telegram notification bot
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

//var clientPath = filepath.Join(file.UserHome(), ".skywire", "manager", "clients.json")
var clientPath = "./test.json"

type ClientConnection struct {
	Label   string `json:"label"`
	NodeKey string `json:"nodeKey"`
	AppKey  string `json:"appKey"`
	Count   int    `json:"count"`
}
type clientConnectionSlice []ClientConnection

// Determines if the specified ClientConnection exists within the clientConnectionSlice
func (c clientConnectionSlice) Exist(rf ClientConnection) bool {
	for _, v := range c {
		if v.AppKey == rf.AppKey && v.NodeKey == rf.NodeKey {
			return true
		}
	}
	return false
}

func readClientConnectionConfig() (cfs map[string]clientConnectionSlice, err error) {
	fb, err := ioutil.ReadFile(clientPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfs = nil
			err = nil
			return
		} else {
			return
		}
	}
	cfs = make(map[string]clientConnectionSlice)
	err = json.Unmarshal(fb, &cfs)
	if err != nil {
		return
	}
	return
}

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

	err = watcher.Add(filename)
	if err != nil {
		log.Panic(err)
	}

	log.Infof("Now watching file: %s", filename)
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				msgText := fmt.Sprintf("File watcher(%s) handling event  [%s]", event.Name, event.Op)
				log.Info(msgText)
				eventMsg <- msgText
			} else {
				log.Debugf("File watcher(%s) ignorning event [%s]", event.Name, event.Op)
			}
		case err := <-watcher.Errors:
			log.Errorln("File watcher error:", err)
		}
		time.Sleep(2 * time.Second)
	}
}

// Create new telegram bot using the bot token passed on the cmd line
func startTelegramBot(eventMsg <-chan string) {
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
		break
	}

	for {
		txt := fmt.Sprintf("Event Received: %v", <-eventMsg)
		log.Debugln(txt)
		msg := tgbotapi.NewMessage(config.ChatID, txt)
		telegramBot.Send(msg)
	}
}

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	log.Infoln("Starting Skywire Telegram Notification Bot App.")
	parseFlags()

	cfs, err := readClientConnectionConfig()
	if err == nil {
		for k, _ := range cfs {
			log.Debugf("Client Type: [%s]", k)
			for _, b := range cfs[k] {
				log.Debugf("NodeKey: %s", b.NodeKey)
				log.Debugf("AppKey:  %s", b.AppKey)
				log.Debugln("")
			}
		}
	}

	return

	msgChannel := make(chan string, 1)

	// Start the telegram bot
	go startTelegramBot(msgChannel)

	// Start watching the Skywire Monitors clients.json file
	go watchFile(msgChannel, "./test.json")

	for {
		txt := <-msgChannel
		log.Debugf("Event: %v", txt)
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
