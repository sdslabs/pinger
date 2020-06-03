# Prints help message
help:
	@echo "SDS Status Makefile"
	@echo "build   - Build status"
	@echo "format  - Format code using golangci-lint"
	@echo "help    - Prints help message"
	@echo "install - Install required tools"
	@echo "lint    - Lint code using golangci-lint"
	@echo "proto   - Build proto files"

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

# Generate proto files.
proto:
	@protoc \
	 --proto_path=pkg/proto/protobufs \
	 --go_out=plugins=grpc:pkg/proto \
	 pkg/proto/protobufs/messages.proto \
	 pkg/proto/protobufs/agent.proto \
	 pkg/proto/protobufs/central.proto

.PHONY: build format help install lint
