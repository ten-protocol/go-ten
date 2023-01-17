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

## Deploying

For deployments, we use the hardhat-deploy plugin. It provides the `deploy` task, which determines what folders with deployment scripts need to be executed for the current selected network. Additionally there is the `obscuro:deploy` task that will launch a wallet extension.
For the wallet extension to work, the network needs to have configured the `url` to 127.0.0.1:3000 and the additional `obscuroEncRpcUrl` property to the rpc endpoint of the obscuro node the wallet will connect to.

Scripts are taken from `deployment_scripts` and executed in alphabetic order. Each folder, as ordered in the network config and then inside of it alphabetically. Notice that `func.dependencies = []` defined in deployment functions has the ability to escape the default ordering. If such a deployment function/script is reached, the deploy plugin will first deploy its dependency if it hasn't already.

## Using hardhat node instead of a geth network

The command `npx hardhat node` will run a hardhat node and automatically deploy the L1 bits to it. They will also match the expected addresses.
This allows to start an obscuro node against this node. The following modification of the parameters of the testnet script should work with the default node on macOS: 

```shell
./start-obscuro-node.sh --sgx_enabled=false --l1host=host.docker.internal --l1port=8545 --mgmtcontractaddr=0xeDa66Cc53bd2f26896f6Ba6b736B1Ca325DE04eF --hocerc20addr=0xC0370e0b5C1A41D447BDdA655079A1B977C71aA9 --pocerc20addr=0x51D43a3Ca257584E770B6188232b199E76B022A2 --is_genesis=true --node_type=sequencer
```