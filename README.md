# skywire-telegram-notify-bot
This is currently a work in progress. Details will be provided as the project progresses.

This document is currently being used to focus my design. Most of the below does not exist or work yet. Bare with me... lol

# Overview
The intention of this project is to have a specialised Telegram Bot application (written in Go) to run on a Skycoin Skywire (Skyminer) Manager Node and provide its owner with (near) realtime status updates based on certain events that occur within the Skyminer.

Initially, plans are to monitor and alert (via Telegram private chat) changes to connection state within the Nodes in the Skyminer.

Future plans may incorporate other capabilites - but these are out of scope for the time being.

# Bot Setup
This section is incomplete and requires further work.

Initiate a Telegram discussion with @BotFather to create your new Bot. remeber the access token provided as part of this process. You will need to pass it as a command line parameter to the Bot app.

# Supported Bot Commands
This section is being used for design currently. These capabilities do not exist yet.

/hello
Used to ensure the Bot app is running. Like a ping command. The Bot should always respond with “Hi”.

/chatid
The Bot will respond with the ID of the current Telegram chat.
The Bot is intended to run in private chat only. For security reasons, the Bot will be designed to not accept or respond to any user who is not part of the chatid.

/start
Start monitoring the Skyminer that the Bot is running on.
Once started, near real-time updates will be provided by the Bot when specific states on the Miner or its Nodes change.

/stop
Stop monitoring the Skyminer thay the Bot is running on.

/status
Request current status from the Bot.
