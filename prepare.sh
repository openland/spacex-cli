#!/usr/bin/env bash
set -e
GOOS=linux GOARCH=amd64 go build -o spacex-cli-linux
GOOS=darwin GOARCH=amd64 go build -o spacex-cli-macos
# GOOS=windows GOARCH=amd64 go build -o spacex-cli-windows.exe
mkdir -p ./npm/bin/
cp spacex-cli-linux npm/bin/
cp spacex-cli-macos npm/bin/
# cp spacex-cli-windows.exe npm/bin/