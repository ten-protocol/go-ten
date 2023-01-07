#!/usr/bin/env bash

#
# This script builds all images locally
#


# Ensure any fail is loud and explicit
set -euo pipefail

# Define local usage vars
start_path="$(cd "$(dirname "${0}")" && pwd)"
testnet_path="${start_path}"
root_path="${testnet_path}/.."

# run the builds in parallel - echo the full command to output
command() {
    echo $@ started
     $( "$@" )
    echo $@ completed
}

command docker build -t testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest -f "${testnet_path}/hardhatdeployer.Dockerfile" "${root_path}" &

wait

