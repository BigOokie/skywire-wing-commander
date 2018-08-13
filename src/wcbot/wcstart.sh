#!/bin/sh
# Auto Startup Wing Commander 

export GOPATH="$HOME/go"

cd $GOPATH/bin
./wcbot /dev/null 2>&1 & echo $! > wcbot.pid &
