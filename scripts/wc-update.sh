#!/bin/bash

echo "Updating Wing Commander..."
cd ${GOPATH}/src/github.com/BigOokie/skywire-wing-commander
git reset --hard
git pull origin master
go install -v  ./... 2>> /tmp/wingcommander_install_errors.log

echo "Checking for Wing Commander Process..."
cd ${GOPATH}/bin

WCPID=$(pgrep wcbot)

if [ "$WCPID" =  "" ]; then
    echo "Wing Commander does not appear to be running."
else
    echo "Wing Commander process identifed: $WCPID"
    echo "Terminating the process"
    kill  ${WCPID}
fi

echo ""
echo "Check Wing Commander Version."
./wcbot -v

echo "Starting Wing Commander (background)..."
nohup ./wcbot -upgradecompleted /dev/null 2>&1 & echo $! > wcbot.pid& > /dev/null

echo "Checking Wing Commander started..."
WCPID=$(pgrep wcbot)

if [ "$WCPID" = "" ]; then
    echo "Wing Commander does not appear to be running."
else
    echo "Wing Commander process identified: $WCPID"
    echo "Wing Commander updated and restarted"
fi