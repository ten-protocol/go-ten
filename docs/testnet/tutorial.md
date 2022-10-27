# Getting Started
In this tutorial, you build your own Obscuro dApp from the start. This is a good way to experience a typical development process as you learn how Obscuro dApps are built, design concepts, tools and terminology.

# Lets go
This 'Guessing Game' tutorial provides an introduction to the fundamentals of Obscuro and shows you how to:
- Set up your local Ethereum development environment for Obscuro.
- Use Obscuro concepts to develop a dApp

What is the *Guessing Game*?

The Guessing Game is a dApp that you will build that demonstrates a basic Obscuro use case, which is a simple number guessing game. The contract generates a random secret number when it's deployed, which is never revealed to an operator or end-user because of the privacy benefits of Obscuro. The goal of the game is to guess this number, and each time an attempt is made, an entrance fee of 1 token is paid. If a user correctly guesses the number, the contract will pay out all of the accumulated entrance fees to them, and reset itself with a new random number.

Without Obscuro, it would be possible to look up the internal state of the contract and cheat, and the game wouldn't work.

The dApp has many of the features that you'd expect to find, including:
- Generating a random number known only to the dApp
- Allow users to guess the numbers by depositing a token for each play
- Colder/Warmer functionality to help guide users to the correct number
- Events to alert users whether they have won or not
- Track all users plays
- Administrator privileges to start and stop the game

# Set up your environment
- You'll need to install MetaMask following the instructions [here](https://metamask.io/)
- Then download and run the Obscuro Wallet extension following the instructions [here](https://docs.obscu.ro/wallet-extension/wallet-extension.html)
- Now configure MetaMask following the instructions [here](https://docs.obscu.ro/wallet-extension/configure-metamask.html)
- Finally, check you can open the Remix IDE by visiting the following URL in your browser https://remix.ethereum.org

That's it, you're all set to start building your first dApp on Obscuro.

## 1. To begin, we'll clone a a basic repo from Github through Remix by following these instructions

### a. Click on GitHub from the main screen in Remix
![Step 1 screenshot](https://images.tango.us/workflows/454b9fc4-c6f9-43d1-beef-9cbbff7f0b0b/steps/7cb0a9d9-2e3d-4c33-b10c-4f89dea34206/a3c964d3-6d40-4931-9ddc-8390f5744452.png?crop=focalpoint&fit=crop&fp-x=0.2873&fp-y=0.4869&fp-z=2.9198&w=1200&mark-w=0.2&mark-pad=0&mark64=aHR0cHM6Ly9pbWFnZXMudGFuZ28udXMvc3RhdGljL21hZGUtd2l0aC10YW5nby13YXRlcm1hcmsucG5n&ar=1812%3A1226)


### b. And enter "https://github.com/obscuronet/sample-applications/blob/main/contracts/GuessWarmerColder.sol" into form and hit import
![Step 2 screenshot](https://images.tango.us/workflows/454b9fc4-c6f9-43d1-beef-9cbbff7f0b0b/steps/48d5b32f-eda8-4e7d-aefc-889d64147192/70a4e414-27b0-4238-9529-f1d64c7fbffd.png?crop=focalpoint&fit=crop&fp-x=0.5000&fp-y=0.2292&fp-z=1.8200&w=1200&mark-w=0.2&mark-pad=0&mark64=aHR0cHM6Ly9pbWFnZXMudGFuZ28udXMvc3RhdGljL21hZGUtd2l0aC10YW5nby13YXRlcm1hcmsucG5n&ar=1812%3A1226)


### c. Select Guess.sol and we're ready
![Step 3 screenshot](https://images.tango.us/workflows/454b9fc4-c6f9-43d1-beef-9cbbff7f0b0b/steps/e35c3f04-92f2-4edf-84ce-ea467feb8797/d401a509-0ab6-4ce8-892d-fbce81c0aeff.png?crop=focalpoint&fit=crop&fp-x=0.1076&fp-y=0.2406&fp-z=2.6507&w=1200&mark-w=0.2&mark-pad=0&mark64=aHR0cHM6Ly9pbWFnZXMudGFuZ28udXMvc3RhdGljL21hZGUtd2l0aC10YW5nby13YXRlcm1hcmsucG5n&ar=1812%3A1226)

## 2. Exploring Guess.sol
This a Solidity file and once we're done, it will contain everything we need for the 'Guessing Game'. Inside the file you'll find the classic ERC20 interface and a simple implementation named ERC20Basic representing USDT which we'll use as means of players entering to play and the prize pool.

There's also an empty contract called Guess inside Guess.sol, and this is what we'll be extending

