#!/usr/bin/env bash

#
# This script displays the Obscuro node health status
#

help_and_exit() {
    echo ""
    echo "Usage: "
    echo "      -  $(basename "${0}") "
    echo ""
    exit 1  # Exit with error explicitly
}
# Ensure any fail is loud and explicit
set -euo pipefail

# Define local usage vars
start_path="$(cd "$(dirname "${0}")" && pwd)"
testnet_path="${start_path}"

# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in

            --help)                     help_and_exit ;;
            *)
    esac
done


net_status=$(curl -s --request POST 'http://127.0.0.1:80' \
             --header 'Content-Type: application/json' \
             --data-raw '{ "method":"obscuro_health", "params":null, "id":1, "jsonrpc":"2.0" }')

echo "Health Status: ${net_status}"
echo ""
echo "Container Status:"
docker inspect --format "Container: {{.Name}} - Status: {{.State.Status}} {{println }}Created at: {{.Created}}{{println }}Arguments: {{println }}{{range .Args}}{{println .}}{{end}} " $(docker-compose ps -a -q)
