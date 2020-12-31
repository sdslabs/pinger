#!/usr/bin/env bash

source ./build/util.sh

set -e
set -x

IMPORT_PATH="github.com/sdslabs/pinger"
TARGET_DIR="."
BIN_NAME="pinger"

# If the tag is not found takes the "{{ HEAD commit short name }} (dev)"
CURRENT_TAG=$(git describe --tags --exact-match 2> /dev/null \
  || echo "$(git rev-parse --short HEAD)" "(dev)")

LDFLAGS="-X '${IMPORT_PATH}/cmd.version=${CURRENT_TAG}'"
if [ "${DEBUG}" = "on" ]
then
  LDFLAGS="${LDFLAGS} -X '${IMPORT_PATH}/cmd.debug=true'"
fi

# Create the target directory if it doesn't exist.
test -d "${TARGET_DIR}" || mkdir "${TARGET_DIR}"

# Finally build the binary.
go build -v \
  -ldflags "${LDFLAGS}" \
  -o "${TARGET_DIR}/${BIN_NAME}" \
  "${IMPORT_PATH}/cmd/${BIN_NAME}"

finally "Built as '${TARGET_DIR}/${BIN_NAME}'"
