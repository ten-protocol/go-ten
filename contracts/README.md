# TEN smart contracts

This hardhat project contains the relevant smart contracts for the TEN L2 platform.

For more details check the ``src`` folder.

## Dependencies

NodeJS LTS (v18)
NPM (tested with 8.19.3)

## Installing

First you'll need to install ethereum locally

```shell
brew tap ethereum/ethereum
brew install ethereum
```

Notice that for generating the go bindings, you need to have the `abigen` executable installed, which comes from the ethereum package, but can be independently installed using (recommended for the devcontainer):

```shell
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```

Running the following command will pull all of the relevant dependencies for node and solidity.

```shell
npm install
``` 

## Generating abi bindings for GO

Running the following command will regenerate the bindings in the `generated` directory:

```shell
npx hardhat generate-abi-bindings --output-dir generated
```

The command internally uses the abi and bytecode exporter plugins and searches the path configured in their configs for exporting for relevant files in order to launch the `abigen` executable with the correct parameters. More info on installing `abigen` can be found [here](https://geth.ethereum.org/docs/dapp/abigen)


Additionally you can pass the `noCompile` flag which will disable running the contract compilation beforehand. This allows to build go bindings for abi/bins where the actual solidity source files are missing.

## Deploying

### Deployment Scripts folder structure

 * core - Scripts required to be predeployed for TEN to start.
 * bridge - Scripts that deploy/upgrade ONLY the bridge.
 * messenger - Scripts that enable the relayer functionality. Can contain predeployed libraries too in the future.
 * testnet - Scripts that should only be deployed on the testnet. Tokens, "dev tooling" scripts, etc.

For deployments, we use the hardhat-deploy plugin. It provides the `deploy` task, which determines what folders with deployment scripts need to be executed for the current selected network. Additionally there is the `ten:deploy` task that will launch a wallet extension.
For the wallet extension to work, the network needs to have configured the `url` to 127.0.0.1:3000 and the additional `tenEncRpcUrl` property to the rpc endpoint of the TEN node the wallet will connect to.

Scripts are taken from `deployment_scripts` and executed in alphabetic order. Each folder, as ordered in the network config and then inside of it alphabetically. Notice that `func.dependencies = []` defined in deployment functions has the ability to escape the default ordering. If such a deployment function/script is reached, the deploy plugin will first deploy its dependency if it hasn't already.

## Compilation

The following command compiles the solidity contracts and produces artifacts for them:

```shell
npx hardhat compile
```

[https://www.npmjs.com/package/hardhat-ignore-warnings](The ignore warnings plugin) can be used to configure the behavior of the `compile` task in regards to warnings and errors.
