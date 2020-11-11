#!/usr/bin/env bash

source ./build/util.sh

set -e
set -x

# Tidy up and then vendor deps.
go mod tidy
go mod vendor

finally "Vendor updated"
