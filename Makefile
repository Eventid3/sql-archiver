.PHONY: build test test-verbose test-cover help install

BUILD_DIR = build
EXE = sql-archiver

help:
	@echo "Available commands:"
	@echo "  build         - Compile the project and place the binary in the build directory"
	@echo "  test          - Run all tests"
	@echo "  test-verbose  - Run all tests with verbose output"
	@echo "  test-cover    - Generate a test coverage report"
	@echo "  install       - Build the project and install the binary to the Go bin directory"

build:
	go build -o $(BUILD_DIR)/$(EXE) .

test:
	go test ./...

test-verbose:
	go test -v ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report written to coverage.html"

install: build
	cp $(BUILD_DIR)/$(EXE) $(HOME)/go/bin
