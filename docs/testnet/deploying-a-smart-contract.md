# Deploying a Smart Contract to Obscuro Testnet
Using the steps below you will add an extension to your MetaMask wallet so it can connect to Obscuro Testnet, then using Remix you will deploy your smart contract to Testnet. Remember there are no gas fees on Obscuro Testnet so you do not need to load your account with OBX tokens in advance.

## Prerequisites
* [MetaMask](https://metamask.io/) wallet installed in your browser.
* A local copy of the [Obscuro MetaMask wallet extension](../wallet-extension/wallet-extension.md)

## Prepare Your MetaMask Wallet for Obscuro Testnet
An essential part of how Obscuro provides full privacy is the encryption of communication between an Obscuro application and Obscuro nodes on the network.

Follow the steps [here](../wallet-extension/wallet-extension.md) to configure and start the wallet extension and 
generate a viewing key. The wallet extension should be started with the following flags:

```
--nodeHost=testnet.obscu.ro --nodePortHTTP=13000 --nodePortWS=13001
```

If you do not have the Obscuro wallet extension running, MetaMask will not be able to communicate with the Obscuro 
Testnet.

> **_TIP_**  Every time you restart the wallet extension, you must generate a new viewing key. This is because the 
  private-key part of the viewing key is only held in memory and never persisted to disk, for security reasons.

Your wallet is now configured for the Obscuro Testnet which allows you to view encrypted traffic for your wallet only. 
You can now go ahead and deploy your smart contract to the Obscuro Testnet.

## Deploy Your Smart Contract Using an IDE
1. Browse to the popular Solidity-compatible Integrated Development Environment called [Remix](https://remix.ethereum.org/).

    > **_TIP_**  Take a look at the [Remix docs](https://remix-ide.readthedocs.io/en/latest/create_deploy.html) for guidance on using the IDE.

1. Load your solidity smart contract into the Remix workspace.

1. Compile your smart contract using the Remix Solidity Compiler.

1. Log in to MetaMask and confirm you are connected to Obscuro Testnet network. The parameters for the Obscuro Testnet can be found [here](./essentials.md).

1. In the _Deploy & Run Transactions_ section of Remix change the Environment to _Injected Web3_. This tells Remix to use the network settings currently configured in your MetaMask wallet, which in this case is the Obscuro Testnet. If the connection to Obscuro Testnet is successful you will see the text _Custom (777) network_ displayed under _Injected Web3_.

1. Click the _Deploy_ button to deploy your smart contract to the Obscuro Testnet.

1. MetaMask will automatically open and ask you to confirm the transaction. Click _Confirm_.

1. Wait for your transaction to be confirmed. Notifications will be shown in MetaMask and Remix to indicate that the transaction was successful.

Congratulations, your smart contract is now deployed to Obscuro Testnet!

> **Be prepared to redeploy your smart contract to Obscuro Testnet often. The Testnet is likely to be restarted without warning!** As Testnet matures and new features are added it will become more and more stable, but in the early days do expect some turbulence.

Because Obscuro provides full privacy, the details of your transaction are encrypted and only visible to you, as the holder of your wallet's private key.

Now head over to the [ObscuroScan page](./obscuroscan.md) to see how you can view the transaction details.
