# Skywire Wing Commander
<img src="assets/icons/WingCommanderLogoFull-600x600.png" width=250 height=250>

**Note:** The Skycoin Cloud logo (above) is the property of the [Skycoin project](https://skycoin.net).

| Build Status |  |
|--------------|--|
|**Master**|[![Build Status](https://travis-ci.org/BigOokie/skywire-wing-commander.svg?branch=master)](https://travis-ci.org/BigOokie/skywire-wing-commander)|
|**Dev**|[![Build Status](https://travis-ci.org/BigOokie/skywire-wing-commander.svg?branch=dev)](https://travis-ci.org/BigOokie/skywire-wing-commander)|

# Contents
- [Changelog](CHANGELOG.md)
- [Credits](CREDITS.md)
- [Overview](#overview)
- [Wing Comamander Setup](#wing-comamander-setup)
    - [Configuration File](#configuration-file)
    - [Create your Bot](#create-your-bot)
    - [Install and Build](#install-and-build)
    - [Update and Rebuild](#update-and-rebuild)
- [Running Wing Commander](#run-wing-commander)
- [Stopping Wing Commander](#stop-wing-commander)
- [Wing Commander Bot Commands](#wing-commander-bot-commands)
- [Known Issues](#known-issues)
- [Donations](#Donations)

# Overview
**Wing Commander** is a Telegram Bot written in `Go` designed to help the **[Skyfleet](https://skycoin.net)** community monitor and manage their Skyminers and associated Nodes.

This is currently a *Work In Progress (WIP)* and has been released as an early *Alpha* to select group for testing and feedback. More details will be provided as the project progresses.

Please note that this **is not an official [Skycoin](https://skycoin.net) project**. If you have issues or questions - please **do not bother the Skycoin or Skywire teams** - raise any issues or feature requests  in [GitHub](https://github.com/BigOokie/skywire-wing-commander/issues). Also note that this is not my job - I am doing this as an active member of the Skyfleet community and will endevor to get back to you and resolve issues when I can. Please have patience and bare with me.

The intention of this project is to have a specialised Telegram Bot application (written in Go) to run on a Skycoin Skywire (Skyminer) Manager Node and provide its owner with realtime management and monitoring capabilities.

# Wing Commander Setup
This section is incomplete and requires further work. It should e sufficient however for those interested in working with the *Alpha* release to get it running.

## Configuration File
You MUST provide a valid configuration file for the Bot or it will not launch. The config file must reside in the following location `$HOME\.wingcommander\config.toml`

You will need to create the `.wingcommander` folder and then place the `config.toml` file into it.
```sh
cd ~
mkdir .wincommander
```

Refer to the provided example configution file: `BigOokie\skywire-wing-commander\src\wcbot\config.example.toml` file for details of all requied settings. I suggest copying this as a template and then using a text editor such as `nano` or `vi` to edit the details. 

## Create your Bot
<img src="assets/images/Telegram-BotFather.jpg" width=150 height=150>

Initiate a Telegram chat with the `@BotFather`. The `@BotFather` is a Bot provided by Telegram and will guide you through the process of creating a new Telegram account for your Bot. 

I won't cover specifics here - I suggest you Google it.

Once you have created a new Bot, the `@BotFather` will provide you with details similar to the below:
```
@BotFather:
Done! Congratulations on your new bot. You will find it at t.me/{YOUR-BOT-NAME_bot}. 
You can now add a description, about section and profile picture for your bot, see /help for a list of commands. By the way, when you've finished creating your cool bot, ping our Bot Support if you want a better username for it. Just make sure the bot is fully operational before you do this.

Use this token to access the HTTP API:
NNNNNNNNN:ABCD1234abc1_aBCD12345EFghiJK1234ab

For a description of the Bot API, see this page: https://core.telegram.org/bots/api
```

You will need to paste the (API) `token` provided by the `@BotFather` into your `config.toml` file .

### Suggested Bot Settings
The `@BotFather` allows you to control certain settings for the Bot - including its ability to participate in Group Chats. At this stage the Bot has been designed for use in private chats only. Further, I would strongly recommend not allowing it to be used by anyone other than youself or within Group context

## Install and Build
The **Wing Commander** Bot is designed to operate on your Skywire Manager Node - it is not intended or expected to function on a subordinate Node that does not run the Skywire Manager.

It is expected that you have `Go v1.10.x` installed. I will leave the installation of this to you. The below steps assume Go is correctly installed and your `$GOPATH` is correctly defined.

To get and build the code use the following commands - note paths are case-sensative:
```sh
mkdir -p $GOPATH/src/github.com/BigOokie
cd $GOPATH/src/github.com/BigOokie
git clone https://github.com/BigOokie/skywire-wing-commander.git
cd skywire-wing-commander
go install -v ./...
```

## Update and rebuild
To update and rebuild, use the following commands:
```sh
cd $GOPATH/src/github.com/BigOokie/skywire-wing-commander
git pull origin master 
go install -v ./...
```

## Run Wing Commander
To run the **Wing Commander** Bot you must have a `config.toml` file setup. At present, the `config.toml` file MUST be placed into the same folder that the `wcbot` application (i.e. `$GOPATH/bin/`). This will be moved into the users folder at some point in the future.

There is an example configuration file provided with the source (`config.example.toml`). Easiest way to start is to copy this and replace the required values.

Key elements you will need in the `config.toml` file are:
- The Bot token, provided by the `@BotFather`
- Your Bots `ChatID`.
- Your Telegram `@` user name

### Find your ChatID
To get your `ChatID` go into Telegram and send a chat message to your newly created Bot (it will not respond). Once you have initiated a chat with your bot, then enter the following URL into your browser:
```sh
https://api.telegram.org/bot<YourBOTToken>/getUpdates
```
The above URL should produce `JSON` output for your Bot, including the `ChatID`. Paste your `ChatID` into your `config.toml` file.

### Run as background process
To run **Wing Commander** as a background process (detached from the terminal):
```sh
cd $GOPATH/bin
nohup ./wcbot /dev/null 2>&1 & echo $! > wcbot.pid &
```

### Run as forground process
To run **Wing Commander** as a forground process (debug info logged to the terminal):
```sh
cd $GOPATH/bin
./wcbot
```

Once the **Wing Commander** Bot is running, start a private chat with the Bot you can make use for the commands to run it (see following section)

## Stop Wing Commander 
```sh
cd $GOPATH/bin
pkill -F wcbot.pid
```

Alternativly, if you are running **Wing Commander** interactivly from the command line, you can pres `CTRL+C` to shut it down gracefully.

# Wing Commander Commands
This section outlines the Telegram Bot commands that are currently supported by **Wing Commander**:
- [Help](#help)
- [About](#about)
- [Status](#status)
- [Start](#start)
- [Stop](#stop)

## Help
`/help`

Provides the user with information on the commands supported by **Wing Commander** (listed below).

## About
`/about`

Shows information and credits about the creator of **Wing Commander** and any key contributors.

## Status
`/status`

Manually request current status of the **Wing Commander** Bot.

## Start
`/start`

**Wing Commander** will start monitoring the **Skyminer** that it is running on.
Once started, **Wing Commander** will provide notification updates via Telegram when any Node managed by the Skyminer connects or disconnects.
Additionally, the `/start` command will initiate a Heartbeat which will provide a status update on a configurable cycle (interval set in `config.toml`).  The Heartbeat will help you to ensure that the Bot and/or the Skyminer itself is still running. If you stop recieving the Heartbeat - you need to check whats going on.

## Stop
`/stop`

**Wing Commander** will stop monitoring the Skyminer. This will also stop the Heartbeat.

# Known Issues
The following section outlines some known issues that need to be taken into consideration by anyone running this software:
- DNS level blocking/filtering of telegram API domain on the Skywire nodes (or routers). [More detail provided below](#opendns-errors).

- Repository Renamed. The Repository has recently been renamed. It was previously `skywire-telegram-notify-bot`. It is now called `skywire-wing-commander` (note that this is case sensative). Please make sure you update any references to the new repo name.

I have built and tested this on the following setups - but please note it is still considered *ALPHA*: 
- DIY Raspberry Pi Miner
- Official Skyminer (using official OrangePi images)
- DIY on MacOS.

**YOU TAKE FULL RESPONSIBILITY**

## OpenDNS Errors
If you are running OpenDNS or other DNS which protects and prevents access to certain domains, you may need to update the settings to ensure that the Telegram API domain (`https://api.telegram.org`) is not blocked.
In OpenDNS this is blocked by the `Chat` and `Instant message` areas.

Access to the Telegram API is being blocked if you get the following error when trying to run the Bot
```sh
PANI[0000] Post https://api.telegram.org/bot{BOTTOKEN_FROM_BOTFATHER}/getMe: x509: certificate signed by unknown authority 
```

You can verify this i using `curl`. Run the following cmd and if you get the error show, the Domain is being blocked. If the domain isnt blocked you should get a JSON response from Telegrams API.
```sh
curl https://api.telegram.org/bot{BOTTOKEN_FROM_BOTFATHER}/getMe

curl: (60) SSL certificate problem: unable to get local issuer certificate
```

# Donations
This is not my job, but I enjoy building things for the **Skyfleet** community. If you find my work useful, please consider donating to support it.
```
Skycoin:    ES5LccJDhBCK275APmW9tmQNEgiYwTFKQF

BitCoin:    37rPeTNjosfydkB4nNNN1XKNrrxxfbLcMA
```