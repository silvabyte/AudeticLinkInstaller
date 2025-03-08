.PHONY: all test build clean run deps build-linux release dry-run

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BINARY_NAME=link
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_MAC=$(BINARY_NAME)_mac
BINARY_WIN=$(BINARY_NAME).exe
BINARY_ARM64=$(BINARY_NAME)_arm64
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
	rm -f $(OUT_DIR)/$(BINARY_MAC)
	rm -f $(OUT_DIR)/$(BINARY_WIN)
	rm -f $(OUT_DIR)/$(BINARY_ARM64)

# Run the application
run: build
	./$(OUT_DIR)/$(BINARY_NAME)

# Install dependencies
deps:
	$(GOGET) 

# Cross-compilation for Linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(OUT_DIR)/$(BINARY_UNIX) ./cmd/${BINARY_NAME}

# Cross-compilation for Raspberry Pi (ARM64)
build-arm64:
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(OUT_DIR)/$(BINARY_ARM64) ./cmd/${BINARY_NAME}

# Cross-compilation for all platforms
build-all:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(OUT_DIR)/$(BINARY_MAC) ./cmd/${BINARY_NAME}
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(OUT_DIR)/$(BINARY_UNIX) ./cmd/${BINARY_NAME}
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(OUT_DIR)/$(BINARY_WIN) ./cmd/${BINARY_NAME}
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(OUT_DIR)/$(BINARY_ARM64) ./cmd/${BINARY_NAME}

# Create a new release
release: build-arm64
	@if [ -z "$(VERSION)" ]; then echo "Please specify VERSION=x.x.x when running make release"; exit 1; fi
	@echo "Creating release $(VERSION)"
	gh release create v$(VERSION) \
		$(OUT_DIR)/$(BINARY_ARM64) \
		--title "Release $(VERSION)" \
		--generate-notes

# Run dry-run simulation for local development
dry-run: build
	@echo "Running dry-run simulation for local development..."
	@GITHUB_TOKEN="dummy_token" ./$(OUT_DIR)/$(BINARY_NAME) install --dry-run rpi02w