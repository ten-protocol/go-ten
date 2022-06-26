#!/usr/bin/env bash

#
# This script downloads and builds the obscuro node
#

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}") --host_id=0x0000000000000000000000000000000000000001 --l1host=127.0.0.1"
    echo ""
    echo "  host_id            *Required* Set the node ID"
    echo ""
    echo "  l1host             *Required* Set the l1 host address"
    echo ""
    echo "  l1port             *Optional* Set the l1 port. Defaults to 9000"
    echo ""
    echo "  local_gethnetwork  *Optional* Start a local geth network. Defaults to false"
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
l1port=9000
local_gethnetwork=false

# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --l1host)                   l1host=${value} ;;
            --l1port)                   l1port=${value} ;;
            --host_id)                  host_id=${value} ;;
            --local_gethnetwork)        local_gethnetwork=${value} ;;
            --help)                     help_and_exit ;;
            *)
    esac
done

if [[ -z ${l1host:-} || -z ${host_id:-} ]];
then
    help_and_exit
fi

# Hardcoded geth pre funded private key
pk_address=0x13E23Ca74DE0206C56ebaE8D51b5622EFF1E9944
pk_string=f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb

# set the pk in the env file
echo "PKSTRING=${pk_string}" > "${testnet_path}/.env"
echo "PKADDR=${pk_address}" >> "${testnet_path}/.env"
echo "HOSTID=${host_id}"  >> "${testnet_path}/.env"

# start the geth network
if ${local_gethnetwork}
then
  echo "Starting the gethnetwork.."
  docker network create --driver bridge node_network || true
  docker run --name=gethnetwork -d \
    --network=node_network \
    --entrypoint /home/go-obscuro/integration/gethnetwork/main/main \
     obscuro_gethnetwork:latest \
    --numNodes=3 \
    --startPort=8000 \
    --websocketStartPort=${l1port} \
    --prefundedAddrs=${pk_address}

  echo "Waiting 30s for the network to be up..."
  sleep 30
fi

# set the addresses in the env file
echo "L1HOST=${l1host}" >> "${testnet_path}/.env"
echo "L1PORT=${l1port}" >> "${testnet_path}/.env"

# deploy contracts to the geth network
echo "Deploying contracts to the geth network..."
docker network create --driver bridge node_network || true
docker run --name=contractdeployer \
    --network=node_network \
    --entrypoint /home/go-obscuro/tools/contractdeployer/main/main \
     obscuro_contractdeployer:latest \
    -l1NodeHost=${l1host} \
    --l1NodePort=${l1port} \
    --privateKey=${pk_string}

# storing the contract addresses to the .env file
log_output=$(docker logs --tail 1 contractdeployer)
json_output=$(echo ${log_output} | awk -F"[{}]" '{print "{"$2"}"}')
mgmtContractAddr=$(echo "${json_output}"  | jq .MgmtContractAddr)
erc20ContractAddr=$(echo "${json_output}" | jq .ERC20ContractAddr)

echo "MGMTCONTRACTADDR=${mgmtContractAddr}" >> "${testnet_path}/.env"
echo "ERC20CONTRACTADDR=${erc20ContractAddr}" >> "${testnet_path}/.env"


echo "Starting enclave and host..."
docker compose up enclave host


