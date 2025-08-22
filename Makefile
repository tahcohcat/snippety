.PHONY: build test test-verbose test-coverage clean lint help

# Binary name
BINARY_NAME=snippety
CMD_PATH=./cmd/snippety

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(CMD_PATH)
	@echo "✅ Build complete: $(BINARY_NAME)"

# Run all tests
test:
	@echo "Running tests..."
	@go test ./...
	@echo "✅ All tests passed"

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v ./...

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover ./...
	@echo ""
	@echo "Generating detailed coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete"

# Run linters (if available)
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, running go vet and go fmt"; \
		go vet ./...; \
		go fmt ./...; \
	fi
	@echo "✅ Linting complete"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY_NAME)-linux-amd64 $(CMD_PATH)
	@GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY_NAME)-darwin-amd64 $(CMD_PATH)
	@GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY_NAME)-darwin-arm64 $(CMD_PATH)
	@GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY_NAME)-windows-amd64.exe $(CMD_PATH)
	@echo "✅ Multi-platform build complete"

# Development setup
dev: deps build test
	@echo "✅ Development setup complete"

# CI pipeline
ci: deps lint test-coverage
	@echo "✅ CI pipeline complete"

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the binary"
	@echo "  test          - Run all tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  lint          - Run linters"
	@echo "  deps          - Install dependencies"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  dev           - Complete development setup"
	@echo "  ci            - Run CI pipeline"
	@echo "  help          - Show this help message"

# Default target
default: help