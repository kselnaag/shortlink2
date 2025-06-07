#!/usr/bin/env bash
set -xe
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ./bin/shortlink2 ./cmd/main.go
