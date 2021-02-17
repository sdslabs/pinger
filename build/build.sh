#!/usr/bin/env bash

source ./build/util.sh

set -e
set -x

IMPORT_PATH="github.com/sdslabs/pinger"
TARGET_DIR="."
BIN_NAME="pinger"

# set the default version to undefined
if [ -z ${VERSION} ]
then
  CURRENT_TAG="undefined"
else
  CURRENT_TAG="${VERSION}"
fi

LDFLAGS="-X '${IMPORT_PATH}/cmd.version=${CURRENT_TAG}'"
if [ "${DEV}" = "on" ]
then
  TAGS="dev"
fi

# Create the target directory if it doesn't exist.
test -d "${TARGET_DIR}" || mkdir "${TARGET_DIR}"

# Finally build the binary.
go build -v \
  -tags "${TAGS}" \
  -ldflags "${LDFLAGS}" \
  -o "${TARGET_DIR}/${BIN_NAME}" \
  "${IMPORT_PATH}/cmd/${BIN_NAME}"

finally "Built as '${TARGET_DIR}/${BIN_NAME}'"
