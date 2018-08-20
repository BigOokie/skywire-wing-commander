#!/bin/bash

echo "Updating Wing Commander..."
cd ${GOPATH}/src/github.com/BigOokie/skywire-wing-commander
git reset --hard
git pull origin master
go install ./... 2>> /tmp/wingcommander_install_errors.log

echo "Kill Wing Commander Process..."
cd ${GOPATH}/bin
[[ -f wcbot.pid ]] && pkill -F wcbot.pid && rm wcbot.pid

echo "Restarting Wing Commander..."
./wcbot /dev/null 2>&1 & echo $! > wcbot.pid &

echo "Wing Commander updated and restarted"