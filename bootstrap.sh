#!/usr/bin/env bash

# Script requires the stow command to be installed

# Exit when any command fails
set -e

# Wrap all commands in a subshell
(
mkdir -p ~/.gimx/config || true
mkdir -p ~/.gimx/macros || true
stow --target="$HOME" gimx
echo "Done"
)
