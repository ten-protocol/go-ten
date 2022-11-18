#!/usr/bin/env bash

#
# This script starts up the obscuroscan server
#

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}")"
    echo ""
    echo "  rpcServerAddress        *Optional* Set the rpc server address (defaults to http://testnet.obscu.ro:13000)"
    echo ""
    echo "  address                 *Optional* Set the obscuroscan endpoint address (defaults to 0.0.0.0:80)"
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
rpcServerAddress='http://testnet-host-1:13000'
address='0.0.0.0:80'
docker_image="testnetobscuronet.azurecr.io/obscuronet/obscuroscan:latest"

# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --rpcServerAddress)         rpcServerAddress=${value} ;;
            --address)                  address=${value} ;;
            --help)                     help_and_exit ;;
            *)
    esac
done

# ensure required fields
if [[ -z ${rpcServerAddress:-} || -z ${address:-} ]];
then
    help_and_exit
fi

# start the container
echo "Starting the obscuroscan server..."
docker network create --driver bridge node_network || true
docker run --name=obscuroscan \
    --network=node_network \
    -p 80:80 \
    --entrypoint /home/go-obscuro/tools/obscuroscan/main/main \
    "${docker_image}" \
    --rpcServerAddress=${rpcServerAddress} \
    --address=${address}
echo ""
