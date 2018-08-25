#!/bin/bash
echo "Updating Wing Commander..."
cd ${GOPATH}/src/github.com/BigOokie/skywire-wing-commander
git reset --hard
git pull origin master
go install -v  ./... 2>> /tmp/wingcommander_install_errors.log

echo "Checking if Wing Commander is already running..."
cd ${GOPATH}/bin

WCPID=$(pgrep wcbot)

if [ "$WCPID" =  "" ]; then
    echo "Wing Commander does not appear to be running."
else
    echo "Wing Commander process identified: $WCPID"
    echo "Terminating the process."
    kill  ${WCPID}
fi

echo ""
echo "New version:"
./wcbot -v

echo "Starting Wing Commander (background)..."
nohup ./wcbot /dev/null 2>&1 & echo $! > wcbot.pid& > /dev/null

echo "Checking Wing Commander started..."
WCPID=$(pgrep wcbot)

if [ "$WCPID" = "" ]; then
    echo "Wing Commander does not appear to be running."
    echo "Something may have gone wrong. Please restart manually."
else
    echo "Wing Commander process identifed: $WCPID"
    echo "Wing Commander updated and restarted successfully."
fi