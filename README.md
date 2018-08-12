# Skywire Wing Commander
<img src="assets/icons/WingCommanderLogoFull-600x600.png" width=250 height=250>

**Note:** The Skycoin Cloud included in the logo above is the property of the [Skycoin project](https://skycoin.net), and has been used here with permission of the Skycoin project.

[![Go Report Card](https://goreportcard.com/badge/github.com/BigOokie/skywire-wing-commander)](https://goreportcard.com/report/github.com/BigOokie/skywire-wing-commander)
[![Coverage Status](https://coveralls.io/repos/github/BigOokie/skywire-wing-commander/badge.svg?branch=master)](https://coveralls.io/github/BigOokie/skywire-wing-commander?branch=master)

| Build Status |  |
|--------------|--|
|**Master**|[![Build Status](https://travis-ci.org/BigOokie/skywire-wing-commander.svg?branch=master)](https://travis-ci.org/BigOokie/skywire-wing-commander)|
|**Dev**|[![Build Status](https://travis-ci.org/BigOokie/skywire-wing-commander.svg?branch=dev)](https://travis-ci.org/BigOokie/skywire-wing-commander)|

# Contents
- [Changelog](CHANGELOG.md)
- [Credits](CREDITS.md)
- [Contributors](CONTRIBUTORS.md)
- [Overview](#overview)
- [Wing Comamander Setup](#wing-comamander-setup)
    - [Create your Bot](#create-your-bot)
    - [Install and Build](#install-and-build)
    - [Update and Rebuild](#update-and-rebuild)
    - [Configuration](#configuration)
- [Running Wing Commander](#running-wing-commander)
    - [Background process](#background-process)
    - [Forground Process](#forground-process)
    - [Automatic restart](#automatic-restart)
- [Stopping Wing Commander](#stopping-wing-commander)
- [Wing Commander Commands](#wing-commander-commands)
- [Known Issues](#known-issues)
- [Donations](#Donations)

***

# Overview
**Wing Commander** is a Telegram bot written in `Go` designed to help the **[Skyfleet](https://skycoin.net)** community monitor and manage their Skyminers and associated Nodes.

This is currently a *Work In Progress (WIP)*. The initial *ALPHA* release phase has been completed, and the project is now being made available more broadly in *BETA* release form.

Please note that this **is not an official [Skycoin](https://skycoin.net) project**. If you have issues or questions - please **do not bother the Skycoin or Skywire teams** - raise any issues or feature requests in [GitHub](https://github.com/BigOokie/skywire-wing-commander/issues/new/choose). Also note that this is not my job - I am doing this as an active member of the **Skyfleet** community and will endeavor to get back to you and resolve issues when I can. Please have patience and bare with me.

The intention of this project is to have a specialised Telegram bot application (written in Go) to run on a Skycoin Skywire (Skyminer) Manager Node and provide its owner with realtime management and monitoring capabilities.

***

# Wing Commander Setup
The following section should contain all the steps requied to setup and build the bot.

## Create your Bot

<img src="assets/images/Telegram-BotFather.jpg" width=150 height=150>

Initiate a Telegram chat with the `@BotFather`. The `@BotFather` is a bot provided by Telegram and will guide you through the process of creating a new Telegram account for your bot. 

I won't cover specifics here - I suggest you Google it.

Once you have created a new bot, the `@BotFather` will provide you with details similar to the below:
```
@BotFather:
Done! Congratulations on your new bot. You will find it at t.me/{YOUR-BOT-NAME_bot}. 
You can now add a description, about section and profile picture for your bot, see /help for a list of commands. By the way, when you've finished creating your cool bot, ping our bot Support if you want a better username for it. Just make sure the bot is fully operational before you do this.

Use this token to access the HTTP API:
NNNNNNNNN:ABCD1234abc1_aBCD12345EFghiJK1234ab

For a description of the bot API, see this page: https://core.telegram.org/bots/api
```

You will need to paste the (API) `token` provided by the `@BotFather` into your `config.toml` file .

### Suggested Bot Settings
The `@BotFather` allows you to control certain settings for the bot - including its ability to participate in group chats. At this stage the bot has been designed for use in private chats only. Further, I would strongly recommend not allowing it to be used by anyone other than youself or within group context.

## Install and Build
The **Wing Commander** bot is designed to operate on your Skywire Manager Node - it is not intended or expected to function on a subordinate Node that does not run the Skywire Manager.

It is expected that you have `Go v1.10.x` installed. I will leave the installation of this to you. The below steps assume Go is correctly installed and your `$GOPATH` is correctly defined.

To get and build the code use the following commands - note paths are case-sensitive:
```sh
mkdir -p $GOPATH/src/github.com/BigOokie
cd $GOPATH/src/github.com/BigOokie
git clone https://github.com/BigOokie/skywire-wing-commander.git
cd skywire-wing-commander
go install -v ./...
```

## Update and Rebuild
To update and rebuild, use the following commands:
```sh
cd $GOPATH/src/github.com/BigOokie/skywire-wing-commander
git pull origin master 
go install -v ./...
```

## Configuration
You MUST provide a valid configuration file for the bot or it will not launch. The config file must reside in the following location `$HOME\.wingcommander\config.toml`

Refer to the following example configuration file: [config.example.toml](\src\wcbot\config.example.toml)
It is recommended to copy the example configuration file to `$HOME\.wingcommander\config.toml`. Use the example file as a template and edit the details as needed. 

The following commands can be used to setup the required folders and copy the example config file template:
```sh
cd ~
mkdir .wingcommander
cd .wingcommander
cp $GOPATH/src/github.com/BigOokie/skywire-wing-commander/src/wcbot/config.example.toml ~/.wingcommander/config.toml
cp $GOPATH/src/github.com/BigOokie/skywire-wing-commander/src/wcbot/TelegramDetails.sh ~/.wingcommander/TelegramDetails.sh
chmod +x TelegramDetails.sh 
```
Next run ```./TelegramDetails.sh``` It will ask you for your Bot API which will look similar to this ```640158980:A1HwlYeM7RWvoHflI3-55518gvETkC-hJro```

Next edit the `config.toml` file and update it with your specific settings.
```sh
nano config.toml
```

Key elements you will need in the `config.toml` file are - but feel free to edit other settings, but the defaults should be fine for most users:
- The bot token, provided by the `@BotFather`
- Your bots `ChatID`.
- Your Telegram `@` username

### Find your ChatID
To get your `ChatID` go into Telegram and send a chat message to your newly created bot (it will not respond). Once you have initiated a chat with your bot, then enter the following URL into your browser:
```
https://api.telegram.org/bot{BOTTOKEN_FROM_BOTFATHER}/getUpdates
```
The above URL should produce `JSON` output for your bot, including the `ChatID`. Paste your `ChatID` into your `config.toml` file.

An example of the `JSON` output is as follows:
```json
{"ok":true,"result":[{"update_id":111111111,
"message":{"message_id":4,"from":{"id":222222222,"is_bot":false,"first_name":"TestUser","last_name":"","username":"TestUSer","language_code":"en-US"},"chat":{"id":000000000,"first_name":"TestUser","last_name":"","username":"TestUser","type":"private"},"date":1533900000,"text":"Hello"}}]}
```
In the example above, `"chat":{"id":000000000` is what you are looking for, and specifically the `id` which in this example is `000000000`.

## Running Wing Commander

### Background process
To run **Wing Commander** as a background process (detached from the terminal):
```sh
cd $GOPATH/bin
nohup ./wcbot /dev/null 2>&1 & echo $! > wcbot.pid &
```

### Forground process
To run **Wing Commander** as a foreground process (debug info logged to the terminal):
```sh
cd $GOPATH/bin
./wcbot
```

### Automatic restart 
Use the following commands to setup an automatic startup script to check and restart the **Wing Commander** bot incase the Manager Node goes offline.
```sh 
cd $GOPATH/src/github.com/BigOokie/skywire-wing-commander/src/wcbot
chmod +x StartCommander.sh 
crontab -e 
```
Go to the bottom of the file and enter the following:

```sh
@reboot /go/src/github.com/BigOokie/skywire-wing-commander/src/wcbot/StartCommander.sh
```
Then press `CTRL+O` & `ENTER` to Save, then press `CTRL+X` to Exit.

## Stopping Wing Commander 
```sh
cd $GOPATH/bin
pkill -F wcbot.pid
```

Alternatively, if you are running **Wing Commander** interactively from the command line, you can press `CTRL+C` to shut it down gracefully.

***

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
Once started, **Wing Commander** will provide notification updates via Telegram when any Node managed by the Skyminer Manager connects or disconnects.
Additionally, the `/start` command will initiate a heartbeat which will provide a status update on a configurable cycle (interval set in `config.toml`).  The heartbeat will help you to ensure that the bot and/or the Skyminer itself is still running. **If you stop receiving the heartbeat, you need to check whats going on.**

### Heartbeat screenshot
<img src="assets/images/WingCommander-Heartbeat.png">


## Stop
`/stop`

**Wing Commander** will stop monitoring the Skyminer. This will also stop the heartbeat.

***

# Known Issues
The following section outlines some known issues that need to be taken into consideration by anyone running this software:
- DNS level blocking/filtering of Telegram API domain on the Skywire nodes (or routers). [More detail provided below](#opendns-errors).

- Repository renamed. The repository has recently been renamed. It was previously `skywire-telegram-notify-bot`. It is now called `skywire-wing-commander` (note that this is case sensitive). Please make sure you update any references to the new repo name.

I have built and tested this on the following setups: 
- DIY Raspberry Pi Miner
- Official Skyminer (using the [official prepared images](https://github.com/skycoin/skywire#ip-presetted-system-images) for the orange pi prime)
- DIY on MacOS.

**YOU TAKE FULL RESPONSIBILITY**

## OpenDNS Errors
If you are running OpenDNS or other DNS which protects and prevents access to certain domains, you may need to update the settings to ensure that the Telegram API domain (`https://api.telegram.org`) is not blocked.
In OpenDNS this is blocked by the `Chat` and `Instant message` areas.

Access to the Telegram API is being blocked if you get the following error when trying to run the Bot:
```sh
PANI[0000] Post https://api.telegram.org/bot{BOTTOKEN_FROM_BOTFATHER}/getMe: x509: certificate signed by unknown authority 
```

You can verify this by using `curl`. Run the following command and if you get the error shown below, the domain is being blocked. If the domain isn't blocked you should get a JSON response from Telegrams API.
```
curl https://api.telegram.org/bot{BOTTOKEN_FROM_BOTFATHER}/getMe

curl: (60) SSL certificate problem: unable to get local issuer certificate
```

## Date-Time set incorrectly
If the Date-time on your Node is set incorrectly you will get certificate issues as well.
If you get the following error messages (when running interactivly) you need to check the date on your node.
```sh
INFO[0000] Initiating Bot instance.
ERRO[0001] Failed to initialize Telegram API: Post https://api.telegram.org/bot{BOTTOKEN_FROM_BOTFATHER}/getMe: x509: certificate has expired or is not yet valid
INFO[0001] Skywire Wing Commander Telegram Bot - Stopped.
```
Check your system date
```sh
date
```

Change your system date with this command (replacing the example date and time with the current):
```sh
timedatectl set-time '2018-08-10 11:57'
```

# Donations
This is not my job, but I enjoy building things for the **Skyfleet** community. If you find my work useful, please consider donating to support it.
```
Skycoin:    ES5LccJDhBCK275APmW9tmQNEgiYwTFKQF
```
