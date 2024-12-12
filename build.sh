#!/usr/bin/env bash
VERSION=$(git describe --tags)
export GOOS="linux"
export GOARCH="arm"
export GOARM="7"

go build -v -ldflags="-X 'main.Version=${VERSION}'" -o dist/server-monitor