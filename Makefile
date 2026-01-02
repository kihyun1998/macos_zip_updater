# ============================================================================
# ACRA Updater Build Script - Go Version
# File: Makefile
# Description: Build universal binary for macOS updater
# ============================================================================

BINARY_NAME=ACRA_Point_Client_Updater
BIN_DIR=bin
AMD64_BINARY=$(BIN_DIR)/$(BINARY_NAME)_AMD64
ARM64_BINARY=$(BIN_DIR)/$(BINARY_NAME)_ARM64
UNIVERSAL_BINARY=$(BIN_DIR)/$(BINARY_NAME)

.PHONY: all build clean help

all: build

help:
	@echo "============================================================================"
	@echo "                    ACRA Updater Build Script"
	@echo "============================================================================"
	@echo "Usage: make [command]"
	@echo ""
	@echo "Commands:"
	@echo "  build    - Build universal binary (AMD64 + ARM64)"
	@echo "  clean    - Clean build artifacts"
	@echo "  help     - Show this help message"
	@echo "============================================================================"

build:
	@echo "============================================================================"
	@echo "Building ACRA Updater Universal Binary..."
	@echo "============================================================================"
	@mkdir -p $(BIN_DIR)
	@echo "[1/3] Building ARM64 binary..."
	GOOS=darwin GOARCH=arm64 go build -o $(ARM64_BINARY) main.go
	@echo "[2/3] Building AMD64 binary..."
	GOOS=darwin GOARCH=amd64 go build -o $(AMD64_BINARY) main.go
	@echo "[3/3] Creating Universal Binary..."
	lipo -create -output $(UNIVERSAL_BINARY) $(AMD64_BINARY) $(ARM64_BINARY)
	@echo "============================================================================"
	@echo "Build completed successfully!"
	@echo "Universal binary: $(UNIVERSAL_BINARY)"
	@echo "============================================================================"
	@lipo -info $(UNIVERSAL_BINARY)

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@echo "Done!"
