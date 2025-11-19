.PHONY: build clean install test help

# Build configuration
BINARY_NAME=system-update
BUILD_DIR=.
CMD_DIR=cmd/system-update

# Detect Go binary location and set GOROOT if using Homebrew
GO := $(shell command -v go 2>/dev/null || echo "/opt/homebrew/bin/go")
GOROOT := $(shell $(GO) env GOROOT 2>/dev/null || echo "/opt/homebrew/opt/go/libexec")

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@GOROOT=$(GOROOT) $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build and install to /usr/local/bin
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@GOROOT=$(GOROOT) $(GO) test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	@GOROOT=$(GOROOT) $(GO) fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@GOROOT=$(GOROOT) $(GO) vet ./...

# Display help
help:
	@echo "Available targets:"
	@echo "  build    - Build the binary (default)"
	@echo "  install  - Build and install to /usr/local/bin"
	@echo "  clean    - Remove build artifacts"
	@echo "  test     - Run tests"
	@echo "  fmt      - Format code"
	@echo "  lint     - Run linter"
	@echo "  help     - Show this help message"

# Default target
.DEFAULT_GOAL := build
