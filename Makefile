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
NEXT_PATCH=$(shell if [ "$(CURRENT_VERSION)" = "v0.0.0" ]; then echo "0.0.1"; else echo $(CURRENT_VERSION) | awk -F. '{ printf("%d.%d.%d", substr($$1,2), $$2, $$3+1) }'; fi)
NEXT_MINOR=$(shell if [ "$(CURRENT_VERSION)" = "v0.0.0" ]; then echo "0.1.0"; else echo $(CURRENT_VERSION) | awk -F. '{ printf("%d.%d.%d", substr($$1,2), $$2+1, 0) }'; fi)
NEXT_MAJOR=$(shell if [ "$(CURRENT_VERSION)" = "v0.0.0" ]; then echo "1.0.0"; else echo $(CURRENT_VERSION) | awk -F. '{ printf("%d.%d.%d", substr($$1,2)+1, 0, 0) }'; fi)

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
		export VERSION=$(NEXT_PATCH); \
		echo "Using next patch version: $$VERSION"; \
		git tag -a "v$$VERSION" -m "Release v$$VERSION"; \
		git push origin "v$$VERSION"; \
	else \
		echo "Creating and pushing tag v$(VERSION)"; \
		git tag -a "v$(VERSION)" -m "Release v$(VERSION)"; \
		git push origin "v$(VERSION)"; \
	fi

# Bump version targets
bump-patch:
	@$(MAKE) release VERSION=$(NEXT_PATCH)

bump-minor:
	@$(MAKE) release VERSION=$(NEXT_MINOR)

bump-major:
	@$(MAKE) release VERSION=$(NEXT_MAJOR)

# Create a new release by creating a tag (CI will handle the release)
release: clean
	@if [ -z "$(VERSION)" ]; then \
		export VERSION=$(NEXT_PATCH); \
		echo "Using next patch version: $$VERSION"; \
		if [ -n "$$(git status --porcelain)" ]; then \
			echo "Error: Working directory is not clean. Please commit or stash changes first."; \
			exit 1; \
		fi; \
		$(MAKE) tag VERSION=$$VERSION; \
		echo "Tag v$$VERSION created and pushed. CI will handle the release."; \
	else \
		echo "Creating release v$(VERSION)"; \
		if [ -n "$$(git status --porcelain)" ]; then \
			echo "Error: Working directory is not clean. Please commit or stash changes first."; \
			exit 1; \
		fi; \
		$(MAKE) tag VERSION=$(VERSION); \
		echo "Tag v$(VERSION) created and pushed. CI will handle the release."; \
	fi

# Test release without publishing
release-dry-run: clean
	goreleaser release --snapshot --clean

# Run dry-run simulation for local development
dry-run: build
	@echo "Running dry-run simulation for local development..."
	@GITHUB_TOKEN="dummy_token" ./$(OUT_DIR)/$(BINARY_NAME) install --dry-run rpi02w