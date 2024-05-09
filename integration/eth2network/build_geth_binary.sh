#!/usr/bin/env bash

#
# This script downloads and builds geth from source
# Requires to specify the version to clone
#
echo "build geth binary"

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}") --version=v1.14.0 --output=path_to_output"
    echo ""
    echo "  version       *Required* Set the version of geth to build"
    echo "  output        *Required* Where to copy the binary"
    echo ""
    echo ""
    exit 1  # Exit with error explicitly
}

# Ensure any fail is loud and explicit
set -euo pipefail

# Define local usage vars
script_path="$(cd "$(dirname "${0}")" && pwd)"
build_path="${script_path}/../.build"
geth_repo_path="${build_path}/geth_repo"
geth_repo_bin_path="${geth_repo_path}/build/bin/geth"

# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --version)          geth_version=${value} ;;
            --output)           geth_path=${value} ;;
            --help)                     help_and_exit ;;
            *)
    esac
done

if [[ -z ${geth_version:-} ]];
then
    help_and_exit
fi

# Make sure .build folder exists
mkdir -p "${build_path}"

# Clone geth source code if the path is empty
if [ -d "${geth_repo_path}" ]
then
    echo "Skipping geth repo clone - Found data in ${geth_repo_path}"
else
    git clone --depth 1 --branch "${geth_version}" https://github.com/ethereum/go-ethereum "${geth_repo_path}"
fi

# Build geth binary
cd "${geth_repo_path}"
export GOROOT=
make geth

# Copy binary to the correct path
mkdir -p "${geth_path}"
cp "${geth_repo_bin_path}" "${geth_path}/geth"

cd ..

# Delete Geth repo
echo "Deleting geth repo clone - Found data in ${geth_repo_path}"
rm -rf "${geth_repo_path}"
