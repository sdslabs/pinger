#!/usr/bin/env bash

source ./build/util.sh

set -e
set -x

DOCS_DIR="./docs"

cd "${DOCS_DIR}"

# Install dependencies. If they are already installed, this step exits
# almost immediately so no extra time is wasted.
bundle install

# Build the documentation. If debug mode is on, serve instead of build.
BUILD_CMD="build"
if [ "${DEBUG}" = "on" ]
then
  BUILD_CMD="serve"
fi
bundle exec jekyll "${BUILD_CMD}"

finally "Built in '${DOCS_DIR}/_site'"
