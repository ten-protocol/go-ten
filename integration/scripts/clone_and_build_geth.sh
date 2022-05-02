#!/usr/bin/env bash

# This script downloads and builds geth from source
# It's possible to specify the version to clone and the geth binary output

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}") --version=v1.10.17 --output=./build/geth_output "
    echo ""
    echo "  version        Set the version of geth to clone + download"
    echo ""
    echo "  output         Set the where the geth binary will be built to"
    echo ""
    echo ""
    exit 1  # Exit with error explicitly
}

# Ensure any fail is loud and explicit
set -euo pipefail

# Define local usage vars ( Don't define CAPS version )
script_path="$(cd "$(dirname "${0}")" && pwd)"
build_path="${script_path}/../.build"
geth_repo_path="${build_path}/geth_repo"

# Define defaults
geth_path="${build_path}/geth_bin"
geth_version="v1.10.17"

# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --version)          geth_version=${value} ;;
            --output)           geth_path=${value} ;;
            --help)                   show_help_and_exit ;;
            *)
    esac
done

# Make sure .build folder exists
mkdir -p "${build_path}"

# Clone geth source code if the path is empty
if [ -d "${geth_repo_path}" ]
then
    echo "Skipping geth repo clone - Found data in ${geth_repo_path}"
else
    git clone --depth 1 --branch "${geth_version}" https://github.com/ethereum/go-ethereum "${geth_repo_path}"
fi

# Build geth
cd "${geth_repo_path}"
make geth

# Copy binary to the correct path
mkdir -p "${geth_path}"
cp "${geth_repo_path}/build/bin/geth" "${geth_path}/geth"
