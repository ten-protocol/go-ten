#!/usr/bin/env bash

# Default port values
BEACON_RPC_PORT=4000
GETH_RPC_PORT=8545
GETH_WS_PORT=8546

# Function to display usage
usage() {
  echo "Usage: $0 [--geth-rpc GETH_RPC_PORT] [--geth-ws GETH_WS_PORT] [--geth-http GETH_HTTP_PORT] [--beacon-rpc BEACON_RPC_PORT]"
  exit 1
}

# Parse command-line arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --beacon-rpc) BEACON_RPC_PORT="$2"; shift ;;
        --geth-rpc) GETH_RPC_PORT="$2"; shift ;;
        --geth-ws) GETH_WS_PORT="$2"; shift ;;
        *) usage ;;
    esac
    shift
done

./geth --datadir=gethdata account import pk.txt
echo "Private key imported"

./geth --datadir=gethdata init genesis.json
echo "Geth genesis initialized"

./prysmctl testnet generate-genesis \
           --fork deneb \
           --num-validators 2 \
           --genesis-time-delay 600 \
           --chain-config-file config.yml \
           --geth-genesis-json-in genesis.json \
           --geth-genesis-json-out genesis.json \
           --output-ssz genesis.ssz
sleep 5
echo "Prysm genesis generated"

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
               --rpc-port=${BEACON_RPC_PORT} \
               --accept-terms-of-use \
               --jwt-secret jwt.hex \
               --suggested-fee-recipient 0x123463a4B065722E99115D6c222f267d9cABb524 \
               --minimum-peers-per-subnet 0 \
               --enable-debug-rpc-endpoints \
               --verbosity=debug \
               --execution-endpoint gethdata/geth.ipc &
#               --execution-endpoint gethdata/geth.ipc > "${prysm_logs}/beacon-chain.log" 2>&1 &

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
./validator --beacon-rpc-provider=127.0.0.1:${BEACON_RPC_PORT} \
            --datadir validatordata \
            --accept-terms-of-use \
            --interop-num-validators 4 \
            --chain-config-file config.yml &
#            --chain-config-file config.yml > "${prysm_logs}/validator.log" 2>&1 &

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
       --http.port ${GETH_RPC_PORT} \
       --ws --ws.api eth,net,web3 \
       --ws.port ${GETH_WS_PORT} \
       --authrpc.jwtsecret jwt.hex \
       --datadir gethdata \
       --nodiscover \
       --syncmode full \
       --allow-insecure-unlock \
       --unlock 0x123463a4b065722e99115d6c222f267d9cabb524 \
       --password ./password.txt &
#       --password ./password.txt > "${prysm_logs}/geth.log" 2>&1 &

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