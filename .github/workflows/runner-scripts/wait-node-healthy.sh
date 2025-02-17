#!/usr/bin/env bash

#
# This script checks and waits for an obscuro node to be healthy
#
#

help_and_exit() {
    echo ""
    echo "Usage: "
    echo "   ex: "
    echo "      -  $(basename "${0}") --host=erpc.uat-testnet.ten.xyz --port=80"
    echo ""
    echo "  node             *Required* Set the host address"
    echo ""
    echo "  port             *Optional* Set the host port. Defaults to 80"
    echo ""
    echo "  timeout          *Optional* Set timeout in seconds. Defaults to 5*60 seconds"
    echo ""
    exit 1  # Exit with error explicitly
}
# Ensure any fail is loud and explicit
set -euo pipefail

# Define local usage vars
start_path="$(cd "$(dirname "${0}")" && pwd)"
testnet_path="${start_path}"

# Defaults
port=80
timeout=5*60

# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --host)                   host=${value} ;;
            --port)                   port=${value} ;;
            --timeout)                timeout=${value} ;;

            --help)                     help_and_exit ;;
            *)
    esac
done

if [[ -z ${host:-} ]];
then
    help_and_exit
fi

net_status=""
time=0

echo "Running health check against http://${host}:${port}"

while ! [[ $net_status = *\"OverallHealth\":true* ]]
do
    net_status=$(curl --request POST "http://${host}:${port}" \
                 --header 'Content-Type: application/json' \
                 --data-raw '{ "method":"ten_health", "params":null, "id":1, "jsonrpc":"2.0" }') || true
    echo $net_status

    sleep 2
    ((time=time+2))

    if [[ $time == $timeout ]] ;
    then
      echo "Node not healthy after timeout!"
      exit 1
    fi
done

echo "Node up and running!"
exit 0