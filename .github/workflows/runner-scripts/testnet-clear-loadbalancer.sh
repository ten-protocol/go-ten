#!/usr/bin/env bash

#
# This script removes any existing backend pool element from the azure load balancer for testnet
#
#

# Ensure any fail is loud and explicit
set -euo pipefail

nic_id=$(az network lb address-pool show \
    --resource-group Testnet \
    --lb-name testnet-loadbalancer \
    --name Backend-Pool-Obscuro-Testnet \
    --query backendIpConfigurations \
    --output tsv | cut -f5 | cut -d "/" -f 9)

ipconfig_id=$(az network lb address-pool show \
    --resource-group Testnet \
    --lb-name testnet-loadbalancer \
    --name Backend-Pool-Obscuro-Testnet \
    --query backendIpConfigurations \
    --output tsv | cut -f5 | cut -d "/" -f 11)

az network nic ip-config address-pool remove \
   --address-pool Backend-Pool-Obscuro-Testnet \
   --ip-config-name "${ipconfig_id}" \
   --nic-name "${nic_id}" \
   --resource-group Testnet \
   --lb-name testnet-loadbalancer

echo 'Load balancer removed sucessfully'
