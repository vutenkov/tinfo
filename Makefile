.PHONY: build install test watch lint
.DEFAULT_GOAL := build

build:
	go build

release:
	go build -ldflags "-s -w"

install:
	go install

test:
	go test

watch:
	ginkgo watch

lint:
	golangci-lint run
