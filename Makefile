# Copyright (c) 2020 SDSLabs
# Use of this source code is governed by an MIT license
# details of which can be found in the LICENSE file.

.PHONY: build format help tools lint proto

.DEFAULT_GOAL := help

GO := go
GOPATH := $(shell go env GOPATH)
GOPATH_BIN := $(GOPATH)/bin
GOLANGCI_LINT := $(GOPATH_BIN)/golangci-lint
BUILD_OUTPUT := ./target/pinger
BUILD_INPUT := cmd/pinger/main.go
GO_PACKAGES := $(shell go list ./... | grep -v vendor)
UNAME := $(shell uname)

all: install lint proto build

help:
	@echo "Pinger Makefile"
	@echo "build   - Build pinger"
	@echo "format  - Format code using golangci-lint"
	@echo "help    - Prints help message"
	@echo "install - Install required tools"
	@echo "lint    - Lint code using golangci-lint"
	@echo "proto   - Build proto files"

build:
	@echo "Building..."
	@test -d target || mkdir target
	@$(GO) build -o $(BUILD_OUTPUT) $(BUILD_INPUT)
	@echo "Built as $(BUILD_OUTPUT)"

format:
	@echo "Formatting..."
	@$(GO) fmt $(GO_PACKAGES)
	@$(GOLANGCI_LINT) run --fix --issues-exit-code 0 > /dev/null 2>&1
	@echo "Code formatted"

lint:
	@echo "Linting..."
	@$(GO) vet $(GO_PACKAGES)
	@$(GOLANGCI_LINT) run
	@echo "No errors found"

proto:
	@echo "Compiling protobufs..."
	@protoc \
	 	--proto_path=pkg/proto/protobufs \
	 	--go_out=plugins=grpc:pkg/proto \
	 	pkg/proto/protobufs/messages.proto \
	 	pkg/proto/protobufs/agent.proto \
	 	pkg/proto/protobufs/central.proto
	@echo "Compiled successfully"

install: install-protoc install-golangcilint

install-protoc:
	@echo "Installing protoc..."
ifeq ($(UNAME), Darwin)
	@brew install protobuf
# TODO: add installation for ubuntu using apt or apt-get
else
	@echo "Install protoc manually, see: https://grpc.io/docs/protoc-installation/"
	@echo "Not required if not changing protobufs."
endif
	@echo "Installed successfully"

install-golangcilint:
	@echo "Installing golangci-lint..."
	@curl -sSfL \
	 	https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
	 	sh -s -- -b $(GOPATH_BIN) v1.24.0
	@echo "Installed successfully"
