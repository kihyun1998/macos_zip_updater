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
ENTITLEMENTS=updater.entitlements

# Code signing identity (set via environment variable or override)
SIGNING_IDENTITY=Developer ID Application: PENTA SYSTEMS TECHNOLOGY INC. (CJ2KJJN35D)

.PHONY: all build clean help sign

all: build

help:
	@echo "============================================================================"
	@echo "                    ACRA Updater Build Script"
	@echo "============================================================================"
	@echo "Usage: make [command]"
	@echo ""
	@echo "Commands:"
	@echo "  build    - Build universal binary (AMD64 + ARM64)"
	@echo "  sign     - Code sign the universal binary"
	@echo "  clean    - Clean build artifacts"
	@echo "  help     - Show this help message"
	@echo ""
	@echo "Environment Variables:"
	@echo "  SIGNING_IDENTITY - Code signing identity (default: Developer ID Application)"
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

sign:
	@echo "============================================================================"
	@echo "Code signing ACRA_Point_Client_Updater..."
	@echo "============================================================================"
	@BINARY="$${BINARY_PATH:-$(UNIVERSAL_BINARY)}"; \
	echo "Signing with: $(SIGNING_IDENTITY)"; \
	echo "Target binary: $$BINARY"; \
	echo "Entitlements: $(ENTITLEMENTS)"; \
	echo ""; \
	if [ ! -f "$$BINARY" ]; then \
		echo "ERROR: Binary not found: $$BINARY"; \
		exit 1; \
	fi; \
	if [ ! -f "$(ENTITLEMENTS)" ]; then \
		echo "ERROR: Entitlements file not found: $(ENTITLEMENTS)"; \
		exit 1; \
	fi; \
	echo "Removing existing signature..."; \
	codesign --remove-signature "$$BINARY" 2>/dev/null || true; \
	echo "Signing binary with entitlements..."; \
	codesign --force --timestamp --options runtime \
		--entitlements "$(ENTITLEMENTS)" \
		--sign "$(SIGNING_IDENTITY)" \
		"$$BINARY"; \
	echo "Verifying signature..."; \
	codesign --verify --deep --strict --verbose=2 "$$BINARY"; \
	echo "============================================================================"; \
	echo "Code signing completed successfully!"; \
	echo "============================================================================"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@echo "Done!"
