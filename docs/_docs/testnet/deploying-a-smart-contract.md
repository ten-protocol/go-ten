---
---
# Deploying a Smart Contract to Ten Testnet
Using the steps below you will add an extension to your MetaMask wallet so it can connect to Ten Testnet, then using Remix you will deploy your smart contract to Testnet. Remember there are no gas fees on Ten Testnet so you do not need to load your account with TEN tokens in advance.

## Prerequisites
* [MetaMask](https://metamask.io/) wallet installed in your browser.
* A local copy of the [Ten MetaMask wallet extension](https://docs.obscu.ro/wallet-extension/wallet-extension/)

## Prepare Your MetaMask Wallet for Ten Testnet
An essential part of how Ten provides full privacy is the encryption of communication between an Ten application and Ten nodes on the network.

Follow the steps [here](https://docs.obscu.ro/wallet-extension/wallet-extension/) to configure and start the wallet extension and 
generate a viewing key. 


If you do not have the Ten wallet extension running, MetaMask will not be able to communicate with the Ten 
Testnet.

> **_TIP_**  Every time you restart the wallet extension, you must generate a new viewing key. This is because the 
  private-key part of the viewing key is only held in memory and never persisted to disk, for security reasons.

Your wallet is now configured for the Ten Testnet which allows you to view encrypted traffic for your wallet only. 
You can now go ahead and deploy your smart contract to the Ten Testnet.

## Deploy Your Smart Contract Using an IDE
1. Browse to the popular Solidity-compatible Integrated Development Environment called [Remix](https://remix.ethereum.org/).

    > **_TIP_**  Take a look at the [Remix docs](https://remix-ide.readthedocs.io/en/latest/create_deploy.html) for guidance on using the IDE.

1. Load your solidity smart contract into the Remix workspace.

1. Compile your smart contract using the Remix Solidity Compiler.

1. Log in to MetaMask and confirm you are connected to Ten Testnet network. The parameters for the Ten Testnet can be found [here](https://docs.obscu.ro/testnet/essentials/).

1. In the _Deploy & Run Transactions_ section of Remix change the Environment to _Injected Web3_. This tells Remix to use the network settings currently configured in your MetaMask wallet, which in this case is the Ten Testnet. If the connection to Ten Testnet is successful you will see the text _Custom (443) network_ displayed under _Injected Web3_.

1. Click the _Deploy_ button to deploy your smart contract to the Ten Testnet.

1. MetaMask will automatically open and ask you to confirm the transaction. Click _Confirm_.

1. Wait for your transaction to be confirmed. Notifications will be shown in MetaMask and Remix to indicate that the transaction was successful.

Congratulations, your smart contract is now deployed to Ten Testnet!

> **Be prepared to redeploy your smart contract to Ten Testnet often. The Testnet is likely to be restarted without warning!** As Testnet matures and new features are added it will become more and more stable, but in the early days do expect some turbulence.

Because Ten provides full privacy, the details of your transaction are encrypted and only visible to you, as the holder of your wallet's private key.

Now head over to the [TenScan page](https://docs.obscu.ro/testnet/obscuroscan/) to see how you can view the transaction details.
