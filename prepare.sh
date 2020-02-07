#!/usr/bin/env bash
set -e
go build
mkdir -p ./npm/bin/
cp spacex-cli npm/bin/