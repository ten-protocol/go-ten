# Deploying a Smart Contract to Obscuro Testnet
Using the steps below you will add an extension to your MetaMask wallet so it can connect to Obscuro Testnet, then using Remix you will deploy your smart contract to Testnet.

## Prerequisites
* [MetaMask](https://metamask.io/) wallet installed in your browser.
* A local copy of the Obscuro MetaMask wallet extension downloaded from the [releases area](https://github.com/obscuronet/go-obscuro/releases) Obscuro repository on GitHub.

## Prepare Your MetaMask Wallet for Obscuro Testnet
An essential part of how Obscuro provides full privacy is the encryption of communication between an Obscuro application and Obscuro nodes on the network.

Follow the steps to config the MetaMask wallet extension [here](wallet-extension.md) then return to this page. If you do not have the Obscuro wallet extension running MetaMask will not be able to communicate with the Obscuro Testnet. 

Now your wallet is configured for the Obscuro Testnet which makes your encrypted traffic viewable. You can now go ahead and deploy your smart contract to the Obscuro Testnet.

## Deploy Your Smart Contract Using an IDE
1. Browse to the popular Solidity-compatible Integrated Development Environment called [Remix](https://remix.ethereum.org/).

1. Load your solidity smart contract into the Remix workspace.

1. Compile your smart contract using the Remix Solidity Compiler.

1. Open MetaMask and confirm you are connected to Obscuro Testnet network. The parameters for the Obscuro Testnet can be found [here](./essentials.md).

1. In the _Deploy & Run Transactions_ section of Remix change the Environment to _Injected Web3_. This tells Remix to use the network settings currently configured in your MetaMask wallet, which in this case is the Obscuro Testnet.

    If the connection to Obscuro Testnet is successful you will see the text _Custom (777) network_ displayed under _Injected Web3_.

1. Click the _Deploy_ button to deploy your smart contract to the Obscuro Testnet.

1. MetaMask will automatically open and ask you to confirm the deployment. Click _Confirm_.

Congratulations, your smart contract is now deployed to Obscuro Testnet!

Because Obscuro provides full privacy the details of any transactions with your smart contract are encrypted and only visible to those with the viewing key.

Now head over to the [ObscuroScan page](./obscuroscan.md) to see how you can view the transaction details.
