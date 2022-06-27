#!/usr/bin/env bash

#
# This script starts the gethnetwork in docker
#

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}") --pkaddress=0x13E23Ca74DE0206C56ebaE8D51b5622EFF1E9944"
    echo ""
    echo "  pkaddress          *Required* Set the prefunded address"
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
            --pkaddress)                 pkaddress=${value} ;;
            --wsport)                   wsport=${value} ;;
            --port)                     port=${value} ;;
            --help)                     help_and_exit ;;
            *)
    esac
done

if [[ -z ${pkaddress:-} ]];
then
    help_and_exit
fi

# start the geth network
echo "Starting the gethnetwork.."
docker network create --driver bridge node_network || true
docker run --name=gethnetwork -d \
  --network=node_network \
  --entrypoint /home/go-obscuro/integration/gethnetwork/main/main \
   testnetobscuronet.azurecr.io/obscuronet/obscuro_gethnetwork:latest \
  --numNodes=3 \
  --startPort=${port} \
  --websocketStartPort=${wsport} \
  --prefundedAddrs=${pkaddress}

echo "Waiting 30s for the network to be up..."
sleep 30
echo "Network should be up and running"