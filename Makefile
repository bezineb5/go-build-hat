# Simple Makefile for go-build-hat project

.PHONY: test build-examples help

help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  test - Run tests"
	@echo "  build-examples - Build all examples"

# Run tests
test:
	go test ./pkg/buildhat/...

# Build all examples
build-examples:
	@mkdir -p bin
	go build -o bin/real_hardware ./examples/real_hardware/
