#!/bin/bash

export GO111MODULE=on

set -ex

go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.17.1
golangci-lint run ./...
for d in $(go list ./... | grep -v vendor | grep -v /cmd/ | grep -v test); do
    out=$(echo $d | cut -c21- | sed "s/\//_/g")
    rm -f coverage$out
    go test --coverprofile=coverage$out -covermode=atomic -race $d
done

