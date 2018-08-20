#!/bin/sh
# Auto Startup Wing Commander 

export GOPATH=$HOME/go

cd $GOPATH/bin
[[ -f wcbot.pid ]] && pkill -F wcbot.pid && rm wcbot.pid