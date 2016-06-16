export GOPATH := $(shell pwd)
export GOBIN := $(GOPATH)/bin

.PHONY: deps build test
all: deps

deps:
	go get -t ./...

build: deps
	go build -o filelint

test: deps
	go test ./...

