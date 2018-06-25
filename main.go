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
	"sync"
	"time"

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
		"*Usage*\n" +
		"- /about | show information and credits about my creator and any contributors\n" +
		"- /help | show this message\n" +
		"- /status | ask me how I'm going.. and if I'm still running\n" +
		"- /chatid | tell me our current telegram chatid\n" +
		"- /start | start me monitoring your Skyminer. Once started, I will start sending notifications\n" +
		"- /stop | stop me monitoring your Skyminer. Once stopped, I won't send any more notifications\n" +
		"\n" +
		"\n" +
		"Note that the Bot is bound to the _conversation_. This means it is between you and me (the bot)"

	// About cmd message
	msgAbout = "Skywire Manager Telegram Monitoring Bot (" + version + ")\n" +
		"\n" +
		"Created by @BigOokie 2018\n" +
		"GitHub: https://github.com/BigOokie/skywire-telegram-notify-bot"

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

// UserHome returns the current user home path
func userHome() string {
	// os/user relies on cgo which is disabled when cross compiling
	// use fallbacks for various OSes instead
	// usr, err := user.Current()
	// if err == nil {
	// 	return usr.HomeDir
	// }
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

var (
	config      botConfig
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

// watchFile will watch the file specified by filename
func watchFile(monitorMsgEvent chan<- string, monitorStopEvent <-chan bool, filename string) {
	log.Debugf("[FWM] Seting up monitoring on file: %s", filename)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panic(err)
	}
	defer watcher.Close()
	defer log.Infof("[FWM] Stop watching file: %s", filename)

	err = watcher.Add(filename)
	if err != nil {
		log.Panic(err)
	}

	log.Infof("[FWM] Start watching file: %s", filename)
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Debugf("[FWM] (%s) handling event  [%s]\n", event.Name, event.Op)
				//msgText := getClientConnectionListString()
				msgText := getClientConnectionCountString()
				log.Debugln("[FWM] Sending monitor message event.")
				monitorMsgEvent <- msgText
				log.Debugln("[FWM] Sent monitor message event.")
			} else {
				log.Debugf("[FWM] (%s) ignorning event [%s]\n", event.Name, event.Op)
			}
		case err := <-watcher.Errors:
			log.Errorln("[FWM] Error:", err)
		case stop := <-monitorStopEvent:
			if stop {
				log.Debugln("[FWM] Stop event recieved.")
				return
			}
		default:
			continue
		}
	}
}

func sendMonitorMsg(monitorMsgEvent <-chan string) {
	telegramBot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Panic(err)
	}
	defer log.Info("[SENDER] Telegram Bot finished.")
	msg := tgbotapi.NewMessage(config.ChatID, "")

	for {
		select {
		case msg.Text = <-monitorMsgEvent:
			log.Debugln("[SENDER] Recieved message from File Monitor")
			log.Debugln("[SENDER] Begin Message")
			log.Debugln(msg.Text)
			log.Debugln("[SENDER] End Message")
			telegramBot.Send(msg)
		default:
			time.Sleep(time.Second)
		}
	}
}

// Create new telegram bot using the bot token passed on the cmd line
func startTelegramBot(botwg *sync.WaitGroup) {
	var watcherRunning = false

	telegramBot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Panic(err)
	}
	defer log.Info("[BOT] Telegram Bot finished.")

	monitorMsgEvent := make(chan string, 1)
	defer close(monitorMsgEvent)
	monitorStopEvent := make(chan bool)
	defer close(monitorStopEvent)

	// Signal the Bot has finished
	defer botwg.Done()

	telegramBot.Debug = config.BotDebug
	log.Infof("[BOT] Telegram Bot connected and authorised on account %s", telegramBot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := telegramBot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			log.Debug("[BOT] Ignoring empty message.")
			continue
		}

		config.ChatID = update.Message.Chat.ID
		msg := tgbotapi.NewMessage(config.ChatID, "")
		log.Debugf("[BOT] Message recieved from ChatID: %v", config.ChatID)

		if update.Message.IsCommand() {
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			log.Debugf("[BOT] Recieved Command: %s", update.Message.Command())
			switch update.Message.Command() {
			case "help":
				msg.Text = msgHelp
			case "about":
				msg.Text = msgAbout
			case "status":
				msg.Text = msgStatus
			case "chatid":
				msg.Text = fmt.Sprintf("ChatID: %v", update.Message.Chat.ID)
			case "start":
				if watcherRunning {
					msg.Text = msgMonitorAlreadyStarted
					log.Debugln(msg.Text)
				} else {
					watcherRunning = true
					msg.Text = msgMonitorStart
					go sendMonitorMsg(monitorMsgEvent)
					// Start watching the Skywire Monitors clients.json file
					go watchFile(monitorMsgEvent, monitorStopEvent, config.ClientFile)
				}
			case "stop":
				msg.Text = msgMonitorStop
				monitorStopEvent <- true
				watcherRunning = false
			default:
				msg.Text = msgDefault
			}
		} else {
			msg.Text = "Sorry. I don't chat much.."
		}

		if len(msg.Text) > 0 {
			log.Debugln("[BOT] Start Response")
			log.Debugln(msg.Text)
			log.Debugln("[BOT] End Response")
			telegramBot.Send(msg)
		}
	}
	//}
}

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	log.Infoln("[MAIN] Starting Skywire Telegram Notification Bot App.")
	defer log.Infoln("[MAIN] Stopping Skywire Telegram Notification Bot App. Bye.")
	parseFlags()

	config.ClientFile = selectClientFile()

	// Setup a waitgroup to sync and wait for the Telegram Bot to end.
	var botwg sync.WaitGroup
	botwg.Add(1)
	// Start the telegram bot
	go startTelegramBot(&botwg)
	botwg.Wait()
}
