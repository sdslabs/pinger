#!/usr/bin/env bash

source ./build/util.sh

set -e
set -x

# If a tag is not specified, set the tag to "pinger:dev"
if [ -z "${TAG}" ]
then
  TAG="pinger:dev"
fi

# extract version from the tag by omitting pinger from string
VERSION=$(echo $TAG | cut -d':' -f 2)

# Build the image using the tag.
docker build -t "${TAG}" .

finally "Built image with tag '${TAG}'"
