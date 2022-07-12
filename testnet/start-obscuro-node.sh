#!/usr/bin/env bash

#
# This script downloads and builds the obscuro node
#

help_and_exit() {
    echo ""
    echo "Usage: "
    echo "   ex: (run locally)"
    echo "      -  $(basename "${0}") --sgx_enabled=false --host_id=0x0000000000000000000000000000000000000001 --l1host=gethnetwork --mgmtcontractaddr=0xeDa66Cc53bd2f26896f6Ba6b736B1Ca325DE04eF --erc20contractaddr=0xC0370e0b5C1A41D447BDdA655079A1B977C71aA9"
    echo ""
    echo "   ex: (run connect external)"
    echo "      -  $(basename "${0}") --sgx_enabled=true --host_id=0x0000000000000000000000000000000000000001 --l1host=testnet-gethnetwork-18.uksouth.azurecontainer.io --mgmtcontractaddr=0x7e440D3F8a82636529b0A4Fb9a4Ff66f8Bc7141F --erc20contractaddr=0xF63035376a11007DDEBed404405b69F079b17836"
    echo ""
    echo "  host_id            *Required* Set the node ID"
    echo ""
    echo "  l1host             *Required* Set the l1 host address"
    echo ""
    echo "  mgmtcontractaddr   *Required* Set the management contract address"
    echo ""
    echo "  erc20contractaddr  *Required* Set the erc20 contract address"
    echo ""
    echo "  sgx_enabled        *Required* Set the execution to run with sgx enabled"
    echo ""
    echo "  l1port             *Optional* Set the l1 port. Defaults to 9000"
    echo ""
    echo "  pkaddress          *Optional* Set the pk address. Defaults to 0x13E23Ca74DE0206C56ebaE8D51b5622EFF1E9944"
    echo ""
    echo "  pkstring           *Optional* Set the pk string. Defaults to f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb"
    echo ""
    echo "  is_genesis         *Optional* Set the node as genesis node. Defaults to false"
    echo ""
    echo "  profiler_enabled   *Optional* Enables the profiler in the host + enclave. Defaults to false"
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
l1_port=9000
is_genesis=false
profiler_enabled=false
pk_address=0x13E23Ca74DE0206C56ebaE8D51b5622EFF1E9944
pk_string=f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb


# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --l1host)                   l1_host=${value} ;;
            --l1port)                   l1_port=${value} ;;
            --host_id)                  host_id=${value} ;;
            --mgmtcontractaddr)         mgmt_contract_addr=${value} ;;
            --erc20contractaddr)        erc20_contract_addr=${value} ;;
            --pkaddress)                pk_address=${value} ;;
            --pkstring)                 pk_string=${value} ;;
            --sgx_enabled)              sgx_enabled=${value} ;;
            --is_genesis)               is_genesis=${value} ;;
            --profiler_enabled)         profiler_enabled=${value} ;;
            --help)                     help_and_exit ;;
            *)
    esac
done

if [[ -z ${l1_host:-} || -z ${host_id:-} || -z ${mgmt_contract_addr:-} || -z ${erc20_contract_addr:-} || -z ${sgx_enabled:-} ]];
then
    help_and_exit
fi


# set the data in the env file
echo "PKSTRING=${pk_string}" > "${testnet_path}/.env"
echo "PKADDR=${pk_address}" >> "${testnet_path}/.env"
echo "HOSTID=${host_id}"  >> "${testnet_path}/.env"
echo "MGMTCONTRACTADDR=${mgmt_contract_addr}"  >> "${testnet_path}/.env"
echo "ERC20CONTRACTADDR=${erc20_contract_addr}"  >> "${testnet_path}/.env"
echo "L1HOST=${l1_host}" >> "${testnet_path}/.env"
echo "L1PORT=${l1_port}" >> "${testnet_path}/.env"
echo "ISGENESIS=${is_genesis}" >> "${testnet_path}/.env"
echo "PROFILERENABLED=${profiler_enabled}" >> "${testnet_path}/.env"


if ${sgx_enabled} ;
then
  echo "Starting enclave with enabled SGX and host..."
  docker compose up enclave host edgelessdb -d
else
  echo "Starting enclave with DISABLED SGX and host..."
  docker compose -f docker-compose.non-sgx.yml up enclave host -d
fi

