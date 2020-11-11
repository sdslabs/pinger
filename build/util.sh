#!/usr/bin/env bash

# finally logs the output for when something finally happens.
function finally {
  set +x
  printf "\n\033[1m"
  printf "%s" "$1"
  printf "\033[0m\n\n"
  set -x
}
