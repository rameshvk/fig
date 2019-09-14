#!/bin/bash

GO111MODULE=on

go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.17.1 &&
    go test ./... &&
    golangci-lint run ./...
