#!/bin/bash

cd ${GOPATH}/src/github.com/skycoin/skywire
git reset --hard > /dev/null 2>&1
git clean -f -d > /dev/null 2>&1
git pull origin master > /dev/null 2>&1


echo "Updating Wing Commander..."
cd ${GOPATH}/src/github.com/BigOokie/skywire-wing-commander
git reset --hard > /dev/null 2>&1
git clean -f -d > /dev/null 2>&1
git pull origin master > /dev/null 2>&1
go install ./... 2>> /tmp/wingcommander_install_errors.log

echo "Kill Wing Commander Process..."
cd ${GOPATH}/bin
pkill -F wcbot.pid
rm wcbot.pid

echo "Restarting Wing Commander..."
nohup ./wcbot /dev/null 2>&1 & echo $! > wcbot.pid& > /dev/null

echo "Wing Commander updated and restarted"

exit 0