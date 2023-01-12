#!/usr/bin/env bash

#
# This script removes any existing backend pool element from the azure load balancer for testnet
#
#

if [[ $1 == "testnet" ]]; then
  lb=testnet-loadbalancer
  pool=Backend-Pool-Obscuro-testnet
else
  lb=devtestnet-loadbalancer
  pool=Backend-Pool-Obscuro-devtestnet
fi

nic_id=$(az network lb address-pool show \
    --resource-group Testnet \
    --lb-name ${lb} \
    --name ${pool} \
    --query backendIpConfigurations \
    --output tsv | cut -f5 | cut -d "/" -f 9)

ipconfig_id=$(az network lb address-pool show \
    --resource-group Testnet \
    --lb-name ${lb} \
    --name ${pool} \
    --query backendIpConfigurations \
    --output tsv | cut -f5 | cut -d "/" -f 11)

if [[ -z "${nic_id}" ]]; then
  echo "No Nic found in the load balancer"
  exit 0
fi

if [[ -z "${ipconfig_id}" ]]; then
    echo "No Ip config found in the load balancer"
    exit 0
fi

az network nic ip-config address-pool remove \
   --address-pool ${pool} \
   --ip-config-name "${ipconfig_id}" \
   --nic-name "${nic_id}" \
   --resource-group Testnet \
   --lb-name ${lb} \
   

echo 'Load balancer removed successfully'
