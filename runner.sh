#!/bin/bash

export GOROOT=/usr/local/go
export GOPATH=$HOME/Go
export PATH=$PATH:$GOROOT/bin

LOGDATEFORMAT='+%Y-%m-%d %H:%M:%S %p'

if ps axu | grep --silent "[h]alloween_player.go"; then
	echo "$(date "${LOGDATEFORMAT}") INFO [${0}] Halloween palyer (halloween_player.go) is already running, exiting"
	exit 0
fi

cd /home/pi/halloween_player

echo "Starting halloween player"

go run halloween_player.go
