#!/usr/bin/env bash

source ./build/util.sh

set -e
set -x

DOCS_DIR="./docs"

cd "${DOCS_DIR}"

# Build the documentation. If development mode is on, serve instead of build.
BUILD_CMD="build"
if [ "${DEV}" = "on" ]
then
  BUILD_CMD="serve"
fi
mdbook "${BUILD_CMD}"

finally "Built in '${DOCS_DIR}/book'"
