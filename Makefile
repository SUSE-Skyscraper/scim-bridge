.ONESHELL:
.PHONY: lint fmt build test

lint:
	golangci-lint run

fmt:
	go mod tidy
	go fmt ./cmd/... ./internal/...

build:
	go build -v ./cmd/main.go

test:
	go test -v ./cmd/... ./internal/... -coverprofile=coverage.out -covermode=atomic
