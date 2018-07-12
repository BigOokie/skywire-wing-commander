# Skywire Wing Commander
**Build Status**

**Master** 

[![Build Status](https://travis-ci.org/BigOokie/Skywire-Wing-Commander.svg?branch=master)](https://travis-ci.org/BigOokie/Skywire-Wing-Commander)

**Dev**

[![Build Status](https://travis-ci.org/BigOokie/Skywire-Wing-Commander.svg?branch=dev)](https://travis-ci.org/BigOokie/Skywire-Wing-Commander)

## Overview
Skywire Wing Commander is a Telegram Bot written in Go designed to help the [Skyfleet](https://skycoin.net) community monitor and manage their SkyMiners and associated Nodes.

This is currently a Work In Progress (WIP) and has been released as an ealry alpha to select group members for testing and feedback. More details will be provided as the project progresses.

Please note that this is not an official [SkyCoin](https://skycoin.net) project. If you have issues or questions - please do not bother the official team - raise them with me by logging an issue in GitHub. Also note that this is not my job - I am doing this as an active member of the Skyfleet community and will endevor to get back to you or resolve issues when I can. So please have patience and bare with me.

The intention of this project is to have a specialised Telegram Bot application (written in Go) to run on a Skycoin Skywire (Skyminer) Manager Node and provide its owner with realtime management and monitoring capabilities.

High level design intention (WIP) is covered here:
* [Wing Commander - High Level Design](https://github.com/BigOokie/skywire-telegram-notify-bot/blob/Manager-Node-Connection-Monitoring/docs/DESIGN.md)


# Known Issues
The following section outlines some known issues that need to be taken into consideration by anyone running this software:
* DNS level blocking/filtering of telegram API domain on the Skywire nodes (or routers). More info below.

* Repository Renamed. The Repository has recently been renamed. it was previously skywire-telegram-notify-bot. It is now called Skywire-Wing-Commander. Please make sure you update any references to the new repo name.

I have built and tested this on a DIY Miner (RasPi), Official SkyMiner and OSX.

# Wing Comamander Bot Setup
This section is incomplete and requires further work. Alpha and Beta testers can follow and provide feedback.

Initiate a Telegram discussion with @BotFather to create your new Bot. I wont cover specifics here - possibly this wil be covered later, but there are many articles out there on Google. 

Remeber the access token provided by @BotFather as part of this process. You will need to pass it as a command line parameter to the Bot app.

# Install and build the Wing Commander Bot
The Wing Commander Bot is designed to operate on your Manager Node - it is not expected to function on a subordinate node that does not run the Skywire Manager.

It is expected that you have Go v1.10.x installed. I will leave the installation of this to you.

To get and build the code use the following cmds:
```
mkdir -p $GOPATH/src/github.com/BigOokie
cd $GOPATH/src/github.com/BigOokie
git clone https://github.com/BigOokie/skywire-wing-commander.git

go install -v ./...
```

# Update and rebuild
To update and rebuild, use the following cmds:
```

cd $GOPATH/src/github.com/BigOokie/skywire-wing-commander
git pull origin master 
go install -v ./...
```

# Run the Bot
To run the Bot you will need the have the auth key provided by the @BotFather. If you dont have this or dont know how to have this - please dont continue (this is for Alpha testers who know what to do!)

Run the following cmd to get the Bot running on your Master Node - replacing `{BOTTOKEN_FROM_BOTFATHER}` with the token provided by the @BotFather 

To run as a background process (detached from the terminal):
```
cd $GOPATH/bin

nohup ./swwc -bottoken {BOTTOKEN_FROM_BOTFATHER} /dev/null 2>&1 & echo $! > swwc.pid
```

To run as a forground process (debug info logged to the terminal):
```
cd $GOPATH/bin

./swwc -bottoken {BOTTOKEN_FROM_BOTFATHER}
```

Once this is running, you should be able to start a private chat with your new Bot based on detailed provided by the @BotFather.

Once you are in a private chat with the bot you can make use for the commands to run it (see following section)

#Shutting down the Bot 
```
cd $GOPATH/bin

pkill -F swwc.pid
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

`/help`

Provides the user with information on the  commands supported by Wing Commander (those listed here).

`/about`

Shows information and credits about the creator of the Wing Commander and any key contributors.

`/status`

Manually request current status of the Wing Commander Bot.

`/heartbeat`

Start an automated heatbeat - the Wing Commander Bot will notify you that it is still running every 2 hours.

`/start`

Wing Commander will start monitoring the Skyminer that it is running on.
Once started, near real-time updates will be provided when specific states on the Miner or its Nodes change.

`/stop`

Wing Commander will stop monitoring the Skyminer that it is running on.


# Wing Commander Output
Wing Commander will monitor your SkyMiner and the Node that connect to it. Periodically Wing Commander will provide you with status updates about your Nodes and their connection status with the SkyMiner Manager. Additionally, Wing Commander will also monitor and notify you of changes to connection status with external (3rd party Nodes).

This can help you to manage the up-time of your own Nodes as well as monitor when external parties (Nodes) are connecting to your Miner over the the Skywire network.

# Donations most welcome
This is not my job, but I enjoy building things for the Skyfleet community. If you find my work useful, please consider donating to support it.
```
Skycoin:    ES5LccJDhBCK275APmW9tmQNEgiYwTFKQF

BitCoin:    37rPeTNjosfydkB4nNNN1XKNrrxxfbLcMA
```