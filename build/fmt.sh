#!/usr/bin/env bash

source ./build/util.sh

set -e
set -x

# At-least format using gofmt.
# shellcheck disable=SC2086
go fmt ./...

# Now try running golangci-lint. This fails if the build fails in any way.
golangci-lint run --fix

finally "Formatted successfully"
