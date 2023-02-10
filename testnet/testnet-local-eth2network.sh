#!/usr/bin/env bash

#
# This script starts the eth2network in docker
#

#
# Prefunded PK used for contract deployment
# f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb -> 0x13E23Ca74DE0206C56ebaE8D51b5622EFF1E9944
#
# Prefunded PK used by the obscuro nodes for rollups issuance
# 8ead642ca80dadb0f346a66cd6aa13e08a8ac7b5c6f7578d4bac96f5db01ac99 -> 0x0654D8B60033144D567f25bF41baC1FB0D60F23B

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}") --pkaddresses=0x13E23Ca74DE0206C56ebaE8D51b5622EFF1E9944,0x0654D8B60033144D567f25bF41baC1FB0D60F23B"
    echo ""
    echo "  pkaddresses          *Required* Set the prefunded addresses"
    echo ""
    echo "  wsport             *Optional* Set web socket start port. Defaults to 9000"
    echo ""
    echo "  port               *Optional* Set geth http port. Defaults to 8025"
    echo ""
    echo ""
    echo ""
    exit 1  # Exit with error explicitly
}
# Ensure any fail is loud and explicit
set -euo pipefail

# Set default options
wsport=9000
port=8025

# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --pkaddresses)                 pkaddresses=${value} ;;
            --wsport)                   wsport=${value} ;;
            --port)                     port=${value} ;;
            --help)                     help_and_exit ;;
            *)
    esac
done

if [[ -z ${pkaddresses:-} ]];
then
    help_and_exit
fi

# start the geth network
echo "Starting the eth2network.."
docker network create --driver bridge node_network || true
docker run --name=eth2network -d \
  --network=node_network \
  -p 8025:8025 -p 8026:8026 -p 9000:9000 -p 9001:9001 \
  --entrypoint /home/obscuro/go-obscuro/integration/eth2network/main/main \
   testnetobscuronet.azurecr.io/obscuronet/eth2network:latest \
  --numNodes=1 \
  --gethHTTPStartPort=${port} \
  --gethWSStartPort=${wsport} \
  --prefundedAddrs=${pkaddresses}

echo "Network should be up and running"