#!/usr/bin/env bash

#
# This script builds the ManagementContract from solidity to a go generated package
#

# Ensure any fail is loud and explicit
set -euo pipefail

# Define local usage vars
contract_name="ManagementContract"

script_path="$(cd "$(dirname "${0}")" && pwd)"
contract_path="${script_path}/${contract_name}.sol"
abi_path="${script_path}/abi"
bin_path="${script_path}/bin"
libs_path="${script_path}/libs"
generated_path="${script_path}/generated/${contract_name}"

# ensure folder exists
mkdir -p "${generated_path}"

# generate the abi
solc --base-path "${script_path}" --include-path "${libs_path}" --abi -o "${abi_path}" "${contract_path}" --overwrite
solc --base-path "${script_path}" --include-path "${libs_path}" --bin -o "${bin_path}" "${contract_path}" --overwrite

# generates the golang package
abigen --abi="${abi_path}/${contract_name}.abi" --bin="${bin_path}/${contract_name}.bin" --pkg="${contract_name}" --out="${generated_path}/${contract_name}.go"
