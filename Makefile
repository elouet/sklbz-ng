.PHONY: all build run clean test lint

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary name
BINARY_NAME=sklbz-ng

all: build

build: 
	$(GOBUILD) -o $(BINARY_NAME) -v .

run: 
	$(GOBUILD) -o $(BINARY_NAME) -v . && ./$(BINARY_NAME)

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run with custom port
default:
	@echo "Running on default port 8080"
	$(GOBUILD) -o $(BINARY_NAME) -v . && ./$(BINARY_NAME) -port 8080

# Test the application
test: 
	$(GOTEST) -v ./...

# Download dependencies
deps: 
	$(GOMOD) download
	$(GOMOD) tidy

# Lint the code
lint: 
	golint ./...

# Run with custom configuration
run-dev:
	@echo "Running in development mode on port 8080"
	$(GOBUILD) -o $(BINARY_NAME) -v . && ./$(BINARY_NAME) -port 8080 -data ./data -static ./static

# Build for production
build-prod:
	@echo "Building for production"
	$(GOBUILD) -o $(BINARY_NAME) -ldflags="-s -w" .

# Install the application
install: 
	$(GOBUILD) -o $(GOPATH)/bin/$(BINARY_NAME) -v .

# Show help
help: 
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  run        - Build and run the application"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  deps       - Download dependencies"
	@echo "  lint       - Run linter"
	@echo "  run-dev    - Run in development mode"
	@echo "  build-prod - Build for production"
	@echo "  install    - Install the application"
	@echo "  help       - Show this help message"
