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

# enable GPIO
sudo echo "1" > /sys/class/gpio/export 2>&1 || true

echo "Starting halloween player"
go run halloween_player.go
