#!/bin/sh
# Auto Startup Wing Commander 

export GOPATH=$HOME/go

cd $GOPATH/bin
nohup ./wcbot /dev/null 2>&1 & echo $! > wcbot.pid& > /dev/null
