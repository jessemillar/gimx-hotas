#!/usr/bin/env bash

(
cd voice-recognition || return
sudo chmod +0666 /dev/uinput
go build main.go

# Configure the Azure Speech API
SPEECHSDK_ROOT="$HOME/.speechsdk"
CGO_CFLAGS="-I$SPEECHSDK_ROOT/include/c_api"
CGO_LDFLAGS="-L$SPEECHSDK_ROOT/lib/arm32 -lMicrosoft.CognitiveServices.Speech.core"
LD_LIBRARY_PATH="$SPEECHSDK_ROOT/lib/arm32:$LD_LIBRARY_PATH"

sudo SPEECHSDK_ROOT="$SPEECHSDK_ROOT" CGO_CFLAGS="$CGO_CFLAGS" CGO_LDFLAGS="$CGO_LDFLAGS" LD_LIBRARY_PATH="$LD_LIBRARY_PATH" ./main
)
