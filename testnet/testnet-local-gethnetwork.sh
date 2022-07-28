#!/usr/bin/env bash

#
# This script starts the gethnetwork in docker
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
    echo "  port               *Optional* Set start port. Defaults to 8000"
    echo ""
    echo ""
    echo ""
    exit 1  # Exit with error explicitly
}
# Ensure any fail is loud and explicit
set -euo pipefail

# Set default options
wsport=9000
port=8000

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
echo "Starting the gethnetwork.."
docker network create --driver bridge node_network || true
docker run --name=gethnetwork -d \
  --network=node_network \
  -p 8025:8025 -p 8026:8026 -p 8027:8027 -p 9000:9000 -p 9001:9001 -p 9002:9002 \
  --entrypoint /home/go-obscuro/integration/gethnetwork/main/main \
   testnetobscuronet.azurecr.io/obscuronet/obscuro_gethnetwork:latest \
  --numNodes=3 \
  --startPort=${port} \
  --websocketStartPort=${wsport} \
  --prefundedAddrs=${pkaddresses}

echo "Waiting 30s for the network to be up..."
sleep 30
echo "Network should be up and running"