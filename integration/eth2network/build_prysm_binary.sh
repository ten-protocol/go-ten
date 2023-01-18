#!/usr/bin/env bash

#
# This script downloads and builds prysm from source
# Requires to specify the version to clone
#

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}") --version=v3.2.0 "
    echo ""
    echo "  version       *Required* Set the version of prysm to build"
    echo ""
    echo ""
    exit 1  # Exit with error explicitly
}

# Ensure any fail is loud and explicit
set -euo pipefail

# Define local usage vars
script_path="$(cd "$(dirname "${0}")" && pwd)"
build_path="${script_path}/../.build"
prysm_repo_path="${build_path}/prysm_repo"
prysm_repo_bin_path="${prysm_repo_path}/build/bin/prysm"

# Define defaults
prism_path="${build_path}/prysm_bin"

# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --version)          prysm_version=${value} ;;
            --help)                     help_and_exit ;;
            *)
    esac
done

if [[ -z ${prysm_version:-} ]];
then
    help_and_exit
fi

# Make sure .build folder exists
mkdir -p "${build_path}"

# Only download geth code if binary does not exist
if [ -f "${prism_path}/prysm-${prysm_version}" ]
then
    echo "Skipping prysm build - Found binary at ${prism_path}/prysm-${prysm_version}"
    exit 0
fi

# Clone prysm source code if the path is empty
if [ -d "${prysm_repo_path}" ]
then
    echo "Skipping prysm repo clone - Found data in ${prysm_repo_path}"
else
    git clone --depth 1 --branch "${prysm_version}" https://github.com/prysmaticlabs/prysm.git "${prysm_repo_path}"
fi

# Build geth binary
cd "${prysm_repo_path}"
go build -o=../beacon-chain ./cmd/beacon-chain
go build -o=../validator ./cmd/validator
go build -o=../prysmctl ./cmd/prysmctl

## Copy binary to the correct path
#mkdir -p "${prism_path}"
#cp "${prism_repo_bin_path}" "${prism_path}/geth-${prysm_version}"
#
#cd ..
#
## Delete Geth repo
#echo "Deleting geth repo clone - Found data in ${prism_repo_path}"
#rm -rf "${prism_repo_path}"
