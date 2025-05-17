# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=personal-agent

# Project structure
MAIN_PACKAGE=./cmd/cli

all: test build

# Build the project
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v $(MAIN_PACKAGE)

# Run tests
test:
	$(GOTEST) -v ./...

# Run the application
run:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v $(MAIN_PACKAGE)
	./bin/$(BINARY_NAME)

# Clean build files
clean:
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)

# Vet examines Go source code and reports suspicious constructs
vet:
	$(GOCMD) vet ./...

.PHONY: all build test run clean vet



