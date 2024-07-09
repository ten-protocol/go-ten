#!/usr/bin/env bash

# Default port values
BEACON_RPC_PORT=4000
GETH_HTTP_PORT=8545
GETH_WS_PORT=8546
BUILD_DIR="./build"
GETH_BINARY="./geth"
BEACON_BINARY="./beacon-chain"
PRYSMCTL_BINARY="./prysmctl"
VALIDATOR_BINARY="./validator"
BEACON_LOG_FILE="./beacon-chain.log"
VALIDATOR_LOG_FILE="./validator.log"
GETH_LOG_FILE="./geth.log"
GETHDATA_DIR="/gethdata"
BEACONDATA_DIR="/beacondata"
VALIDATORDATA_DIR="/validatordata"

# Function to display usage
usage() {
    echo "Usage: $0 [--geth-http GETH_HTTP_PORT] [--geth-ws GETH_WS_PORT] [--beacon-rpc BEACON_RPC_PORT] [--build-dir BUILD_DIR ]
    [--beacon-log BEACON_LOG_FILE] [--validator-log VALIDATOR_LOG_FILE] [--geth-log GETH_LOG_FILE]
    [--geth-binary GETH_BINARY] [--beacon-binary BEACON_BINARY] [--prysmctl-binary PRYSMCTL_BINARY]
    [--validator-binary VALIDATOR_BINARY] [--gethdata-dir GETHDATA_DIR] [--beacondata-dir BEACONDATA_DIR]
    [--validatordata-dir VALIDATORDATA_DIR]"
    exit 1
}

# Parse command-line arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --beacon-rpc) BEACON_RPC_PORT="$2"; shift ;;
        --geth-http) GETH_HTTP_PORT="$2"; shift ;;
        --geth-ws) GETH_WS_PORT="$2"; shift ;;
        --build-dir) BUILD_DIR="$2"; shift ;;
        --geth-binary) GETH_BINARY="$2"; shift ;;
        --beacon-binary) BEACON_BINARY="$2"; shift ;;
        --prysmctl-binary) PRYSMCTL_BINARY="$2"; shift ;;
        --validator-binary) VALIDATOR_BINARY="$2"; shift ;;
        --beacon-log) BEACON_LOG_FILE="$2"; shift ;;
        --validator-log) VALIDATOR_LOG_FILE="$2"; shift ;;
        --geth-log) GETH_LOG_FILE="$2"; shift ;;
        --gethdata-dir) GETHDATA_DIR="$2"; shift ;;
        --beacondata-dir) BEACONDATA_DIR="$2"; shift ;;
        --validatordata-dir) VALIDATORDATA_DIR="$2"; shift ;;
        *) usage ;;
    esac
    shift
done

mkdir -p "$(dirname "${BEACON_LOG_FILE}")"
mkdir -p "$(dirname "${VALIDATOR_LOG_FILE}")"
mkdir -p "$(dirname "${GETH_LOG_FILE}")"

echo "Beacon RPC Port: ${BEACON_RPC_PORT}"
echo "Geth HTTP Port: ${GETH_HTTP_PORT}"
echo "Geth WS Port: ${GETH_WS_PORT}"
echo "Build Directory: ${BUILD_DIR}"
echo "Geth Data Directory: ${GETHDATA_DIR}"
echo "Beacon Data Directory: ${BEACONDATA_DIR}"
echo "Validator Data Directory: ${VALIDATORDATA_DIR}"
echo "Geth Log: ${GETH_LOG_FILE}"
echo "Beacon Log: ${BEACON_LOG_FILE}"
echo "Validator lod: ${VALIDATOR_LOG_FILE}"

if [ ! -f "${BEACON_BINARY}" ]; then
    echo "Error: Beacon binary not found at ${BEACON_BINARY}"
    exit 1
fi

if [ ! -f "${PRYSMCTL_BINARY}" ]; then
    echo "Error: Prysmctl binary not found at ${PRYSMCTL_BINARY}"
    exit 1
fi

if [ ! -f "${VALIDATOR_BINARY}" ]; then
    echo "Error: Validator binary not found at ${VALIDATOR_BINARY}"
    exit 1
fi

# Run the Prysm beacon node
${BEACON_BINARY} --datadir="${BEACONDATA_DIR}" \
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
               --execution-endpoint "${GETHDATA_DIR}/geth.ipc" > "${BEACON_LOG_FILE}" 2>&1 &
echo "Beacon node started"

# Run Prysm validator client
${VALIDATOR_BINARY} --beacon-rpc-provider=127.0.0.1:4000 \
            --datadir="${VALIDATORDATA_DIR}" \
            --accept-terms-of-use \
            --interop-num-validators 2 \
            --chain-config-file config.yml > "${VALIDATOR_LOG_FILE}" 2>&1 &
echo "Validator client started"

# Run go-ethereum
${GETH_BINARY} --http \
       --http.api eth,net,web3 \
       --ws --ws.api eth,net,web3 \
       --authrpc.jwtsecret jwt.hex \
       --datadir="${GETHDATA_DIR}" \
       --nodiscover \
       --syncmode full \
       --allow-insecure-unlock \
       --unlock 0x123463a4b065722e99115d6c222f267d9cabb524 \
       --password ./password.txt > "${GETH_LOG_FILE}" 2>&1 &

echo "Geth network started"
sleep 600