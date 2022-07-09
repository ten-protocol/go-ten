#!/usr/bin/env bash

#
# This script builds the management contract from solidity to a go generated package
#

# Ensure any fail is loud and explicit
set -euo pipefail

# Define local usage vars
script_path="$(cd "$(dirname "${0}")" && pwd)"
contract_path="${script_path}/management_contract.sol"
abi_path="${script_path}/management_contract.abi"
bin_path="${script_path}/management_contract.bin"
libs_path="${script_path}/libs"
package_path="${script_path}/compiledcontracts"
mgmt_contract_package="ManagementContract"
management_package_path="${package_path}/generated${mgmt_contract_package}"
mgmt_contract_name="ManagementContract"

# ensure folder exists
mkdir -p "${management_package_path}"

# generate the abi
solc --abi -o "${abi_path}" "${contract_path}" --overwrite
solc --bin -o "${bin_path}" "${contract_path}" --overwrite
# generates the golang package
abigen --abi="${abi_path}/${mgmt_contract_name}.abi" --bin="${bin_path}/${mgmt_contract_name}.bin"  --pkg="${mgmt_contract_package}" --out="${management_package_path}/ManagementContract.go"
