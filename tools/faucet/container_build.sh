#!/usr/bin/env bash

#
# This script builds the faucet server image locally
#

# Ensure any fail is loud and explicit
set -euo pipefail

# run the build
docker build -t obscuronet/faucet_testnet:latest -f ./Dockerfile .
