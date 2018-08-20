#!/bin/sh
# Auto Startup Wing Commander 

export GOPATH=$HOME/go

cd ${GOPATH}/bin
pkill -F wcbot.pid
rm wcbot.pid