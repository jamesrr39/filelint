export GOPATH := $(shell pwd)
export GOBIN := $(GOPATH)/bin

.PHONY: deps build test
all: deps

deps:
	go get ./...

build: deps
	go build -o filelint

test: deps
	go test

