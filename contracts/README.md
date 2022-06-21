# Network contracts folder

This folder holds the contracts that are deployed into the L1 network.

## Contracts

* `management_contract.sol` - stores the rollups and the list of attested aggregators

## Contract deployment

Currently, the contract is automatically deployed in the simulation.
The deployment is based off the compilation of the contract and using the bytecode output.

Contract compilation is being done using https://remix.ethereum.org, as we are yet to automate it.
