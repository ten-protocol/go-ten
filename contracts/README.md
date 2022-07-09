# Network contracts folder

This folder holds the contracts that are deployed onto the L1 network.

## Contracts

* `management_contract.sol` - stores the rollups and the list of attested aggregators

## Contract deployment

Currently, the contract is automatically deployed in the simulation. The deployment is based off the compilation of the contract and using the bytecode output.

The contract can be compiled by running the `build_mgmt_contract.sh` script. This script requires `solc` and `abigen` (`version` >= `1.10.20-stable`).
