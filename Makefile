PROJECT = lorem
OUTPUT_DIR = build
INSTALL_DIR = /usr/local/bin

# List of OS and architecture pairs to build for
OS_ARCH_PAIRS = \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64

# List of unsupported OS/ARCH pairs for UPX compression
UNSUPPORTED_UPX_PAIRS = \
	darwin/amd64 \
	darwin/arm64 \
	windows/arm64

.PHONY: all build clean help

# Default target
all: clean build

# Determine output name based on OS and architecture
define build_target
	OS=$(1); \
	ARCH=$(2); \
	OUTPUT_NAME=$(OUTPUT_DIR)/$(PROJECT)-$$OS-$$ARCH; \
	if [ "$$OS" = "windows" ]; then OUTPUT_NAME=$$OUTPUT_NAME.exe; fi; \
	echo "Building $$OUTPUT_NAME..."; \
	GOOS=$$OS GOARCH=$$ARCH go build -o $$OUTPUT_NAME; \
	$(call compress_target,$$OUTPUT_NAME,$$OS,$$ARCH)
endef

# Compress the binary if supported
define compress_target
	OUTPUT_NAME=$(1); \
	OS=$(2); \
	ARCH=$(3); \
	if echo "$(UNSUPPORTED_UPX_PAIRS)" | grep -q "$$OS/$$ARCH"; then \
		echo "Skipping UPX compression for $$OUTPUT_NAME (unsupported platform)"; \
	else \
		echo "Compressing $$OUTPUT_NAME with UPX..."; \
		upx --best $$OUTPUT_NAME || echo "UPX compression failed for $$OUTPUT_NAME"; \
	fi
endef

# Build the project for all architectures
build:
	@echo "Building ${PROJECT} for multiple platforms..."
	@mkdir -p $(OUTPUT_DIR)
	@for os_arch in $(OS_ARCH_PAIRS); do \
		OS=$$(echo $$os_arch | cut -d'/' -f1); \
		ARCH=$$(echo $$os_arch | cut -d'/' -f2); \
		$(call build_target,$$OS,$$ARCH); \
	done
	@echo "All builds and compressions completed."

# Clean the built files
clean:
	@echo "Cleaning built files..."
	@rm -rf $(OUTPUT_DIR)
	@echo "Done"

# Display help screen
help:
	@echo "\nUsage: make [target]\n"
	@echo "Targets:\n"
	@echo "  all         : Clean and build the project"
	@echo "  build       : Build the project for all architectures and compress with UPX"
	@echo "  clean       : Remove the built executables"
	@echo "  help        : Show this help message\n"
