# Begin your Journey

First, thanks for taking your precious time to help us ensure testnet is all it can be. 🙏 

## Features and setting expectations

In this release, you will see the following features:

- Privacy across the entire lifecycle of your smart contract by encrypting specific information / API requests considered sensitive, for example, eth_getBalance. The full list of encrypted requests and responses is below.

https://docs.obscu.ro/wallet-extension/handling-sensitive-data.html#sensitive-json-rpc-api-methods

* An Ethereum-like Layer 1. Soon Obscuro Testnet will be hooked into a well-known Ethereum Testnet.

* Straightforward integration with MetaMask using a Wallet Extension binary run locally to hook your MetaMask into Obscuro Testnet.

https://docs.obscu.ro/wallet-extension/wallet-extension.html

* A rudimentary block explorer called ObscuroScan. For Testnet ONLY, this shows the details of transactions in a decrypted format by default to make it easier to understand what is happening on Testnet.

http://testnet-obscuroscan.uksouth.azurecontainer.io/ 

Current Limitations:

We're in our first release, and much more is yet to come. However, some of the current limitations you can expect:

* The gas feature is not implemented yet, nor is native value transfer.

* There is no streaming of events. This will limit how your UI updates based on events being listened to.

* Testnet might crash unexpectedly as sharp edges are smoothed away and edge cases fixed.

* Expect to lose all data when Testnet crashes or is restarted. 

## Challenges

Below are a set of challenges to get you started. For each challenge you successfully complete, you'll be rewarded task points which get transformed into OBX tokens as per the [Tokenomics](https://github.com/obscuronet/obscuro-project/wiki/Tokenomics)

### Warmup - 20 points
Get set up and deploy a simple contract.
* Configure the Wallet extension and Metamask following the instructions [here](https://docs.obscu.ro/wallet-extension/wallet-extension.html)
* Try deploying any contract using Remix (recommended) (e.g. the default Storage contract). Remember to select Injected Provider as your Environment ![Step 1 screenshot](https://images.tango.us/public/screenshot_017a4a1c-b655-4eda-8c6e-7a298e1faa70.png?crop=focalpoint&fit=crop&fp-x=0.0945&fp-y=0.1827&fp-z=2.5115&w=1200&mark-w=0.2&mark-pad=0&mark64=aHR0cHM6Ly9pbWFnZXMudGFuZ28udXMvc3RhdGljL21hZGUtd2l0aC10YW5nby13YXRlcm1hcmsucG5n&ar=2027%3A1103)

### Deploy the Guessing game - 20 points
The guessing game is a very basic contract that showcases some of Obscuro's capabilities. The goal of the game is to simply guess the secret number the game has selected. Each time an attempt is made, an entrance fee of 1 token is paid. If a player guesses correctly, they win the entire pot and a new random number is selected.

Before we get into how the game works, let's see it in action and play. See the instructions [here](https://docs.obscu.ro/testnet/example-dapps.html)

Now that we've played the game let's try deploying the Guessing game to the Obscuro testnet. You can clone the repo from [here](https://github.com/obscuronet/number-guessing-game). 

#### How the Guessing game works
There are two main functions of interest. First, the *attempt* function allows a player to guess the number and check if they've correctly guessed the secret number. If they have, the prize pot is paid out to them.
![attempt function](../../assets/images/guessing.png)

The other is the *setNewTarget* function. This checks the game is in a state ready to be reset and then sets a secret number based on various block parameters.
![setNewTarget function](../../assets/images/setnewtarget.png)

### Extend the Guessing game with Warmer/Colder functionality - 50 points
Currently, after each guess, you have no idea how far you are from the winner number. The game will be more fun if players know with each play whether they're closer (warmer) or further (colder) from their last guess.

The challenge here is to extend the game to include this functionality. You'll need to think about how you store previous guesses for each player.

### Extend the Guessing game with a proper win notification - 75 points
You'll have noticed immediately that on winning the game, there is no notification! The only thing that happens is the game transfers to you all the tokens in the pot. 

The challenge here is to extend the game with events (or anything else) to inform the winning player that they've won and how much they've won.

You'll need to write a UI for this.

### Deploy a DEX - 100 points
A key promise of Obscuro is that existing Dapps built for EVM-based chains will also work on Obscuro. So the challenge here is to deploy an existing evm-based DEX of your choice to Obscuro.

### Deploy a lending protocol - 100 points
Building on the previous challenge, can you deploy an existing EVM-based lending protocol to Obscuro.

### Make liquidations private and everything else public in the lending protocol - 150 points
Obscuro allows developers at a granular level to select what parts of their Dapp they wish to keep private vs public with powerful results. Read this [blog post](https://medium.com/obscuro-labs/the-obscuro-experience-26d1697a5378) to understand how this works. The challenge here is to take an existing Dapp, such as a DEX or lending protocol, and change it such that some parts are private and others are public for a compelling reason, e.g. keeping liquidation levels private in a lending protocol while ensuring everything else is public.

### Make guessing game truly random - 150 points
Currently, the Guessing game isn't truly random. Any developer deploying the game can calculate the secret number depending on when the game was initialised. The challenge is to extend the Guessing game to truly make the secret number random such that even the developers won't know. You'll need to co-ordinate with the Obscuro Labs team to build this. 

## Support

We're here to help. If you have issues, please post your questions on [Discord](https://discord.gg/obscuro).