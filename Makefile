.PHONY: all build install clean test tidy run build-all deps

# Binary name
BINARY_NAME=jsmon

# Build directory
BUILD_DIR=bin

# Version info
VERSION?=1.0.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOINSTALL=$(GOCMD) install

all: tidy build

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Build the project
build: deps
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v ./cmd/jsmon

# Build for all platforms
build-all: deps
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)

	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/jsmon

	@echo "Building for Linux (arm64)..."
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/jsmon

	@echo "Building for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/jsmon

	@echo "Building for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/jsmon

	@echo "Building for Windows (amd64)..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/jsmon

	@echo "All builds completed!"

# Install the binary to GOPATH/bin
install: deps
	@echo "Installing $(BINARY_NAME)..."
	$(GOINSTALL) ./cmd/jsmon

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f jsmon.json
	rm -rf downloads/*

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Tidy up go.mod
tidy:
	@echo "Tidying go.mod..."
	$(GOMOD) tidy

# Help
help:
	@echo "JSMon - Makefile targets:"
	@echo "  make build       - Build the binary"
	@echo "  make build-all   - Build for all platforms (Linux, macOS, Windows)"
	@echo "  make install     - Install to GOPATH/bin"
	@echo "  make run         - Build and run"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make tidy        - Tidy go.mod"
	@echo "  make deps        - Download and install dependencies"
