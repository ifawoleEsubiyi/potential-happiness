# Makefile for validatord

.PHONY: build test test-race lint fmt vet clean clean-scripts deps verify check security help default

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
	@echo "Cleaning Go build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	@echo "Cleaning JavaScript artifacts..."
	find . -name "node_modules" -type d -exec rm -rf {} + 2>/dev/null || true
	find . -name "package-lock.json" -type f -delete 2>/dev/null || true
	find . -name "*.log" -type f -delete 2>/dev/null || true
	@echo "Cleaning temporary files..."
	find . -name "*.tmp" -type f -delete 2>/dev/null || true
	find . -name "*.temp" -type f -delete 2>/dev/null || true
	find . -name "*~" -type f -delete 2>/dev/null || true
	@echo "Clean complete!"

# Clean scripts directory specifically
clean-scripts:
	@echo "Cleaning scripts directory..."
	cd scripts && rm -f *.log *.tmp *.temp *~ 2>/dev/null || true
	@echo "Scripts directory cleaned!"

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
	@echo "  clean        - Clean all build artifacts (Go, JavaScript, temp files)"
	@echo "  clean-scripts - Clean scripts directory only"
	@echo "  deps      - Download and tidy dependencies"
	@echo "  verify    - Verify dependencies"
	@echo "  check     - Run all checks (fmt, vet, lint, test)"
	@echo "  security  - Run gosec security scan"
	@echo "  help      - Show this help message"

# Default target
default: build
