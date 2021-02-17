BOLD := \033[1m
NORMAL := \033[0m

.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "$(BOLD)Usage:$(NORMAL)"
	@echo "  make [target...] [KEY1=VAL1 KEY2=VAL2]"
	@echo ""
	@echo "$(BOLD)Targets:$(NORMAL)"
	@printf "  make %s\n    %s\n" \
		"all" "Builds binary and documentation" \
		"build" "Builds the binary" \
		"build DEV=on" "Build the binary in development mode" \
		"build VERSION=v1.2.3" "Builds the binary with version v1.2.3" \
		"docker" "Builds the docker image" \
		"docker TAG=abc:def" "Builds the docker image with given tag" \
		"docs" "Builds the documentation" \
		"docs DEV=on" "Serves documentation on local server" \
		"fmt" "Formats the code" \
		"lint" "Checks for errors in code" \
		"proto" "Compiles the protobufs into go code" \
		"vendor" "Cleans up and updates vendor"
	@echo

.PHONY: all
all: build docs

.PHONY: build
build:
	DEV=$(DEV) VERSION=$(VERSION) ./build/build.sh

.PHONY: docker
docker:
	TAG=$(TAG) ./build/docker.sh

.PHONY: docs
docs:
	DEV=$(DEV) ./build/docs.sh

.PHONY: fmt
fmt:
	./build/fmt.sh

.PHONY: lint
lint:
	./build/lint.sh

.PHONY: proto
proto:
	./build/proto.sh

.PHONY: vendor
vendor:
	./build/vendor.sh
