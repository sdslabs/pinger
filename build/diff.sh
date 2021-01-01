#!/usr/bin/env bash

source ./build/util.sh

set -e
set -x

DIFF_FILE="tmp.diff"

# Add files so we can diff any files if created
git add .

# Diff of specified files
git diff --staged > "${DIFF_FILE}"

if [ -s "${DIFF_FILE}" ]
then
  cat "${DIFF_FILE}"

  rm "${DIFF_FILE}"
  finally "Files differ"
  exit 1
fi

rm "${DIFF_FILE}"

finally "No difference between files"
