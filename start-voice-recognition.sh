#!/usr/bin/env bash

(
cd voice-recognition || return
# Allow the GIMX adapter to work
sudo chmod +0666 /dev/uinput

# Configure the Azure Speech API
export SPEECHSDK_ROOT="/home/pi/.speechsdk"
export CGO_CFLAGS="-I$SPEECHSDK_ROOT/include/c_api"
export CGO_LDFLAGS="-L$SPEECHSDK_ROOT/lib/arm32 -lMicrosoft.CognitiveServices.Speech.core"
export LD_LIBRARY_PATH="$SPEECHSDK_ROOT/lib/arm32:$LD_LIBRARY_PATH"

# Build binary
CGO_CFLAGS="$CGO_CFLAGS" CGO_LDFLAGS="$CGO_LDFLAGS" LD_LIBRARY_PATH="$LD_LIBRARY_PATH" go build main.go

# Run binary
sudo -E CGO_CFLAGS="$CGO_CFLAGS" CGO_LDFLAGS="$CGO_LDFLAGS" LD_LIBRARY_PATH="$LD_LIBRARY_PATH" ./main
)
