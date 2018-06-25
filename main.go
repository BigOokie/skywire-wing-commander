// Telegram notification bot
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

// Define constants used by the application
const (
	version = "v0.0.3-alpha"

	// Bot command messages:
	// Help message
	msgHelp = "I will notify you of connections made to or from your SkyMiner nodes.\n\n" +
		"*Usage:*\n" +
		"- /about - show information and credits about my creator and any contributors\n" +
		"- /help - show this message\n" +
		"- /status - ask me how I'm going.. and if I'm still running\n" +
		"- /start - start me monitoring your Skyminer. Once started, I will start sending notifications\n" +
		"- /stop - stop me monitoring your Skyminer. Once stopped, I won't send any more notifications\n" +
		"\n" +
		"\n" +
		"Note that the Bot is bound to the _conversation_. This means it is between you and me (the bot)"

	// About cmd message
	msgAbout = "Skywire Manager Telegram Monitoring Bot (" + version + ")\n" +
		"\n" +
		"Created by @BigOokie 2018\n" +
		"GitHub: https://github.com/BigOokie/skywire-telegram-notify-bot\n" +
		"Twitter: https://twitter.com/BigOokie\n" +
		"\n" +
		"Donations most welcome 👍\n" +
		"Skycoin: 2aAprdFyxV3bqYB5yix2WsjsH1wqLKaoLhq\n" +
		"BitCoin: 37rPeTNjosfydkB4nNNN1XKNrrxxfbLcMA\n"

	// Status cmd message
	msgStatus = "I'm fine. Sill running 👍"

	// Start cmd messages
	msgMonitorAlreadyStarted = "The monitor has already been started"
	msgMonitorStart          = "Monitor starting"

	// Stop cmd message
	msgMonitorStop = "Monitor stopping"

	// Default cmd message (unhandled)
	msgDefault = "Sorry. I don't know that command."
)

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
	sendBotMsg(m, msgMonitorStart, false)
	go watchFileLoop(selectClientFile(), monitorStopEvent)
}

// stopMonitor sends the message responce for the /start cmd
func stopMonitor(m *tgbotapi.Message, monitorStopEvent chan<- bool) {
	sendBotMsg(m, msgMonitorStop, false)
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
		log.Debugln("[handleBotMessage] Hanling /start command")
		startMonitor(m, monitorStopEvent)
		break

	case "stop":
		log.Debugln("[handleBotMessage] Hanling /stop command")
		stopMonitor(m, monitorStopEvent)
		break

	case "help":
		log.Debugln("[handleBotMessage] Hanling /help command")
		sendBotHelpMessage(m)
		break

	case "about":
		log.Debugln("[handleBotMessage] Hanling /about command")
		sendBotAboutMessage(m)
		break

	case "status":
		log.Debugln("[handleBotMessage] Hanling /status command")
		sendBotStatusMessage(m)
		break

	default:
		log.Debugln("[handleBotMessage] Unhandled command recieved.")
	}
}

// UserHome returns the current user home path
func userHome() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}

	return os.Getenv("HOME")
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Warnf("File does not exist: %s", filename)
		return false
	} else {
		log.Debugf("File exists: %s", filename)
		return true
	}
}

func selectClientFile() string {
	// Default to the Users home folder - but lets check
	clientfile := filepath.Join(userHome(), ".skywire", "manager", "clients.json")
	log.Debugf("[selectClientFile] Checking if file exists: %s", clientfile)

	// If we cant find the file in the users home folder - check the $GOPATH
	if fileExists(clientfile) == false {
		gopath := os.Getenv("GOPATH")
		clientfile = filepath.Join(gopath, "bin", ".skywire", "manager", "clients.json")
		log.Debugf("[selectClientFile] Checking if file exists: %s", clientfile)
		// Can't find what we are looking for
		if fileExists(clientfile) == false {
			log.Panicln("Unable to detect location of client file.")
		}
	}

	log.Debugf("[selectClientFile] Selected file: %s", clientfile)
	return clientfile
}

// clientConnection is a structure that represents the JSON file structure
// of the clients.json file (from Skywire project)
type clientConnection struct {
	Label   string `json:"label"`
	NodeKey string `json:"nodeKey"`
	AppKey  string `json:"appKey"`
	Count   int    `json:"count"`
}

// Defines an in-memory slice (dynamic array) based on the ClientConnection struct
type clientConnectionSlice []clientConnection

// Determines if the specified ClientConnection exists within the clientConnectionSlice
func (c clientConnectionSlice) Exist(rf clientConnection) bool {
	for _, v := range c {
		if v.AppKey == rf.AppKey && v.NodeKey == rf.NodeKey {
			return true
		}
	}
	return false
}

// Reads the physical Skywire Clients.JSON file into an in-memory structure
func readClientConnectionConfig() (cfs map[string]clientConnectionSlice, err error) {
	fb, err := ioutil.ReadFile(config.ClientFile)
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

// getClientConnectionListString will iterate over the ClientConnectionConfig JSON
// file and return a formatted string for all Clients and their Nodes
func getClientConnectionListString() string {
	var clientsb strings.Builder
	var condir string

	// Read the Client Connection Config (JSON) into ccc
	ccc, err := readClientConnectionConfig()
	if err == nil {
		// Iterate ccc reading the Keys (k)
		for k := range ccc {
			// Output to our string builder the current Client Type (from K)
			// Add an newline if this isnt the first itteration
			if clientsb.String() != "" {
				clientsb.WriteString("\n")
			}

			switch k {
			case "socket":
				condir = "Outbound"
			case "socksc":
				condir = "Inbound"
			default:
				condir = "??"
			}

			clientsb.WriteString(fmt.Sprintf("ClientType: [%s](%s)\n", k, condir))
			// Iterate all Nodes in the current client type (ccc[k])
			for _, b := range ccc[k] {
				// Output the details for each node of this client type
				clientsb.WriteString(fmt.Sprintf("Label:   %s\n", b.Label))
				clientsb.WriteString(fmt.Sprintf("NodeKey: %s\n", b.NodeKey))
				clientsb.WriteString(fmt.Sprintf("AppKey:  %s\n", b.AppKey))
				clientsb.WriteString("\n")
			}
		}
	}
	// Return the built string
	return clientsb.String()
}

// getClientConnectionCountString will iterate over the ClientConnectionConfig JSON
// file and return a formatted string containing the count of Clients connected
// Both Outbound and Inbound
func getClientConnectionCountString() string {
	var clientsb strings.Builder
	var condir string

	// Read the Client Connection Config (JSON) into ccc
	ccc, err := readClientConnectionConfig()
	if err == nil {
		// Iterate ccc reading the Keys (k)
		for k := range ccc {
			// Output to our string builder the current Client Type (from K)
			// Add an newline if this isnt the first itteration
			if clientsb.String() != "" {
				clientsb.WriteString("\n")
			}

			switch k {
			case "socket":
				condir = "Outbound"
			case "socksc":
				condir = "Inbound"
			default:
				condir = "??"
			}

			clientsb.WriteString(fmt.Sprintf("ClientType: [%s](%s)  Count:%v\n", k, condir, len(ccc[k])))
		}
	}
	// Return the built string
	return clientsb.String()
}

// The BotConfig struct is used to store run-time configuration
// information for the bot application.
type botConfig struct {
	BotToken   string `json:"bot_token"`
	ChatID     int64  `json:"chat_id"`
	Locked     bool   `json:"locked"`
	BotDebug   bool   `json:"botdebug"`
	ClientFile string `json:"clientfile"`
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
	log.Debugf("[MFL] Seting up monitoring on file: %s", filename)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panic(err)
	}
	defer watcher.Close()
	defer log.Infof("[MFL] Stop watching file: %s", filename)

	err = watcher.Add(filename)
	if err != nil {
		log.Panic(err)
	}

	log.Infof("[MFL] Start watching file: %s", filename)
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Debugf("[MFL] (%s) handling event  [%s]\n", event.Name, event.Op)
				msgText := getClientConnectionCountString()
				sendBotMsgToChatID(config.ChatID, msgText)
			} else {
				log.Debugf("[MFL] (%s) ignorning event [%s]\n", event.Name, event.Op)
			}

		case err := <-watcher.Errors:
			log.Errorln("[MFL] Error:", err)

		case stop := <-monitorStopEvent:
			if stop {
				log.Debugln("[MFL] Stop event recieved.")
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
	log.Infoln("[Starting Skywire Telegram Notification Bot App.")
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
