#!/usr/bin/env bash

#
# This script downloads and builds the obscuro node
#

help_and_exit() {
    echo ""
    echo "Usage: "
    echo "   ex: (run locally)"
    echo "      -  $(basename "${0}") --host_id=0x0000000000000000000000000000000000000001 --l1host=gethnetwork --local_gethnetwork=true --mgmtcontractaddr=0x7e440D3F8a82636529b0A4Fb9a4Ff66f8Bc7141F --erc20contractaddr=0xF63035376a11007DDEBed404405b69F079b17836"
    echo ""
    echo "   ex: (run connect external)"
    echo "      -  $(basename "${0}") --host_id=0x0000000000000000000000000000000000000001 --l1host=testnet-gethnetwork-18.uksouth.azurecontainer.io --mgmtcontractaddr=0x7e440D3F8a82636529b0A4Fb9a4Ff66f8Bc7141F --erc20contractaddr=0xF63035376a11007DDEBed404405b69F079b17836"
    echo ""
    echo "  host_id            *Required* Set the node ID"
    echo ""
    echo "  l1host             *Required* Set the l1 host address"
    echo ""
    echo "  mgmtcontractaddr   *Required* Set the management contract address"
    echo ""
    echo "  erc20contractaddr  *Required* Set the erc20 contract address"
    echo ""
    echo "  l1port             *Optional* Set the l1 port. Defaults to 9000"
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

# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --l1host)                   l1host=${value} ;;
            --l1port)                   l1port=${value} ;;
            --host_id)                  host_id=${value} ;;
            --mgmtcontractaddr)         mgmtcontractaddr=${value} ;;
            --erc20contractaddr)        erc20contractaddr=${value} ;;
            --help)                     help_and_exit ;;
            *)
    esac
done

if [[ -z ${l1host:-} || -z ${host_id:-} || -z ${mgmtcontractaddr:-} || -z ${erc20contractaddr:-} ]];
then
    help_and_exit
fi

# Hardcoded geth pre funded private key
pk_address=0x13E23Ca74DE0206C56ebaE8D51b5622EFF1E9944
pk_string=f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb

# set the data in the env file
echo "PKSTRING=${pk_string}" > "${testnet_path}/.env"
echo "PKADDR=${pk_address}" >> "${testnet_path}/.env"
echo "HOSTID=${host_id}"  >> "${testnet_path}/.env"
echo "MGMTCONTRACTADDR=${mgmtcontractaddr}"  >> "${testnet_path}/.env"
echo "ERC20CONTRACTADDR=${erc20contractaddr}"  >> "${testnet_path}/.env"
echo "L1HOST=${l1host}" >> "${testnet_path}/.env"
echo "L1PORT=${l1port}" >> "${testnet_path}/.env"


echo "Starting enclave and host..."
docker compose up enclave host -d


