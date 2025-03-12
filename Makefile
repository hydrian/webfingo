.PHONY: build compile test

build: test compile

test:
	go test -v ./...

compile:
	go build -o ./bin/webfingo ./cmd/webfingo.go

