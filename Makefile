.DEFAULT_GOAL := run
SHELL := /bin/bash
APP ?= $(shell basename $$(pwd))
COMMIT_SHA = $(shell git rev-parse --short HEAD)

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: run
## run: runs the application
run: build
	scripts/run.sh

.PHONY: dev
## dev: builds and runs the application
dev:
	scripts/dev.sh

.PHONY: build
## build: builds the application
build: clean
	@echo "Building binary ..."
	go build -o ${APP}

.PHONY: clean
## clean: cleans up binary files
clean:
	@echo "Cleaning up ..."
	@go clean

.PHONY: test
## test: runs go test with the race detector
test:
	@source .env; GOARCH=amd64 GOOS=linux go test -v -race ./...

.PHONY: init
## init: sets up go modules
init:
	@echo "Setting up modules ..."
	@go mod init 2>/dev/null; go mod tidy && go mod vendor

.PHONY: push
## push: pushes the application onto CF
push: test build
	cf push

.PHONY: setup
## setup: downloads glide and gin
setup:
	go get -v -u github.com/codegangsta/gin
	go get -v -u github.com/Masterminds/glide

.PHONY: glide
## glide: runs glide install
glide:
	glide install --force

.PHONY: connect
## connect: connects to postgres backend with CLI
connect:
	scripts/connect.sh
