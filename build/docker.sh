#!/usr/bin/env bash

source ./build/util.sh

set -e
set -x

# If a tag is not specified, set the tag to "pinger:dev"
if [ -z "${TAG}" ]
then
  TAG="pinger:dev"
fi

# Build the image using the tag.
docker build -t "${TAG}" .

finally "Built image with tag '${TAG}'"
