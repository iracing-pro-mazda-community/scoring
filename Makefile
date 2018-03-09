.PHONY: run binary setup glide test
SHELL := /bin/bash

all: run

run: binary
	./scoring

binary:
	GOOS=linux go build -i -o scoring

setup:
	go get -v -u github.com/Masterminds/glide

glide:
	glide install --force

test:
	GOARCH=amd64 GOOS=linux go test $$(go list ./... | grep -v /vendor/)
