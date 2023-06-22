#!/usr/bin/env bash

#
# This script starts Obscuroscan for testnet
#

# Ensure any fail is loud and explicit
set -euo pipefail

# Define defaults
nodeID=""
rpcServerAddress="http://testnet.obscu.ro:13000"
address="127.0.0.1:3000"
logPath="obscuroscan_logs.txt"
port=8080
image="testnetobscuronet.azurecr.io/obscuronet/obscuroscan_testnet:latest"

# Parse the options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --nodeID)            nodeID=${value} ;;
            --rpcServerAddress)  rpcServerAddress=${value} ;;
            --address)           address=${value} ;;
            --logPath)           logPath=${value} ;;
            --port)              port=${value} ;;
            --image)             image=${value} ;;
            *)
    esac
done

# Stop and remove any running container, and then star
echo "Force stopping any existing container ... "
docker rm -f  local-testnet-obscuroscan 2>/dev/null

echo "Starting the Obscuroscan..."
docker run --env PORT=${port} -p 8080:${port} --name=local-testnet-obscuroscan \
    --detach \
    --network=node_network \
    --entrypoint ./tools/obscuroscan/main/main \
     ${image} \
    --nodeID=${nodeID} \
    --rpcServerAddress=${rpcServerAddress} \
    --address=${address} \
    --logPath=${logPath}