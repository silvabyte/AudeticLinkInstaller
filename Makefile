.PHONY: all test build clean run deps build-linux release dry-run tag get-version bump-major bump-minor bump-patch

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

# Version management
CURRENT_VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
NEXT_PATCH=$(shell echo $(CURRENT_VERSION) | awk -F. '{ printf("v%d.%d.%d", $$1, $$2, $$3+1) }' | sed 's/v\([0-9]\)/\1/')
NEXT_MINOR=$(shell echo $(CURRENT_VERSION) | awk -F. '{ printf("v%d.%d.%d", $$1, $$2+1, 0) }' | sed 's/v\([0-9]\)/\1/')
NEXT_MAJOR=$(shell echo $(CURRENT_VERSION) | awk -F. '{ printf("v%d.%d.%d", $$1+1, 0, 0) }' | sed 's/v\([0-9]\)/\1/')

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
	rm -rf dist/

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

# Get current version
get-version:
	@echo "Current version: $(CURRENT_VERSION)"
	@echo "Next patch: $(NEXT_PATCH)"
	@echo "Next minor: $(NEXT_MINOR)"
	@echo "Next major: $(NEXT_MAJOR)"

# Create and push a new tag
tag:
	@if [ -z "$(VERSION)" ]; then \
		echo "No version specified, using next patch version: $(NEXT_PATCH)"; \
		VERSION=$(NEXT_PATCH); \
	fi
	@echo "Creating and pushing tag v$(VERSION)"
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)

# Bump version targets
bump-patch: VERSION=$(NEXT_PATCH)
bump-patch: release

bump-minor: VERSION=$(NEXT_MINOR)
bump-minor: release

bump-major: VERSION=$(NEXT_MAJOR)
bump-major: release

# Create a new release using GoReleaser
release: clean
	@if [ -z "$(VERSION)" ]; then \
		echo "No version specified, using next patch version: $(NEXT_PATCH)"; \
		VERSION=$(NEXT_PATCH); \
	fi
	@echo "Creating release v$(VERSION)"
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "Error: Working directory is not clean. Please commit or stash changes first."; \
		exit 1; \
	fi
	make tag VERSION=$(VERSION)
	goreleaser release --clean

# Test release without publishing
release-dry-run: clean
	goreleaser release --snapshot --clean

# Run dry-run simulation for local development
dry-run: build
	@echo "Running dry-run simulation for local development..."
	@GITHUB_TOKEN="dummy_token" ./$(OUT_DIR)/$(BINARY_NAME) install --dry-run rpi02w