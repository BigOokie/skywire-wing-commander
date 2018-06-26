package main

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
		"Skycoin: ES5LccJDhBCK275APmW9tmQNEgiYwTFKQF\n" +
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
