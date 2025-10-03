# Build variables
BINARY_NAME=linkedin-scraper
BUILD_DIR=bin
GO_FILES=$(shell find . -name "*.go" -type f)

# Default target
.PHONY: all
all: clean build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/scraper/main.go

# Build for different platforms
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux cmd/scraper/main.go

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows.exe cmd/scraper/main.go

.PHONY: build-mac
build-mac:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-mac cmd/scraper/main.go

# Run the application
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) -profile "$(PROFILE_URL)" -config configs/config.yaml

# Run with example profile
.PHONY: run-example
run-example: build
	@echo "Running with example profile..."
	./$(BUILD_DIR)/$(BINARY_NAME) -profile "https://www.linkedin.com/in/avasantlal/" -config configs/config.yaml

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Test with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run

# Setup development environment
.PHONY: setup
setup:
	@echo "Setting up development environment..."
	go mod init linkedin-scraper || true
	go mod tidy
	mkdir -p data/output
	mkdir -p logs
	mkdir -p configs
	@echo "Setup complete! Don't forget to:"
	@echo "1. Install ChromeDriver"
	@echo "2. Update configs/config.yaml with your credentials"
	@echo "3. Run 'make build' to build the application"

# Create directories
.PHONY: dirs
dirs:
	@echo "Creating directories..."
	mkdir -p data/output
	mkdir -p logs
	mkdir -p bin

# Install ChromeDriver (macOS)
.PHONY: install-chromedriver-mac
install-chromedriver-mac:
	@echo "Installing ChromeDriver for macOS..."
	brew install chromedriver

# Install ChromeDriver (Linux)
.PHONY: install-chromedriver-linux
install-chromedriver-linux:
	@echo "Installing ChromeDriver for Linux..."
	@echo "Please install ChromeDriver manually:"
	@echo "1. Download from https://chromedriver.chromium.org/"
	@echo "2. Extract and place in /usr/local/bin/"
	@echo "3. Make executable: chmod +x /usr/local/bin/chromedriver"

# Run with different output formats
.PHONY: run-csv
run-csv: build
	./$(BUILD_DIR)/$(BINARY_NAME) -profile "$(PROFILE_URL)" -format csv -output data/output/profile.csv

.PHONY: run-json
run-json: build
	./$(BUILD_DIR)/$(BINARY_NAME) -profile "$(PROFILE_URL)" -format json -output data/output/profile.json

.PHONY: run-both
run-both: build
	./$(BUILD_DIR)/$(BINARY_NAME) -profile "$(PROFILE_URL)" -format both -output data/output/profile.csv

# Development helpers
.PHONY: dev-run
dev-run:
	@echo "Running in development mode..."
	go run cmd/scraper/main.go -profile "$(PROFILE_URL)" -config configs/config.yaml

# Show help
.PHONY: help
help:
	@echo "LinkedIn Profile Scraper - Makefile Commands"
	@echo ""
	@echo "Build Commands:"
	@echo "  build              Build the application for current platform"
	@echo "  build-linux        Build for Linux"
	@echo "  build-windows      Build for Windows"
	@echo "  build-mac          Build for macOS"
	@echo ""
	@echo "Run Commands:"
	@echo "  run PROFILE_URL=<url>    Run the scraper with specified profile URL"
	@echo "  run-example              Run with example profile"
	@echo "  run-csv PROFILE_URL=<url> Run and export to CSV"
	@echo "  run-json PROFILE_URL=<url> Run and export to JSON"
	@echo "  dev-run PROFILE_URL=<url> Run in development mode"
	@echo ""
	@echo "Development Commands:"
	@echo "  setup              Setup development environment"
	@echo "  deps               Install dependencies"
	@echo "  fmt                Format code"
	@echo "  lint               Lint code"
	@echo "  test               Run tests"
	@echo "  test-coverage      Run tests with coverage"
	@echo ""
	@echo "Utility Commands:"
	@echo "  clean              Clean build artifacts"
	@echo "  dirs               Create necessary directories"
	@echo "  install-chromedriver-mac   Install ChromeDriver on macOS"
	@echo "  help               Show this help message"
	@echo ""
	@echo "Example Usage:"
	@echo "  make run PROFILE_URL=\"https://www.linkedin.com/in/username/\""