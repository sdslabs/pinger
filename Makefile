# Copyright (c) 2020 SDSLabs
# Use of this source code is governed by an MIT license
# details of which can be found in the LICENSE file.

.PHONY: build docker docs format help install lint proto vendor

.DEFAULT_GOAL := help

GO := go
GOPATH := $(shell go env GOPATH)
GOPATH_BIN := $(GOPATH)/bin
GOLANGCI_LINT := $(GOPATH_BIN)/golangci-lint
BUILD_OUTPUT := ./target/pinger
BUILD_INPUT := cmd/pinger/main.go
GO_PACKAGES := $(shell go list ./... | grep -v vendor)
UNAME := $(shell uname)

all: install lint proto docs build docker

help:
	@echo "Pinger Makefile"
	@echo "build   - Build pinger"
	@echo "docker  - Build docker image"
	@echo "docs    - Build documentation"
	@echo "format  - Format code using golangci-lint"
	@echo "help    - Prints help message"
	@echo "install - Install required tools"
	@echo "lint    - Lint code using golangci-lint"
	@echo "proto   - Build proto files"
	@echo "vendor  - Vendor dependencies and tidy up"

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

vendor:
	@echo "Tidy up go.mod..."
	@$(GO) mod tidy
	@echo "Vendoring..."
	@$(GO) mod vendor
	@echo "Done!"

proto:
	@echo "Compiling protobufs..."
	@protoc \
	 	--proto_path=pkg/components/agent/proto/protobufs \
	 	--go_out=plugins=grpc:pkg/components/agent/proto \
	 	pkg/components/agent/proto/protobufs/messages.proto \
	 	pkg/components/agent/proto/protobufs/agent.proto 
	@echo "Compiled successfully"

install: install-protoc install-golangcilint

ifeq ($(TAG),)
TAG := "pinger:dev"
endif
docker:
	@echo "Building docker image..."
	@docker build -t $(TAG) .
	@echo "Built with tag $(TAG)"

install-protoc:
	@echo "Installing protoc..."
ifeq ($(UNAME), Darwin)
	@brew install protobuf
# TODO: add installation for ubuntu using apt or apt-get
else
	@echo "Install protoc manually, see: https://grpc.io/docs/protoc-installation/"
	@echo "Not required if not changing protobufs."
endif
	@echo "Installing protoc-gen-go"
	@go get github.com/golang/protobuf/protoc-gen-go
	@echo "Installed successfully"

install-golangcilint:
	@echo "Installing golangci-lint..."
	@curl -sSfL \
	 	https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
	 	sh -s -- -b $(GOPATH_BIN) v1.31.0
	@echo "Installed successfully"

docs: docs-install docs-build

docs-install:
	@echo "Installing documentation dependencies"
	@cd docs && bundle install
	@echo "Dependencies installed!"

docs-build:
	@echo "Building documentation..."
	@cd docs && bundle exec jekyll build
	@echo "Built into ./docs/_site"

docs-watch:
	@echo "Building documentation in watch mode..."
	@cd docs && bundle exec jekyll build --watch

docs-serve:
	@echo "Serving documentation on :4000"
	@cd docs && bundle exec jekyll serve
