# Prints help message
help:
	@echo "SDS Status Makefile"
	@echo "build   - Build status"
	@echo "format  - Format code using golangci-lint"
	@echo "help    - Prints help message"
	@echo "install - Install required tools"
	@echo "lint    - Lint code using golangci-lint"

# Build status
build:
	@./scripts/build.sh

# Format code using golangci-lint
format:
	@./scripts/format.sh

# Install required tools
install:
	@./scripts/install.sh

# Lint code using golangci-lint
lint:
	@./scripts/lint.sh

.PHONY: build format help install lint
