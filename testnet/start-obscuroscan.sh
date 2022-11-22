#!/usr/bin/env bash

#
# This script starts up the obscuroscan server
#

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}") --rpcServerAddress=http://testnet-host-1:13000 --receivingPort=80"
    echo ""
    echo "  rpcServerAddress        *Optional* Set the rpc server address (defaults to http://testnet-host-1:13000)"
    echo ""
    echo "  receivingPort           *Optional* Set the ObscuroScan server receiving port (defaults to 80)"
    echo ""
    echo ""
    echo ""
    exit 1  # Exit with error explicitly
}
# Ensure any fail is loud and explicit
set -euo pipefail

# Define local usage vars
start_path="$(cd "$(dirname "${0}")" && pwd)"
testnet_path="${start_path}"

# Define defaults
receivingPort=80
rpcServerAddress='http://testnet-host-1:13000'
docker_image="testnetobscuronet.azurecr.io/obscuronet/obscuroscan:latest"

# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --rpcServerAddress)         rpcServerAddress=${value} ;;
            --receivingPort)            receivingPort=${value} ;;
            --help)                     help_and_exit ;;
            *)
    esac
done

# ensure required fields
if [[ -z ${rpcServerAddress:-} || -z ${receivingPort:-} ]];
then
    help_and_exit
fi

# start the container
echo "Starting the obscuroscan server..."
docker network create --driver bridge node_network || true
docker run --name=obscuroscan \
    --detach \
    --network=node_network \
    -p $receivingPort:$receivingPort \
    --entrypoint /home/go-obscuro/tools/obscuroscan/main/main \
    "${docker_image}" \
    --address='0.0.0.0:'$receivingPort \
    --rpcServerAddress=${rpcServerAddress}
echo ""
