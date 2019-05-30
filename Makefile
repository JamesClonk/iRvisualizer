.PHONY: run dev build push setup glide test connect
SHELL := /bin/bash

all: run

run: build
	scripts/run.sh

dev:
	scripts/dev.sh

build:
	rm -f iRvisualizer || true
	GOARCH=amd64 GOOS=linux go build -i -o iRvisualizer

push: test build
	cf push

setup:
	go get -v -u github.com/codegangsta/gin
	go get -v -u github.com/Masterminds/glide

glide:
	glide install --force

test:
	GOARCH=amd64 GOOS=linux go test $$(go list ./... | grep -v /vendor/)

connect:
	scripts/connect.sh
