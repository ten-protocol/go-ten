#!/usr/bin/env bash

#
# This script checks and waits for the eth2network to be post-merge ready
#
#

help_and_exit() {
    echo ""
    echo "Usage: "
    echo "   ex: "
    echo "      -  $(basename "${0}") --host=127.0.0.1 --port=8025"
    echo ""
    echo "  host             *Required* Set the host address"
    echo ""
    echo "  port             *Optional* Set the http host port. Defaults to 8025"
    echo ""
    echo "  timeout          *Optional* Set timeout in seconds. Defaults to 5*60 seconds"
    echo ""
    exit 1  # Exit with error explicitly
}

# Ensure any fail is loud and explicit
set -euo pipefail

# Defaults
port=8025
max_attempts=600
merge_block_height=7


# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --host)                   host=${value} ;;
            --port)                   port=${value} ;;
            --timeout)                max_attempts=${value} ;;

            --help)                     help_and_exit ;;
            *)
    esac
done

if [[ -z ${host:-} ]];
then
    help_and_exit
fi

url="http://${host}:${port}"
payload='{"jsonrpc": "2.0", "method": "eth_blockNumber", "params": [], "id": 1}'

attempts=0
while [ $attempts -lt $max_attempts ]; do
  response=$(curl -s -X POST -H "Content-Type: application/json" -d "$payload" "$url") || true
  status=$(echo $response | grep -o "result" ) || echo false

  if [ $status ]; then
    result=$(echo $response | grep -oE '0x[0-9a-fA-F]+')
    if [ -n "$result" ] && [ $((16#${result:2})) -gt $merge_block_height ]; then
      echo "Success: Response is 200 OK and block height is greater than $merge_block_height"
      break
    else
      echo "Failed: block height field not found or is not greater than $merge_block_height"
      echo "Response: $response"
      attempts=$((attempts + 1))
      sleep 1
    fi
  else
      echo "Failed: No 200 OK from the eth2network"
      attempts=$((attempts + 1))
      sleep 1
    fi
done

if [ $attempts -eq $max_attempts ]; then
  echo "Exceeded maximum number of attempts, giving up"
fi

echo "Node up and running!"
exit 0