.DEFAULT_GOAL := build

build:
	go build
.PHONY: build

install:
	go install
.PHONY: install

watch:
	ginkgo watch
.PHONY: watch

test:
	go test
.PHONY: test

lint:
	golangci-lint run
.PHONY: lint

pack: clean linux darwin windows
.PHONY: pack

clean:
	rm -f tinfo *.exe *.tgz
.PHONY: clean

darwin: GOOS=darwin
darwin: GOARCH=amd64
darwin:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags '-w -s'
	tar cvzf tinfo-$(TRAVIS_TAG)-$(GOOS)-$(GOARCH).tgz tinfo
.PHONY: darwin

linux: GOOS=linux
linux: GOARCH=amd64
linux:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags '-w -s'
	tar cvzf tinfo-$(TRAVIS_TAG)-$(GOOS)-$(GOARCH).tgz tinfo
.PHONY: linux

windows: GOOS=windows
windows: GOARCH=amd64
windows:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags '-w -s'
	tar cvzf tinfo-$(TRAVIS_TAG)-$(GOOS)-$(GOARCH).tgz tinfo.exe
.PHONY: windows

ci: lint test
.PHONY: ci
