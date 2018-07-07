# skywire-telegram-notify-bot

[![Build Status](https://travis-ci.org/BigOokie/skywire-telegram-notify-bot.svg?branch=master)](https://travis-ci.org/BigOokie/skywire-telegram-notify-bot)


This is currently a work in progress and has been released as an ealry alpha to select group members for testing and feedback. 

More details will be provided as the project progresses.

This document is currently being used to focus my design. Most of the below does not exist or work yet. Bare with me... lol

# Overview
The intention of this project is to have a specialised Telegram Bot application (written in Go) to run on a Skycoin Skywire (Skyminer) Manager Node and provide its owner with (near) realtime status updates based on certain events that occur within the Skyminer. Initially the focus is on incomming connections made to any of the Nodes managed by the Skyminer Manager.

High level design intention (WIP) is covered here:
* [High Level Design](https://github.com/BigOokie/skywire-telegram-notify-bot/blob/Manager-Node-Connection-Monitoring/docs/DESIGN.md)


# Known Issues
* Works for DIY builds only at the moment. Support for official is coming very soon.
* DNS level blocking/filtering of telegram API domain on the Skywire nodes (or routers). More info below.
* Raw dump of connection data when connection state changes. It is a "notification" but it could be more elegant.

I have built and tested this on both RasPi DIY and OSX. It does build on Official Miner also, but there is a coding error that will prevent it from running.

# Bot Setup
This section is incomplete and requires further work. Alpha and Beta testers can follow and provide feedback.

Initiate a Telegram discussion with @BotFather to create your new Bot. I wont cover specifics here - possibly this wil be covered later, but there are many articles out there on Google. 

Remeber the access token provided by @BotFather as part of this process. You will need to pass it as a command line parameter to the Bot app.

# Install and build the Bot
Get the code on your Skycoin Skywire Manager node. I have not tested this on other nodes and it is not expectred to work.

It is expected that you have Go v1.10.x installed. I will leave this to you.

To get and build the code use the following cmds:
```
mkdir -p $GOPATH/src/github.com/BigOokie
cd $GOPATH/src/github.com/BigOokie
git clone https://github.com/BigOokie/skywire-telegram-notify-bot.git

go install -v ./...
```

# Update and rebuild the bot
To update and rebuild the Bot, use the following cmds:
```

cd $GOPATH/src/github.com/BigOokie/skywire-telegram-notify-bot
git pull origin master 
go install -v ./...
```

# Run the Bot
To run the Bot you will need the have the auth key provided by the @BotFather. If you dont have this or dont know how to have this - please dont continue (this is for Alpha testers who know what to do!)

Run the following cmd to get the Bot running on your Master Node - replacing `{BOTTOKEN_FROM_BOTFATHER}` with the token provided by the @BotFather 

```
cd $GOPATH/bin

nohup ./skywire-telegram-notify-bot -bottoken {BOTTOKEN_FROM_BOTFATHER} /dev/null 2>&1 & echo $! > skywire-bot.pid
```

Once this is running, you should be able to start a private chat with your new Bot based on detailed provided by the @BotFather.

Once you are in a private chat with the bot you can make use for the commands to run it (see following section)

#Shutting down the Bot 
```
cd $GOPATH/bin

pkill -F skywire-bot.pid
```

## OpenDNS Errors (Domain Blocked)
If you are running OpenDNS or other DNS which protects and prevents access to certain domains, you may need to update the settings to ensure that the Telegram API domain (https://api.telegram.org) is not blocked.
In OpenDNS this is blocked by the `Chat` and `Instant message` areas.

Access to the Telegram API is being blocked if you get the following error when trying to run the Bot
```
PANI[0000] Post https://api.telegram.org/bot{BOTTOKEN_FROM_BOTFATHER}/getMe: x509: certificate signed by unknown authority 
```

You can verify this i using `curl`. Run the following cmd and if you get the error show, the Domain is being blocked. If the domain isnt blocked you should get a JSON response from Telegrams API.
```
curl https://api.telegram.org/bot{BOTTOKEN_FROM_BOTFATHER}/getMe

curl: (60) SSL certificate problem: unable to get local issuer certificate
```

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


# Bot Output
When the Bot has been started, if a Node connects to you (or if you connect to a Node) - you should be notified with something similar to the following (Note that Node and App Keys are not real and have been made up below):

```
ClientType: [socket](Outbound)  Count:2
ClientType: [socksc](Inbound)  Count:1
```

What you are interested in is the bottom section where the `ClientType [socksc]`. This tells you WHO has connected to YOU (inbound).

The top section `ClientType [socket]` tells you WHO YOU have connected to (outbound).

# Donation-ware
If you found my tips useful, consider providing a tip of your own ;-)
```
Skycoin:    ES5LccJDhBCK275APmW9tmQNEgiYwTFKQF

BitCoin:    37rPeTNjosfydkB4nNNN1XKNrrxxfbLcMA
```