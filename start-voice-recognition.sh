#!/usr/bin/env bash

(
cd voice-control || return
sudo chmod +0666 /dev/uinput && go install && go run main.go
)
