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
	@echo "🧹 Cleaning all artifacts..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	find . -type f -name "*.test" -delete
	find . -type f -name "*.out" -delete
	find . -type f -name "*.tmp" -delete
	find . -type f -name "*.temp" -delete
	find . -type f -name "*.log" -delete
	find . -type f -name "*.prof" -delete
	find . -type f -name "*.pprof" -delete
	rm -rf /tmp/fluffy-check
	@echo "✅ Cleanup complete!"

# Clean scripts directory outputs
clean-scripts:
	@echo "🧹 Cleaning scripts directory outputs..."
	rm -rf /tmp/fluffy-check
	@echo "✅ Scripts cleanup complete!"

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
	@echo "  build        - Build the binary"
	@echo "  test         - Run all tests"
	@echo "  test-race    - Run tests with race detector and coverage"
	@echo "  lint         - Run golangci-lint"
	@echo "  fmt          - Format code"
	@echo "  vet          - Run go vet"
	@echo "  clean        - Clean all artifacts"
	@echo "  clean-scripts - Clean scripts directory outputs"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  verify       - Verify dependencies"
	@echo "  check        - Run all checks (fmt, vet, lint, test)"
	@echo "  security     - Run gosec security scan"
	@echo "  help         - Show this help message"

# Default target
default: build
