#!/usr/bin/env bash

#
# This script builds the faucet server image locally
#

# Ensure any fail is loud and explicit
set -euo pipefail

# cd to root of the repository
cd ../../

pwd

# run the build
docker build -t obscuronet/obscuro_gateway_testnet:latest -f ./tools/walletextension/Dockerfile .
