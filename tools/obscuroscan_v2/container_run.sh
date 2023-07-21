#!/usr/bin/env bash

#
# This script starts up the obscuroscan container to connect to a local testnet
#

# Define defaults
nodeHostAddress="http://testnet.obscu.ro:13001"
image="testnetobscuronet.azurecr.io/obscuronet/obscuro_scan_testnet:latest"

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}")"
    echo ""
    echo "  nodeHostAddress          *Optional* Set the l2 host address ( ${nodeHostAddress}"
    echo ""
    echo "  image                    *Optional* Set image to use ( ${image} )"
    echo ""
    echo ""
    echo ""
    exit 1  # Exit with error explicitly
}

# Ensure any fail is loud and explicit
set -euo pipefail

# Parse the options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --nodeHostAddress)        nodeHostAddress=${value} ;;
            --image)           image=${value} ;;
            *)
    esac
done

echo "Starting the obscuroscan server..."
docker run -p 81:80 -p 82:8080 --name=local-obscuro-scan \
    --detach \
    --network=node_network \
    --env NODEHOSTADDRESS=${nodeHostAddress} \
    --entrypoint=./entrypoint.sh \
     ${image}
