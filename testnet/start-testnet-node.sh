#!/usr/bin/env bash

#
# This script downloads and builds the obscuro node
#

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}") "
    echo ""
    echo ""
    echo ""
    exit 1  # Exit with error explicitly
}

# Ensure any fail is loud and explicit
set -euo pipefail

# Define local usage vars
start_path="$(cd "$(dirname "${0}")" && pwd)"

#sudo apt-get update
#sudo apt-get install -y docker.io docker-compose

PKADDR=0x13E23Ca74DE0206C56ebaE8D51b5622EFF1E9944
PKSTRING=f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb

echo "PKSTRING=${PKSTRING}" > .env
echo "PKADDR=${PKADDR}" >> .env

echo "Starting the gethnetwork.."
docker network create --driver bridge node_network || true
docker run --name=gethnetwork -d \
  --network=node_network \
  --entrypoint /home/go-obscuro/integration/gethnetwork/main/main obscuro_gethnetwork:latest \
  --numNodes=3 \
  --startPort=8000 \
  --websocketStartPort=9000 \
  --prefundedAddrs=${PKADDR}

#docker compose up -d gethnetwork
echo "Waiting 30s for the network to be up..."
sleep 10
docker compose up contractdeployer
log_output=$(docker-compose logs --no-color --no-log-prefix --tail 1 contractdeployer)
json_output=$(echo ${log_output} | awk -F"[{}]" '{print "{"$2"}"}')
mgmtContractAddr=$(echo "${json_output}"  | jq .MgmtContractAddr)
erc20ContractAddr=$(echo "${json_output}" | jq .ERC20ContractAddr)

echo "MGMTCONTRACTADDR=${mgmtContractAddr}" >> .env
echo "ERC20CONTRACTADDR=${erc20ContractAddr}" >> .env

echo "Starting enclave and host..."
docker compose up enclave host


