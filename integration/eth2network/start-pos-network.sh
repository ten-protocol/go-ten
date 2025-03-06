#!/usr/bin/env bash

# Default port values
GETH_NETWORK_PORT=30303
BEACON_P2P_PORT=12000
GETH_HTTP_PORT=8025
GETH_WS_PORT=9000
GETH_RPC_PORT=8552
BEACON_RPC_PORT=4000
BEACON_GATEWAY_PORT=3500
CHAIN_ID=1337
BUILD_DIR="./build"
BASE_PATH="./"
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
TEST_LOG_FILE="./test.log"

# Function to display usage
usage() {
    echo "Usage: $0
    [--geth-http GETH_HTTP_PORT]
    [--geth-ws GETH_WS_PORT]
    [--geth-rpc GETH_RPC_PORT]
    [--beacon-rpc BEACON_RPC_PORT]
    [--chainid CHAIN_ID ]
    [--build-dir BUILD_DIR ]
    [--base-path BASE_PATH ]
    [--beacon-log BEACON_LOG_FILE]
    [--validator-log VALIDATOR_LOG_FILE]
    [--geth-log GETH_LOG_FILE]
    [--geth-binary GETH_BINARY]
    [--beacon-binary BEACON_BINARY]
    [--prysmctl-binary PRYSMCTL_BINARY]
    [--validator-binary VALIDATOR_BINARY]
    [--gethdata-dir GETHDATA_DIR]
    [--beacondata-dir BEACONDATA_DIR]
    [--validatordata-dir VALIDATORDATA_DIR] "
    exit 1
}

# Parse command-line arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --geth-network) GETH_NETWORK_PORT="$2"; shift ;;
        --beacon-p2p) BEACON_P2P_PORT="$2"; shift ;;
        --beacon-rpc) BEACON_RPC_PORT="$2"; shift ;;
        --grpc-gateway-port) BEACON_GATEWAY_PORT="$2"; shift ;;
        --geth-http) GETH_HTTP_PORT="$2"; shift ;;
        --geth-ws) GETH_WS_PORT="$2"; shift ;;
        --geth-rpc) GETH_RPC_PORT="$2"; shift ;;
        --chainid) CHAIN_ID="$2"; shift ;;
        --build-dir) BUILD_DIR="$2"; shift ;;
        --base-path) BASE_PATH="$2"; shift ;;
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
        --test-log) TEST_LOG_FILE="$2"; shift ;;
        *) usage ;;
    esac
    shift
done

mkdir -p "$(dirname "${BEACON_LOG_FILE}")"
mkdir -p "$(dirname "${VALIDATOR_LOG_FILE}")"
mkdir -p "$(dirname "${GETH_LOG_FILE}")"
mkdir -p "$(dirname "${TEST_LOG_FILE}")"

echo "Test" > "${TEST_LOG_FILE}" 2>&1 &

${PRYSMCTL_BINARY} testnet generate-genesis \
           --fork deneb \
           --num-validators 2 \
	         --genesis-time-delay 15 \
	         --output-ssz "${BEACONDATA_DIR}/genesis.ssz" \
           --chain-config-file "${BASE_PATH}/config.yml" \
           --geth-genesis-json-in "${BUILD_DIR}/genesis.json" \

sleep 1
echo "Prysm genesis generated"

echo -e "\n\n" | ${GETH_BINARY} --datadir="${GETHDATA_DIR}" account import "${BASE_PATH}/pk.txt"
echo "Private key imported into gethdata"

${GETH_BINARY} --datadir="${GETHDATA_DIR}" init  --state.scheme=hash "${BUILD_DIR}/genesis.json"
sleep 1
echo "Geth genesis initialized"

# First, check if any of these ports are already in use
if netstat -tuln | grep -q ":${GETH_HTTP_PORT} "; then
    echo "Port ${GETH_HTTP_PORT} is already in use. Please stop any existing Geth instances or use a different port."
    exit 1
fi

if netstat -tuln | grep -q ":${GETH_WS_PORT} "; then
    echo "Port ${GETH_WS_PORT} is already in use. Please stop any existing Geth instances or use a different port."
    exit 1
fi

if netstat -tuln | grep -q ":${GETH_RPC_PORT} "; then
    echo "Port ${GETH_RPC_PORT} is already in use. Please stop any existing Geth instances or use a different port."
    exit 1
fi

# Update Geth command to match the Docker configuration
${GETH_BINARY} --http \
      --datadir="${GETHDATA_DIR}" \
      --http.addr=0.0.0.0 \
      --authrpc.vhosts=* \
      --authrpc.port="${GETH_RPC_PORT}" \
      --authrpc.addr=0.0.0.0 \
      --http.vhosts=* \
      --http.api=engine,eth,net,web3 \
      --http.corsdomain=* \
      --ipcdisable \
      --ws \
      --ws.addr=0.0.0.0 \
      --ws.api=eth,net,web3,debug,txpool \
      --allow-insecure-unlock \
      --unlock=0x3f1Eae7D46d88F08fc2F8ed27FCb2AB183EB2d0E \
      --password "${BASE_PATH}/password.txt" \
      --authrpc.jwtsecret "${BASE_PATH}/jwt.hex" \
      --nodiscover \
      --syncmode=full \
      --state.scheme=hash \
      --dev \
      --dev.period=1 \
      --mine \
      --miner.etherbase=0x3f1Eae7D46d88F08fc2F8ed27FCb2AB183EB2d0E \
      --gcmode=archive  > "${GETH_LOG_FILE}" 2>&1 &
#       --http.api admin,eth,net,web3,debug \
#       --http.addr="0.0.0.0" \
#       --http.port="${GETH_HTTP_PORT}" \
#       --http.corsdomain "*" \
#       --http.vhosts "*" \
#       --ws --ws.api admin,eth,net,web3,debug,txpool \
#       --ws.addr="0.0.0.0" \
#       --ws.port="${GETH_WS_PORT}" \
#       --ws.origins "*" \
#       --authrpc.jwtsecret "${BASE_PATH}/jwt.hex" \
#       --authrpc.port "${GETH_RPC_PORT}" \
#       --authrpc.addr="0.0.0.0" \
#       --authrpc.vhosts "*" \
#       --port="${GETH_NETWORK_PORT}" \
#       --datadir="${GETHDATA_DIR}" \
#       --networkid="${CHAIN_ID}" \
#       --nodiscover \
#       --dev \
#       --dev.period 1 \
#       --syncmode full \
#       --password "${BASE_PATH}/password.txt" > "${GETH_LOG_FILE}" 2>&1 &
geth_pid=$!
echo "GETH PID $geth_pid"

# Run the Prysm beacon node
${BEACON_BINARY} --datadir="${BEACONDATA_DIR}" \
    --rpc-port="${BEACON_RPC_PORT}" \
    --min-sync-peers=0 \
    --genesis-state "${BEACONDATA_DIR}/genesis.ssz" \
    --interop-eth1data-votes \
    --bootstrap-node= \
    --chain-config-file "${BASE_PATH}/config.yml" \
    --rpc-host=0.0.0.0 \
    --grpc-gateway-host=0.0.0.0 \
    --grpc-gateway-port=${BEACON_GATEWAY_PORT} \
    --chain-id "${CHAIN_ID}" \
    --execution-endpoint="http://127.0.0.1:${GETH_RPC_PORT}" \
    --accept-terms-of-use \
    --jwt-secret "${BASE_PATH}/jwt.hex" \
    --suggested-fee-recipient=0x00000000000000000000000000000000 > "${BEACON_LOG_FILE}" 2>&1 &
#               --min-sync-peers 0 \
#               --genesis-state "${BEACONDATA_DIR}/genesis.ssz" \
#               --bootstrap-node= \
#               --interop-eth1data-votes \
#               --chain-config-file "${BASE_PATH}/config.yml" \
#               --contract-deployment-block 0 \
#               --chain-id "${CHAIN_ID}" \
#               --rpc-host=0.0.0.0 \
#               --rpc-port="${BEACON_RPC_PORT}" \
#               --p2p-udp-port="${BEACON_P2P_PORT}" \
#               --grpc-gateway-host=0.0.0.0 \
#               --accept-terms-of-use \
#               --jwt-secret "${BASE_PATH}/jwt.hex" \
#               --suggested-fee-recipient 0x123463a4B065722E99115D6c222f267d9cABb524 \
#               --minimum-peers-per-subnet 0 \
#               --enable-debug-rpc-endpoints \
#               --verbosity=debug \
#               --execution-endpoint "http://127.0.0.1:${GETH_HTTP_PORT}" > "${BEACON_LOG_FILE}" 2>&1 &
beacon_pid=$!
echo "BEACON PID $beacon_pid"

# Run Prysm validator client
${VALIDATOR_BINARY} --beacon-rpc-provider=127.0.0.1:"${BEACON_RPC_PORT}" \
            --datadir="${VALIDATORDATA_DIR}" \
            --accept-terms-of-use \
            --interop-num-validators 2 \
            --interop-start-index 0 \
            --chain-config-file "${BASE_PATH}/config.yml" > "${VALIDATOR_LOG_FILE}" 2>&1 &
validator_pid=$!
echo "VALIDATOR PID $validator_pid"
