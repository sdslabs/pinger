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
		"bin" "Generates static resources and builds binary" \
		"build" "Builds the binary" \
		"build DEBUG=on" "Build the binary in debug mode" \
		"docker" "Builds the docker image" \
		"docker TAG=abc:def" "Builds the docker image with given tag" \
		"docs" "Builds the documentation" \
		"docs DEBUG=on" "Serves documentation on local server" \
		"fmt" "Formats the code" \
		"lint" "Checks for errors in code" \
		"proto" "Compiles the protobufs into go code" \
		"static" "Generate static resource go file" \
		"vendor" "Cleans up and updates vendor"
	@echo

.PHONY: all
all: bin docs

.PHONY: bin
bin: static build

.PHONY: build
build:
	DEBUG=$(DEBUG) ./build/build.sh

.PHONY: docker
docker:
	TAG=$(TAG) ./build/docker.sh

.PHONY: docs
docs:
	DEBUG=$(DEBUG) ./build/docs.sh

.PHONY: fmt
fmt:
	./build/fmt.sh

.PHONY: lint
lint:
	./build/lint.sh

.PHONY: proto
proto:
	./build/proto.sh

.PHONY: static
static:
	./build/static.sh

.PHONY: vendor
vendor:
	./build/vendor.sh
