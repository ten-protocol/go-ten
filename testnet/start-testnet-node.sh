#!/usr/bin/env bash

#
# This script downloads and builds the obscuro node
#

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}") "
    echo ""
    echo ""
    echo ""
    exit 1  # Exit with error explicitly
}

# Ensure any fail is loud and explicit
set -euo pipefail

# Define local usage vars
start_path="$(cd "$(dirname "${0}")" && pwd)"

sudo apt-get update
sudo apt-get install -y docker.io docker-compose
