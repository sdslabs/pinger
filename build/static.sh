#!/usr/bin/env bash

source ./build/util.sh

set -e
set -x

RESOURCE_DIR="./static"
BUNDLE_DIR="./pkg/util/static"

# Generate the resources file
go run -mod=readonly github.com/phogolabs/parcello/cmd/parcello \
 -r -d "${RESOURCE_DIR}" -b "${BUNDLE_DIR}"

finally "Files generated successfully in '${BUNDLE_DIR}/resource.go'"
