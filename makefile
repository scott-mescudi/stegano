# Variables
PKG := ./...
TESTARGS ?= -v

.PHONY: all fmt test lint build clean help

# Default target
all: fmt lint test build ## Format, lint, test, and build the project

# Format Go code
fmt: ## Format the code
	@echo "Formatting Go code..."
	@go fmt $(PKG)

# Test the code
test: ## Run tests
	@echo "Running tests..."
	@go test $(TESTARGS) $(PKG)

# Lint the code
lint: ## Lint the code using go vet
	@echo "Linting code..."
	@go vet $(PKG)

# Build the library
build: ## Build the library
	@echo "Building the library..."
	@go build $(PKG)

# Clean the build cache
clean: ## Clean build cache
	@echo "Cleaning up..."
	@go clean

# Help menu
help: ## Display this help
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'
