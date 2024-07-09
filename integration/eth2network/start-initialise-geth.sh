#!/usr/bin/env bash

GETH_BINARY="./geth"
PRYSMCTL_BINARY="./prysmctl"
GETHDATA_DIR="/gethdata"

# Function to display usage
usage() {
    echo "Usage: $0 [--geth-binary GETH_BINARY] [--prysmctl-binary PRYSMCTL_BINARY] [--gethdata-dir GETHDATA_DIR]"
    exit 1
}

# Parse command-line arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --geth-binary) GETH_BINARY="$2"; shift ;;
        --prysmctl-binary) PRYSMCTL_BINARY="$2"; shift ;;
        --gethdata-dir) GETHDATA_DIR="$2"; shift ;;
        *) usage ;;
    esac
    shift
done

echo "Geth Binary: ${GETH_BINARY}"
echo "Prysm Binary: ${PRYSMCTL_BINARY}"
echo "Geth Data Directory: ${GETHDATA_DIR}"

echo -e "\n\n" | ${GETH_BINARY} --datadir="${GETHDATA_DIR}" account import pk.txt
echo "Private key imported"

${GETH_BINARY} --datadir="${GETHDATA_DIR}" init genesis.json
echo "Geth genesis initialized"

${PRYSMCTL_BINARY} testnet generate-genesis \
           --fork deneb \
           --num-validators 2 \
           --chain-config-file config.yml \
           --output-ssz genesis.ssz
           #           --geth-genesis-json-in genesis.json \
           #           --geth-genesis-json-out genesis.json \
sleep 10