.PHONY: all test build clean run deps build-linux

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BINARY_NAME=link
BINARY_UNIX=$(BINARY_NAME)_unix
OUT_DIR=out

# Run all
all: test build

# Run tests
test: 
	$(GOTEST) 

# Build the project
build: 
	$(GOBUILD) -o $(OUT_DIR)/$(BINARY_NAME) ./cmd/${BINARY_NAME} 

# Clean build files
clean: 
	$(GOCLEAN)
	rm -f $(OUT_DIR)/$(BINARY_NAME)
	rm -f $(OUT_DIR)/$(BINARY_UNIX)

# Run the application
run: build
	./$(OUT_DIR)/$(BINARY_NAME)

# Install dependencies
deps:
	$(GOGET) 

# Cross-compilation for Linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(OUT_DIR)/$(BINARY_UNIX) 