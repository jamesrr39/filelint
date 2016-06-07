export GOPATH := $(shell pwd)
export GOBIN := $GOPATH/bin

all: deps

deps:
	go get ./...

build:
	go build -o filelint

test:
	go test

