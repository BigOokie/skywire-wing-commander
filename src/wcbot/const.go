package main

// Define constants used by the application
const (
	version = "v0.1.0-alpha.1"

	// Bot command messages:
	// Help message
	msgHelpShort = "*Usage:*\n" +
		"- /help - show this message\n" +
		"- /about - show information and credits about my creator and any contributors\n" +
		"- /status - ask how I'm going.. and if I'm still running\n" +
		"- /start - start activly monitoring your SkyMiner. Once started, notifications will be sent to you for events that occur. A heartbeat will also be initiated to let you know if the Bot and the Miner are still running.\n" +
		"- /stop - stop monitoring your SkyMiner. Once stopped, I won't send any more notifications\n"

	msgHelp = "*Wing Commander* here. I will help you to manage and control your SkyMiner and its Nodes.\n\n" +
		msgHelpShort +
		"\n" +
		"\n" +
		"Note: I am bound to this chat."

	// About cmd message
	msgAbout = "*Wing Commander (" + version + ")*\n" +
		"A Telegram Bot written in *Go* designed to help the *Skyfleet* community monitor and manage their *SkyMiners*.\n" +
		"\n" +
		"*Created by:* @BigOokie *2018*\n" +
		"*GitHub:* https://github.com/BigOokie/skywire-wing-commander\n" +
		"*Twitter:* https://twitter.com/BigOokie\n" +
		"\n" +
		"Issues and feature requests must be logged via [GitHub](https://github.com/BigOokie/skywire-wing-commander/issues)\n" +
		"\n" +
		"*SkyCoin*: https://skycoin.net/\n" +
		"\n" +
		"*Donations most welcome* 👍\n" +
		"*Skycoin:* ES5LccJDhBCK275APmW9tmQNEgiYwTFKQF\n" +
		"*BitCoin:* 37rPeTNjosfydkB4nNNN1XKNrrxxfbLcMA\n"

	// Status cmd message
	msgStatus = "*Wing Commander* Ready and reporting for duty 👍"
	// Heartbeat message
	msgHeartbeat = "*Wing Commander Heatbeat* ❤️"

	// Start cmd messages
	msgMonitorAlreadyStarted = "*Wing Commander* Sky Manager Monitoring has already been started."
	msgMonitorStart          = "*Wing Commander* Sky Manager Monitoring starting..."

	// Stop cmd message
	msgMonitorStop       = "*Wing Commander* Sky Manager Monitoring stopping..."
	msgMonitorStopped    = "*Wing Commander* Sky Manager Monitoring stopped..."
	msgMonitorNotRunning = "*Wing Commander* Sky Manager Monitoring is not running..."

	// Default cmd message (unhandled)
	msgDefault = "Sorry. I don't know that command."

	// OS Interupt Signals
	msgOSInteruptSig = "*Wing Commander* OS Interupt Signal Received. Exiting."
	msgOSKillSig     = "*Wing Commander* OS Kill Signal Received. Exiting."
)
