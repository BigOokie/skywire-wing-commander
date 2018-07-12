package wingcommander

// Define constants used by the application
const (
	version = "v0.0.4-alpha"

	// Bot command messages:
	// Help message
	msgHelp = "Wing Commander here. I will help you to manage and control your SkyMiner and its Nodes.\n\n" +
		"*Usage:*\n" +
		"- /about - show information and credits about my creator and any contributors\n" +
		"- /help - show this message\n" +
		"- /status - ask me how I'm going.. and if I'm still running\n" +
		"- /heartbeat - ask me to send you a notification every 2hrs to let you know I’m still running\n" +
		"- /start - start me activly monitoring your SkyMiner. Once started, I will send notifications to you for events that occur\n" +
		"- /stop - stop me monitoring your SkyMiner. Once stopped, I won't send any more notifications\n" +
		"\n" +
		"\n" +
		"Note: I am bound to _this_ conversation."

	// About cmd message
	msgAbout = "Skywire Wing Commander: A Telegram Bot written in Go designed to help the Skyfleet community monitor and manage their SkyMiners and associated Nodes. (" + version + ")\n" +
		"\n" +
		"Created by @BigOokie 2018\n" +
		"GitHub: https://github.com/BigOokie/skywire-wing-commander\n" +
		"Twitter: https://twitter.com/BigOokie\n" +
		"\n" +
		"Donations most welcome 👍\n" +
		"Skycoin: ES5LccJDhBCK275APmW9tmQNEgiYwTFKQF\n" +
		"BitCoin: 37rPeTNjosfydkB4nNNN1XKNrrxxfbLcMA\n"

	// Status cmd message
	msgStatus = "Wing Commander reporting and ready for duty 👍"

	// Start cmd messages
	msgMonitorAlreadyStarted = "Wing Commander Active Monitoring has already been started."
	msgMonitorStart          = "Wing Commander Acrove Monitoring starting..."

	// Stop cmd message
	msgMonitorStop = "Wing Commander Active Monitoring stopping..."

	// Default cmd message (unhandled)
	msgDefault = "Sorry. I don't know that command."
)
