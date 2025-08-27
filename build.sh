#!/usr/bin/env bash
set -ex

go vet ./...
rm -f cmserver
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cmserver