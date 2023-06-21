#!/usr/bin/env bash

#
# This script starts up the Obscuro Gateway
#

# Ensure any fail is loud and explicit
set -euo pipefail

# Define defaults
port=3000
portWS=3001
host="127.0.0.1"
nodeHost="testnet.obscu.ro"
nodePortHTTP=13000
nodePortWS=13001
logPath="wallet_extension_logs.txt"
databasePath=".obscuro/gateway_database.db"

# Parse the options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --host)            host=${value} ;;
            --port)            port=${value} ;;
            --portWS)          portWS=${value} ;;
            --nodeHost)        nodeHost=${value} ;;
            --nodePortHTTP)    nodePortHTTP=${value} ;;
            --nodePortWS)      nodePortWS=${value} ;;
            --logPath)         logPath=${value} ;;
            --databasePath)    databasePath=${value} ;;
            *)
    esac
done

# Stop and remove any running container, and then star
echo "Force stopping any existing container ... "
docker rm -f  obscuro_gateway_testnet 2>/dev/null

echo "Starting Obscuro Gateway..."
docker run -p 3000:"${port}" --name=obscuro_gateway_testnet \
    --detach \
    --entrypoint ./tools/walletextension/main/obscuro-gateway \
     "obscuronet/obscuro_gateway_testnet:latest" \
      -host="${host}" -port="${port}" -portWS="${portWS}" -nodeHost="${nodeHost}" -nodePortHTTP="${nodePortHTTP}" \
      -nodePortWS="${nodePortWS}" -logPath="${logPath}" -databasePath="${databasePath}"