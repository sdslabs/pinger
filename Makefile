# Build status
build:
	@./scripts/build/build.sh

# Format code using golangci-lint
format:
	@./scripts/build/format.sh

# Prints help message
help:
	@echo "SDS Status Makefile"
	@echo "build   - Build status"
	@echo "format  - Format code using golangci-lint"
	@echo "help    - Prints help message"
	@echo "install - Install required tools"
	@echo "lint    - Lint code using golangci-lint"

# Install required tools
install:
	@./scripts/build/install.sh

# Lint code using golangci-lint
lint:
	@./scripts/build/lint.sh

# Setup SDS Status with config.
setup:
	@echo "[*] Setting up SDS Status"
	@cp _examples/sample.config.yml config.yml
	@echo "[+] Done! Edit the config.yml file and get started"

.PHONY: build format help install lint setup
