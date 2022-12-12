# Obscuro smart contracts

This hardhat project contains the relevant smart contracts for the Obscuro L2 platform.

## Dependencies

NodeJS LTS (v18)
NPM 

## Installing

Running the following command will pull all of the relevant dependencies for node and solidity.

```shell
npm install
``` 

## Generating abi bindings for GO

Running the following command will regenerate the bindings in the `generated` directory:

```shell
npx hardhat generate-abi-bindings --output-dir generated
```