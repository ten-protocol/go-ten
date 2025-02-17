#!/usr/bin/env bash

#
# This script starts up the faucet server to connect to a local testnet
#

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}")"
    echo ""
    echo "  nodeHost           *Optional* Set the l2 host address"
    echo ""
    echo "  nodePort           *Optional* Set the l2 host port"
    echo ""
    echo "  port               *Optional* Set the faucet server port"
    echo ""
    echo "  pk                 *Optional* Set the pre-funded private key"
    echo ""
    echo "  jwtSecret          *Optional* Set the jwt secret"
    echo ""
    echo "  image              *Optional* Set image to use, defaults to testnetobscuronet.azurecr.io/obscuronet/faucet_sepolia_testnet:latest"
    echo ""
    echo ""
    echo ""
    exit 1  # Exit with error explicitly
}

# Ensure any fail is loud and explicit
set -euo pipefail

# Define defaults
nodeHost="validator-host"
nodePort=13010
port=80
pk="0x8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b"
jwtSecret="This_is_the_secret"
image="testnetobscuronet.azurecr.io/obscuronet/faucet_sepolia_testnet:latest"

# Parse the options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --nodeHost)        nodeHost=${value} ;;
            --nodePort)        nodePort=${value} ;;
            --port)            port=${value} ;;
            --pk)              pk=${value} ;;
            --jwtSecret)       jwtSecret=${value} ;;
            --image)           image=${value} ;;
            *)
    esac
done

# Stop and remove any running container, and then star
echo "Force stopping any existing container ... "
docker rm -f  local-testnet-faucet 2>/dev/null

echo "Starting the faucet server..."
docker run --env PORT=${port} -p 8080:${port} --name=local-testnet-faucet \
    --detach \
    --network=node_network \
    --entrypoint ./faucet  \
     ${image} \
    --nodeHost=${nodeHost} \
    --nodePort=${nodePort} \
    --pk=${pk}\
    --jwtSecret=${jwtSecret}\
