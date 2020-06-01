#!/usr/bin/env bash

wget -O SpeechSDK-Linux.tar.gz https://aka.ms/csspeech/linuxbinary tar --strip 1 -xzf SpeechSDK-Linux.tar.gz -C "$SPEECHSDK_ROOT"
