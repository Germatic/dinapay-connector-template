.PHONY: test build

test:
	go test ./...

build:
	go build -o bin/connector ./cmd/connector
