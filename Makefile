.PHONY: test verify build release-component

VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse HEAD)
BUILT_AT ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
ENVIRONMENT ?= local
SERVICE ?= dinapay-connector-template
REPOSITORY ?= github.com/Germatic/dinapay-connector-template
MODULE = $(shell go list -m)
LDFLAGS = -s -w \
	-X $(MODULE)/internal/buildinfo.Service=$(SERVICE) \
	-X $(MODULE)/internal/buildinfo.Repository=$(REPOSITORY) \
	-X $(MODULE)/internal/buildinfo.Version=$(VERSION) \
	-X $(MODULE)/internal/buildinfo.Commit=$(COMMIT) \
	-X $(MODULE)/internal/buildinfo.BuiltAt=$(BUILT_AT) \
	-X $(MODULE)/internal/buildinfo.Environment=$(ENVIRONMENT)

test:
	go test ./...

verify:
	./scripts/verify-connector.sh

build:
	mkdir -p bin
	go build -trimpath -ldflags="$(LDFLAGS)" -o bin/connector ./cmd/connector

release-component: build
	PROVIDER="$(PROVIDER)" REPOSITORY="$(REPOSITORY)" ARTIFACT=bin/connector \
	PROCESS_NAME="$(PROCESS_NAME)" PORT="$(PORT)" ./scripts/render-release-component.sh
