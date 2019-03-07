.DEFAULT_GOAL := build

NAME := tinfo

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

clean:
	rm -f tinfo *.exe *.tgz sha256sums.txt
.PHONY: clean

pack: clean
	GOOS=darwin GOARCH=amd64 go build -ldflags '-w -s' -o $(NAME)
	tar cvzf $(NAME)-$(TRAVIS_TAG)-darwin-amd64.tgz $(NAME)
	GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o $(NAME)
	tar cvzf $(NAME)-$(TRAVIS_TAG)-linux-amd64.tgz $(NAME)
	GOOS=windows GOARCH=amd64 go build -ldflags '-w -s' -o $(NAME).exe
	tar cvzf $(NAME)-$(TRAVIS_TAG)-windows-amd64.tgz $(NAME).exe
	sha256sum *.tgz > sha256sums.txt
.PHONY: pack

ci: lint test
.PHONY: ci
