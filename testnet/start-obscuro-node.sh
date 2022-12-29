#!/usr/bin/env bash

#
# This script downloads and builds the Obscuro node.
#
# Note: Be aware that a network MUST always have EXACTLY ONE genesis node (i.e. with flag `is_genesis=true`);
# otherwise, your node wil spin forever waiting for the network secret. In addition, a network MUST always have ONE
# sequencer (i.e. with flag `node_type=sequencer`); otherwise no rollups will be produced. In addition, the genesis
# node MUST be the sequencer.
#

help_and_exit() {
    echo ""
    echo "Usage: "
    echo "   ex: (run locally to internal l1 on local SGX NON capable hardware)"
    echo "      -  $(basename "${0}") --sgx_enabled=false --l1host=gethnetwork --mgmtcontractaddr=0xeDa66Cc53bd2f26896f6Ba6b736B1Ca325DE04eF --hocerc20addr=0xC0370e0b5C1A41D447BDdA655079A1B977C71aA9 --pocerc20addr=0x51D43a3Ca257584E770B6188232b199E76B022A2 --is_genesis=true --node_type=sequencer"
    echo ""
    echo "   ex: (run connected to an external l1 on local SGX capable hardware)"
    echo "      -  $(basename "${0}") --sgx_enabled=true --l1host=testnet-gethnetwork-18.uksouth.azurecontainer.io --mgmtcontractaddr=0xeDa66Cc53bd2f26896f6Ba6b736B1Ca325DE04eF --hocerc20addr=0xC0370e0b5C1A41D447BDdA655079A1B977C71aA9 --pocerc20addr=0x51D43a3Ca257584E770B6188232b199E76B022A2 --node_type=sequencer"
    echo ""
    echo "  l1host             *Required* Set the l1 host address"
    echo ""
    echo "  mgmtcontractaddr   *Required* Set the management contract address"
    echo ""
    echo "  hocerc20addr       *Required* Set the erc20 contract address for HOC"
    echo ""
    echo "  pocerc20addr       *Required* Set the erc20 contract address for POC"
    echo ""
    echo "  sgx_enabled        *Required* Set the execution to run with sgx enabled"
    echo ""
    echo "  sequencerID        *Optional* Set the sequencer address. Defaults to 0x0654D8B60033144D567f25bF41baC1FB0D60F23B"
    echo ""
    echo "  host_id            *Optional* Set the host ID used by the enclave. Defaults to 0x0654D8B60033144D567f25bF41baC1FB0D60F23B"
    echo ""
    echo "  l1port             *Optional* Set the l1 port. Defaults to 9000"
    echo ""
    echo "  pkstring           *Optional* Set the pk string. Defaults to 8ead642ca80dadb0f346a66cd6aa13e08a8ac7b5c6f7578d4bac96f5db01ac99"
    echo ""
    echo "  is_genesis         *Optional* Set the node as genesis node. Defaults to false"
    echo ""
    echo "  node_type          *Optional* Set the node's type. Defaults to validator"
    echo ""
    echo "  log_level          *Optional* Sets the log level. Defaults to 2 (warn)."
    echo ""
    echo "  p2p_public_address *Optional* Set host p2p public address. Defaults to 127.0.0.1:10000"
    echo ""
    echo "  pccs_addr           *Optional* Set the enclave Provision Certificate Cache Service. Defaults to 127.0.0.1:8081"
    echo ""
    echo "  profiler_enabled   *Optional* Enables the profiler in the host + enclave. Defaults to false"
    echo ""
    echo "  debug_enclave      *Optional* Dev mode, with a dlv debugger remote attach on port 2345"
    echo ""
    echo "  dev_testnet        *Optional* Uses dev images for dev testnet"
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
node_type=validator
profiler_enabled=false
p2p_public_address="127.0.0.1:10000"
debug_enclave=false
dev_testnet=false
host_id=0x0654D8B60033144D567f25bF41baC1FB0D60F23B
pk_string=8ead642ca80dadb0f346a66cd6aa13e08a8ac7b5c6f7578d4bac96f5db01ac99
log_level=4
sequencer_id=0x0654D8B60033144D567f25bF41baC1FB0D60F23B
pccs_addr="127.0.0.1:8081"


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
            --hocerc20addr)             hoc_erc20_addr=${value} ;;
            --pocerc20addr)             poc_erc20_addr=${value} ;;
            --pkstring)                 pk_string=${value} ;;
            --sgx_enabled)              sgx_enabled=${value} ;;
            --is_genesis)               is_genesis=${value} ;;
            --node_type)                node_type=${value} ;;
            --log_level)                log_level=${value} ;;
            --profiler_enabled)         profiler_enabled=${value} ;;
            --p2p_public_address)       p2p_public_address=${value} ;;
            --debug_enclave)            debug_enclave=${value} ;;
            --dev_testnet)              dev_testnet=${value} ;;
            --sequencerID)              sequencer_id=${value} ;;
            --pccs_addr)                pccs_addr=${value} ;;

            --help)                     help_and_exit ;;
            *)
    esac
done

if [[ -z ${l1_host:-} || -z ${mgmt_contract_addr:-} || -z ${hoc_erc20_addr:-} || -z ${poc_erc20_addr:-} || -z ${sgx_enabled:-} ]];
then
    help_and_exit
fi


# reset any data in the env file
echo "" > "${testnet_path}/.env"

# set the data in the env file
echo "PKSTRING=${pk_string}" >> "${testnet_path}/.env"
echo "HOSTID=${host_id}"  >> "${testnet_path}/.env"
echo "MGMTCONTRACTADDR=${mgmt_contract_addr}"  >> "${testnet_path}/.env"
echo "HOCERC20ADDR=${hoc_erc20_addr}"  >> "${testnet_path}/.env"
echo "POCERC20ADDR=${poc_erc20_addr}"  >> "${testnet_path}/.env"
echo "L1HOST=${l1_host}" >> "${testnet_path}/.env"
echo "L1PORT=${l1_port}" >> "${testnet_path}/.env"
echo "ISGENESIS=${is_genesis}" >> "${testnet_path}/.env"
echo "NODETYPE=${node_type}" >> "${testnet_path}/.env"
echo "LOGLEVEL=${log_level}" >> "${testnet_path}/.env"
echo "PROFILERENABLED=${profiler_enabled}" >> "${testnet_path}/.env"
echo "P2PPUBLICADDRESS=${p2p_public_address}" >> "${testnet_path}/.env"
echo "SEQUENCERID=${sequencer_id}" >> "${testnet_path}/.env"
echo "PCCS_ADDR=${pccs_addr}" >> "${testnet_path}/.env"


if ${debug_enclave} ;
then
  echo "Starting DEBUG enclave and host..."
  docker compose -f docker-compose.debug.yml up enclave host -d
  exit 0
fi

if ${dev_testnet} ;
then
  echo "Starting enclave and host with dev testnet images..."
  docker compose -f docker-compose.dev-testnet.yml up enclave host edgelessdb -d
  exit 0
fi

if ${sgx_enabled} ;
then
  echo "Starting enclave with enabled SGX and host..."
  docker compose up enclave host edgelessdb -d
  exit 0
fi

echo "Starting enclave with DISABLED SGX and host..."
docker compose -f docker-compose.non-sgx.yml up enclave host -d

echo "Waiting 20s for the node to be up..."
sleep 20
echo "Node should be up and running"