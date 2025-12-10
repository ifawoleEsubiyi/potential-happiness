# Makefile for validatord

.PHONY: build test test-race lint fmt vet clean deps verify check security help default

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOMOD=$(GOCMD) mod
BINARY_NAME=validatord

# Build the binary
build:
	$(GOBUILD) -v -o $(BINARY_NAME) .

# Run all tests
test:
	$(GOTEST) -v ./...

# Run tests with race detector and coverage
test-race:
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic ./...

# Run golangci-lint
lint:
	golangci-lint run

# Format code
fmt:
	$(GOFMT) ./...
	gofmt -s -w .

# Run go vet
vet:
	$(GOVET) ./...

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Verify dependencies
verify:
	$(GOMOD) verify

# Run all checks (fmt, vet, lint, test)
check: fmt vet lint test

# Security scan with gosec
security:
	gosec ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  build     - Build the binary"
	@echo "  test      - Run all tests"
	@echo "  test-race - Run tests with race detector and coverage"
	@echo "  lint      - Run golangci-lint"
	@echo "  fmt       - Format code"
	@echo "  vet       - Run go vet"
	@echo "  clean     - Clean build artifacts"
	@echo "  deps      - Download and tidy dependencies"
	@echo "  verify    - Verify dependencies"
	@echo "  check     - Run all checks (fmt, vet, lint, test)"
	@echo "  security  - Run gosec security scan"
	@echo "  help      - Show this help message"

# Default target
default: build
