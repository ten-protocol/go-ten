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

# run the builds in parallel
command() {
    echo $4 start
     $( "$@" )
    echo $4 complete
}

command docker build -t testnetobscuronet.azurecr.io/obscuronet/obscuro_gethnetwork:latest -f "${testnet_path}/gethnetwork.Dockerfile" "${root_path}" &
command docker build -t testnetobscuronet.azurecr.io/obscuronet/obscuro_host:latest -f "${root_path}/dockerfiles/host.Dockerfile" "${root_path}" &
command docker build -t testnetobscuronet.azurecr.io/obscuronet/obscuro_contractdeployer:latest -f "${testnet_path}/contractdeployer.Dockerfile" "${root_path}" &
command docker build -t testnetobscuronet.azurecr.io/obscuronet/obscuro_enclave:latest -f "${root_path}/dockerfiles/enclave.Dockerfile" "${root_path}" &
command docker build -t testnetobscuronet.azurecr.io/obscuronet/obscuro_enclave_debug:latest -f "${root_path}/dockerfiles/enclave.debug.Dockerfile" "${root_path}" &

wait