#!/usr/bin/env bash

#
# This script deploys contracts to testnet
#

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}") --l2host=testnet-host-1"
    echo ""
    echo "  l2host             *Required* Set the l2 host address"
    echo ""
    echo "  hocpkstring        *Optional* Set the pkstring to deploy HOC contract"
    echo ""
    echo "  pocpkstring        *Optional* Set the pkstring to deploy POC contract"
    echo ""
    echo "  l2port             *Optional* Set the l2 port. Defaults to 10000"
    echo ""
    echo "  docker_image       *Optional* Sets the docker image to use. Defaults to testnetobscuronet.azurecr.io/obscuronet/obscuro_contractdeployer:latest"
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
l2port=13001
# todo: get rid of these defaults and require them to be passed in, using github secrets for testnet values (requires bridge.go changes)
hocpkstring="6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682"
pocpkstring="4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8"
hocerc20address="0xf3a8bd422097bFdd9B3519Eaeb533393a1c561aC"
docker_image="testnetobscuronet.azurecr.io/obscuronet/obscuro_contractdeployer:latest"

# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --l2host)                   l2host=${value} ;;
            --l2port)                   l2port=${value} ;;
            --hocpkstring)              hocpkstring=${value} ;;
            --pocpkstring)              pocpkstring=${value} ;;
            --docker_image)             docker_image=${value} ;;
            --help)                     help_and_exit ;;
            *)
    esac
done

# ensure required fields
if [[ -z ${l2host:-} || -z ${hocpkstring:-} || -z ${pocpkstring:-}  ]];
then
    help_and_exit
fi

# deploy contracts to the obscuro network
echo "Deploying HOC ERC20 contract to the obscuro network..."
docker network create --driver bridge node_network || true
docker run --name=hocL2deployer \
    --network=node_network \
    --entrypoint /home/go-obscuro/tools/contractdeployer/main/main \
     "${docker_image}" \
    --nodeHost=${l2host} \
    --nodePort=${l2port} \
    --contractName="L2ERC20" \
    --privateKey=${hocpkstring}\
    --constructorParams="Hocus,HOC,1000000000000000000000000000000"
echo ""

echo "Deploying POC ERC20 contract to the obscuro network..."
docker run --name=pocL2deployer \
    --network=node_network \
    --entrypoint /home/go-obscuro/tools/contractdeployer/main/main \
     "${docker_image}" \
    --nodeHost=${l2host} \
    --nodePort=${l2port} \
    --contractName="L2ERC20" \
    --privateKey=${pocpkstring}\
    --constructorParams="Pocus,POC,1000000000000000000000000000000"
echo ""

echo "Deploying Guessing game contract to the obscuro network..."
docker run --name=guessingGameL2deployer \
    --network=node_network \
    --entrypoint /home/go-obscuro/tools/contractdeployer/main/main \
     "${docker_image}" \
    --nodeHost=${l2host} \
    --nodePort=${l2port} \
    --contractName="GUESS" \
    --privateKey=${pocpkstring}\
    --constructorParams="100,${hocerc20address}"
echo ""
