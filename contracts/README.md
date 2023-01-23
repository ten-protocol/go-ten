# Obscuro smart contracts

This hardhat project contains the relevant smart contracts for the Obscuro L2 platform.

## Dependencies

NodeJS LTS (v18)
NPM (tested with 8.19.3)

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

The command internally uses the abi and bytecode exporter plugins and searches the path configured in their configs for exporting for relevant files in order to launch the `abigen` executable with the correct paramaters. More info on installing `abigen` can be found [here](https://geth.ethereum.org/docs/dapp/abigen)


Additionally you can pass the `noCompile` flag which will disable running the contract compilation beforehand. This allows to build go bindings for abi/bins where the actual solidity source files are missing.

## Compilation

The following command compiles the solidity contracts and produces artifacts for them:

```shell
npx hardhat compile
```

[https://www.npmjs.com/package/hardhat-ignore-warnings](The ignore warnings plugin) can be used to configure the behavior of the `compile` task in regards to warnings and errors.