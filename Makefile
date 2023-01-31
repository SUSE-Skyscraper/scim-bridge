.ONESHELL:
.SHELLFLAGS := -ec
SHELL := /bin/bash

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fmt
fmt:
	go mod tidy
	go fmt ./example/... ./v2/...

build:
	cd example
	go build -v ./cmd/main.go

.PHONY: test
test:
	go test -v ./example/... ./v2/... -coverprofile=coverage.out -covermode=atomic
