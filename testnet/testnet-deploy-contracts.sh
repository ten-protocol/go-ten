#!/usr/bin/env bash

#
# This script deploys contracts to the L1 of testnet
#

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}") --l1host=gethnetwork --pkstring=f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb"
    echo ""
    echo "  l1host             *Required* Set the l1 host address"
    echo ""
    echo "  pkstring           *Required* Set the pkstring to deploy contracts"
    echo ""
    echo "  l1port             *Optional* Set the l1 port. Defaults to 9000"
    echo ""
    echo "  docker_image       *Optional* Sets the docker image to use. Defaults to testnetobscuronet.azurecr.io/obscuronet/contractdeployer:latest"
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
l1port=8025
docker_image="testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest"
pkstring="f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb"
    
# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --l1host)                   l1host=${value} ;;
            --l1port)                   l1port=${value} ;;
            --pkstring)                 pkstring=${value} ;;
            --docker_image)             docker_image=${value} ;;
            --help)                     help_and_exit ;;
            *)
    esac
done

# ensure required fields
if [[ -z ${l1host:-} || -z ${pkstring:-}  ]];
then
    help_and_exit
fi

network_cfg='{ 
        "layer1" : {
            "url" : '"\"http://${l1host}:${l1port}\""',
            "live" : false,
            "saveDeployments" : true,
            "deploy": [ "deployment_scripts/layer1", "deployment_scripts/testnet" ],
            "accounts": [ "'${pkstring}'" ]
        }
    }'

# deploy contracts to the geth network
echo "Creating docker network..."
docker network create --driver bridge node_network || true

echo "Deploying contracts to Layer 1 using obscuro hardhat container..."
docker run --name=hh-l1-deployer \
    --network=node_network \
    -e NETWORK_JSON="${network_cfg}" \
    "${docker_image}" \
    deploy \
    --network layer1

# --tail 5 gets the last 5 lines of the deployment; grep -e '' gives us the line matching the pattern; cut takes out the address where the contract has been deployed
# The standard output from the hh deploy plugin looks like
#  deploying "ManagementContract" (tx: 0xcb6e341c9f30e1b86214542bcd1c930f202201b4483801df5cd3c1f53c4b55f8)...: deployed at 0xeDa66Cc53bd2f26896f6Ba6b736B1Ca325DE04eF with 2533700 gas
mgmtContractAddr=$(docker logs --tail 5 hh-l1-deployer | grep -e 'ManagementContract' | cut -c 121-162)

echo "MGMTCONTRACTADDR=${mgmtContractAddr}" > "${testnet_path}/.env"