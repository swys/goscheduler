PKGS := $(shell go list ./... | grep -v /vendor)
GOOS := linux
GOARC := amd64
BIN_NAME := go-scheduler
BIN_DIR := $(shell go env GOPATH)/bin

.PHONY: all
all: clean lint test build

.PHONY: clean
clean:
	rm -f $(BIN_NAME)

.PHONY: golangci-lint-install
golangci-lint-install:
	golangci-lint --version || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.44.2

.PHONY: lint
lint: golangci-lint-install
	golangci-lint run -v

.PHONY: test
test:
	go test -v -race

build:
	go build -a -o $(BIN_NAME) .