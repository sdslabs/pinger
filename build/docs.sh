#!/usr/bin/env bash

source ./build/util.sh

set -e
set -x

DOCS_DIR="./docs"

cd "${DOCS_DIR}"

# Build the documentation. If debug mode is on, serve instead of build.
# This also installs dependencies when debug mode is off. When running
# for the first time, don't run in debug mode to install all dependencies.
BUILD_CMD="build"
if [ "${DEBUG}" = "on" ]
then
  BUILD_CMD="serve"
else
  bundle install
fi
bundle exec jekyll "${BUILD_CMD}"

finally "Built in '${DOCS_DIR}/_site'"
