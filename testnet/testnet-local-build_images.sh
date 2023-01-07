#!/usr/bin/env bash

#
# This script builds all images locally - note the script must be executed
# in the go-obscuro/testnet folder
#


# Ensure any fail is loud and explicit
set -euo pipefail

# Define local usage vars
start_path="$(cd "$(dirname "${0}")" && pwd)"
testnet_path="${start_path}"
root_path="${testnet_path}/.."

# Run parallel builds
ROOT_PATH=$root_path docker compose -f $testnet_path/docker-compose.local.yml build --parallel




