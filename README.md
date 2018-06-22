# skywire-telegram-notify-bot
This is currently a work in progress and has been released as an ealry alpha to select group members for testing and feedback. 

More details will be provided as the project progresses.

This document is currently being used to focus my design. Most of the below does not exist or work yet. Bare with me... lol

# Overview
The intention of this project is to have a specialised Telegram Bot application (written in Go) to run on a Skycoin Skywire (Skyminer) Manager Node and provide its owner with (near) realtime status updates based on certain events that occur within the Skyminer. Initially the focus is on incomming connections made to any of the Nodes managed by the Skyminer Manager.

Future plans may incorporate other capabilites - but these are out of scope for the time being.

# Bot Setup
This section is incomplete and requires further work. Alpha and Beta testers can follow and provide feedback.

Initiate a Telegram discussion with @BotFather to create your new Bot. I wont cover specifics here - possibly this wil be covered later, but there are many articles out there on Google. 

Remeber the access token provided by @BotFather as part of this process. You will need to pass it as a command line parameter to the Bot app.

#Build the Bot
Get the code on your Skycoin Skywire Manager node. I have not tested this on other nodes and it is not expectred to work.

It is expected that you have Go v1.10.x installed. I will leave this to you.

To get and build the code use the following cmds:
```
mkdir -p $GOPATH/src/github.com/BigOokie
cd $GOPATH/src/github.com/BigOokie
git clone https://github.com/BigOokie/skywire-telegram-notify-bot.git

go install -v main.go
```

# Run the Bot
To run the Bot you will need the have the auth key provided by the @BotFather. If you dont have this or dont know how to have this - please dont continue (this is for Alpha testers who know what to do!)

Run the following cmd to get the Bot running on your Master Node

```
cd $GOPATH/bin
./skywire-telegram-notify-bot -bottoken [TOKEN_FROM_BOTFATHER]
```

Once this is running, you should be able to start a private chat with your new Bot based on detailed provided by the @BotFather.

Once you are in a private chat with the bot you can make use for the commands to run it (see following section)

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
Request current status from the Bot


# Donation-ware
If you found my tips useful, consider providing a tip of your own ;-)
```
Skycoin:    2aAprdFyxV3bqYB5yix2WsjsH1wqLKaoLhq

BitCoin:    37rPeTNjosfydkB4nNNN1XKNrrxxfbLcMA
```