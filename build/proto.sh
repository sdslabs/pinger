#!/usr/bin/env bash

source ./build/util.sh

set -e
set -x

PROTO_PATH="pkg/components/agent/proto/protobufs"
GO_OUT="pkg/components/agent/proto"
PROTO_FILES=$(echo ${PROTO_PATH}/*.proto)

# Compile proto files into Go code.
# shellcheck disable=SC2086
protoc \
	"--proto_path=${PROTO_PATH}" \
	"--go_out=plugins=grpc:${GO_OUT}" \
	${PROTO_FILES}

finally "Compiled in '${GO_OUT}'"
