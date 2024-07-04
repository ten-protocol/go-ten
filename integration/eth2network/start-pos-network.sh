#!/bin/bash

./geth --datadir=gethdata account import secret.txt
echo "Private key imported"

./geth --datadir=gethdata init genesis.json
echo "Geth genesis initialized"

# Set the log directory to the current directory
prysm_logs="$(pwd)/logs"

# Run the Prysm beacon node
./beacon-chain --datadir beacondata \
               --min-sync-peers 0 \
               --genesis-state genesis.ssz \
               --bootstrap-node= \
               --interop-eth1data-votes \
               --chain-config-file config.yml \
               --contract-deployment-block 0 \
               --chain-id 32382 \
               --rpc-host=127.0.0.1 \
               --rpc-port=4000 \
               --accept-terms-of-use \
               --jwt-secret jwt.hex \
               --suggested-fee-recipient 0x123463a4B065722E99115D6c222f267d9cABb524 \
               --minimum-peers-per-subnet 0 \
               --enable-debug-rpc-endpoints \
               --verbosity=debug \
               --execution-endpoint gethdata/geth.ipc > "${prysm_logs}/beacon-chain.log" 2>&1 &

echo "Beacon node started"

## Allow time for the beacon node to start
#sleep 30
#
## Check if beacon node started successfully
#if ! pgrep -f beacon-chain > /dev/null; then
#    echo "Failed to start beacon node"
#    exit 1
#fi

# Run Prysm validator client
./validator --beacon-rpc-provider=127.0.0.1:4000 \
            --datadir validatordata \
            --accept-terms-of-use \
            --interop-num-validators 4 \
            --chain-config-file config.yml > "${prysm_logs}/validator.log" 2>&1 &

echo "Validator client started"

## Allow time for the validator client to start
#sleep 30
#
## Check if validator client started successfully
#if ! pgrep -f validator > /dev/null; then
#    echo "Failed to start validator client"
#    exit 1
#fi

# Run go-ethereum
./geth --http \
       --http.api eth,net,web3 \
       --ws --ws.api eth,net,web3 \
       --authrpc.jwtsecret jwt.hex \
       --datadir gethdata \
       --nodiscover \
       --syncmode full \
       --allow-insecure-unlock \
       --unlock 0x123463a4b065722e99115d6c222f267d9cabb524 \
       --password ./password.txt > "${prysm_logs}/geth.log" 2>&1 &

echo "Geth network started"

## Allow time for geth to start
#sleep 30
#
## Check if geth started successfully
#if ! pgrep -f geth > /dev/null; then
#    echo "Failed to start geth"
#    exit 1
#fi

echo "Running ..."