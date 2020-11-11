#!/usr/bin/env bash

source ./build/util.sh

set -e
set -x

# Vet the packages before running golangci-lint since it fails with very
# ambiguous error messages.
# shellcheck disable=SC2086
go vet ./...

# Now run golangci-lint.
golangci-lint run

finally "No errors found"
