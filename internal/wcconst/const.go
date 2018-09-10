// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package wcconst

// Define constants used by the application
const (
	BotVersion    = "v0.2.0-beta.11"
	BotAppVersion = "Wing Commander " + BotVersion

	AnalyticsID = "UA-125471033-1"

	AppInstanceID = "wing-commander-84F95320-8C2D-4236-9252-A322F01B91A7"

	MsgAppInstErr = "Another instance of Wing Commander has been detected running on this system.\n\n" +
		"To identify and terminate (kill) ALL instances of Wing Commander on this system, run:\n\n" +
		"   pgrep wcbot | xargs kill\n\n" +
		"Exiting\n"

	// ScriptPath defines the path to the Scripts folder
	ScriptPath = "/src/github.com/BigOokie/skywire-wing-commander/scripts/"

	MsgCmdLineHelp = "Wing Commander Help\n" +
		"Command line flags:\n" +
		"  -v       display application version information.\n" +
		"  -config  display application configuration information.\n" +
		"  -help    display this message.\n" +
		"  -about   display information about the application and its author.\n\n\n" +
		MsgHelpShort

	// Bot command messages:
	// Help message
	MsgHelpShort = "*Telegram Usage:*\n" +
		"- /help - show this message.\n" +
		"- /about - show information and credits about my creator and any contributors.\n" +
		"- /status - request a status update. This provides the same information as the Heartbeat.\n" +
		"- /showconfig - display runtime configuration (from config.toml).\n" +
		"- /start - start activly monitoring your Skyminer. Once started, notifications will be sent to you for events that occur. A heartbeat will also be initiated to let you know if the bot and the Miner are still running.\n" +
		"- /stop - stop monitoring your Skyminer. Once stopped, I won't send any more notifications.\n" +
		"- /checkupdate - check GitHub for new updates.\n" +
		"- /update - attempt to update *Wing Commander* to the latest version from GitHub source.\n" +
		"- /uptime - dynamically generate a link to the Skywirenc.com site to check uptime for locally connected Nodes."

	MsgHelp = "*Wing Commander* here. I will help you to manage and monitor your Skyminer and its Nodes.\n\n" +
		MsgHelpShort +
		"\n" +
		"\n" +
		"Note: I am bound to this chat. I will only respond to commands from my configured Admin (%s)."

	// About cmd message
	MsgAbout = "*Wing Commander (" + BotVersion + ")*\n" +
		"A Telegram bot written in *Go* designed to help the *Skyfleet* community monitor and manage their *Skyminers*.\n" +
		"\n" +
		"*Created by:* @BigOokie *2018*\n" +
		"*GitHub:* https://github.com/BigOokie/skywire-wing-commander\n" +
		"*Twitter:* https://twitter.com/BigOokie\n" +
		"\n" +
		"Issues and feature requests must be logged via [GitHub](https://github.com/BigOokie/skywire-wing-commander/issues/new/choose)\n" +
		"\n" +
		"*Skycoin*: https://www.skycoin.net/\n" +
		"\n" +
		"*Donations most welcome* 👍\n" +
		"*Skycoin:* ES5LccJDhBCK275APmW9tmQNEgiYwTFKQF"

	MsgShowConfig = "Wing Commander Configuration\n" +
		"```\n%s\n```\n"

	MsgErrorGetNodes     = "⚠️ Problem getting connected Nodes from the Manager."
	MsgErrorGetDiscNodes = "Problem checking the Discovery Server."
	MsgDiscSomeNodes     = "Some Nodes are not connected to the Discovery Server."

	MsgConnectedNodes = "*Connected Nodes:* %v"
	MsgDiscConnNodes  = "*Discovery Connected Nodes:* %v"

	// Status cmd message
	MsgStatus = "%v*Wing Commander Status*\n" + MsgConnectedNodes + "\n" + MsgDiscConnNodes + "\n%s"
	// Heartbeat message
	MsgHeartbeat = "%v*Wing Commander Heatbeat* ❣️\n" + MsgConnectedNodes + "\n" + MsgDiscConnNodes + "\n%s"

	// Node Connect/Disconnect Event Messages
	MsgNodeConnected    = "*Node Connected:* %s\n\n" + MsgConnectedNodes
	MsgNodeDisconnected = "‼ *Node Disconnected:* %s\n\n" + MsgConnectedNodes

	// Start cmd messages
	MsgMonitorAlreadyStarted = "️️*Wing Commander* Monitoring has already been started."
	MsgMonitorStart          = "*Wing Commander* Monitoring starting..."

	// Stop cmd message
	MsgMonitorStop       = "*Wing Commander* Monitoring stopping..."
	MsgMonitorStopped    = "*Wing Commander* Monitoring stopped..."
	MsgMonitorNotRunning = "*Wing Commander* Monitoring is not running..."

	// Default cmd message (unhandled)
	//msgDefault = "Sorry. I don't know that command."

	// OS Interrupt Signals
	MsgOSInteruptSig = "*Wing Commander* OS Interupt Signal Received. Exiting."
	//MsgOSKillSig     = "*Wing Commander* OS Kill Signal Received. Exiting."
)
