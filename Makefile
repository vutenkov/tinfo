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

clean:
	rm -f tinfo *.exe *.tgz
.PHONY: clean

pack: clean
	GOOS=darwin GOARCH=amd64 go build -ldflags '-w -s' -o hello
	tar cvzf hello-$(TRAVIS_TAG)-darwin-amd64.tgz hello
	GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o hello
	tar cvzf hello-$(TRAVIS_TAG)-linux-amd64.tgz hello
	GOOS=windows GOARCH=amd64 go build -ldflags '-w -s' -o hello.exe
	tar cvzf hello-$(TRAVIS_TAG)-windows-amd64.tgz hello.exe
	sha256sum *.tgz > sha256sums.txt
.PHONY: pack

ci: lint test
.PHONY: ci
