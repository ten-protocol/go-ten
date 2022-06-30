# Deploying a Smart Contract to Obscuro Testnet
Using the steps below you will add an extension to your MetaMask wallet so it can connect to Obscuro Testnet, then using Remix you will deploy your smart contract to Testnet.

## Prerequisites
* MetaMask installed in your browser.
* A local copy of the Obscuro MetaMask wallet extension downloaded from the [Obscuro repo on GitHub](https://github.com/obscuronet/go-obscuro/tree/main/tools/walletextension).

## Prepare Your MetaMask Wallet for Obscuro Testnet
An essential part of how Obscuro provides full privacy is the encryption of communication between an Obscuro application and Obscuro nodes on the network.

Follow the steps to config the MetaMask wallet extension [here](https://docs.obscu.ro/testnet/wallet-extension.html) then return to this page.

Now your wallet is configured to make your encrypted traffic viewable you can go ahead and deploy your smart contract to the Obscuro Testnet.

## Deploy Your Smart Contract Using an IDE
1. Open your favourite Solidity-compatible Integrated Development Environment. The steps below will assume you are using [Remix](https://github.com/ethereum/remix-ide).

1. Open your solidity smart contract in Remix.

1. Open MetaMask and confirm you are connected to Obscuro Testnet network.

1. In the _Deploy & Run Transactions_ section of Remix change the Environment to _Injected Web3_. This tells Remix to use the network settings currently configured in your MetaMask wallet, which in this case is the Obscuro Testnet.

1. Click the _Deploy_ button to deploy your smart contract to the Obscuro Testnet.
