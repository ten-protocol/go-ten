#!/usr/bin/env bash

#
# This script gracefully stops the Obscuro node.
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


docker-compose stop