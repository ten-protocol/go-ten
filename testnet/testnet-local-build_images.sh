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
tools_path="${root_path}/tools"

parallel=true
# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --parallel)                 parallel=${value} ;;
            *)
    esac
done

if ${parallel} ;
  then
    echo "Running parallel builds with docker compose"
    ROOT_PATH=$root_path docker compose -f $testnet_path/docker-compose.local.yml build --parallel
    exit 0
fi


# run the builds in parallel - echo the full command to output
echo "Running parallel builds with regular docker"
command() {
    echo $@
     $( "$@" )
    echo $@ completed
}

command docker build -t testnetobscuronet.azurecr.io/obscuronet/eth2network:latest -f "${testnet_path}/eth2network.Dockerfile" "${root_path}" &
command docker build -t testnetobscuronet.azurecr.io/obscuronet/host:latest -f "${root_path}/dockerfiles/host.Dockerfile" "${root_path}" &
command docker build -t testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest -f "${tools_path}/hardhatdeployer/Dockerfile" "${root_path}" &
command docker build -t testnetobscuronet.azurecr.io/obscuronet/enclave:latest --build-arg TESTMODE=true -f "${root_path}/dockerfiles/enclave.Dockerfile" "${root_path}" &
command docker build -t testnetobscuronet.azurecr.io/obscuronet/tenscan:latest -f "${tools_path}/tenscan/Dockerfile" "${root_path}" &
command docker build -t testnetobscuronet.azurecr.io/obscuronet/faucet:latest -f "${tools_path}/faucet/Dockerfile" "${root_path}" &
command docker build -t testnetobscuronet.azurecr.io/obscuronet/obscuro_gateway:latest -f "${tools_path}/walletextension/Dockerfile" "${root_path}" &

wait

