---
---
# Using the Testnet Token Faucet
Using the steps below you will request testnet OBX tokens from the faucet server.

## Prerequisites
* Your wallet address.
* Access to a command shell with the curl cli installed.

## Requesting Testnet OBX Tokens
1. Make a note of your wallet address or copy it to your clipboard.
2. Open a command shell and issue the below command, where `<address>` should be replaced with the value stored in your clipboard (e.g. `0x75Ad715443e1E2EBdaFA33ABB3B08443966019A6`). The faucet server will credit 100,000 OBX by default.

```bash
curl --location --request POST 'http://testnet-faucet.uksouth.azurecontainer.io/fund/obx' --header 'Content-Type: application/json' --data-raw '{ "address":"<your address>" }'
```

3. After a short period of time the curl command will return `{"status":"ok"}` confirming OBX tokens have been credited to your wallet.

## Viewing Your Wallet Balance
To view the balance of your wallet you will need to establish a connection from your wallet to the Obscuro Testnet. An essential part of how Obscuro provides full privacy is the encryption of communication between an Obscuro application and Obscuro nodes on the network. As a result, you will need to use the wallet extension to allow your wallet to communication with the Obscuro Testnet.

Use the steps [here](https://docs.obscu.ro/testnet/deploying-a-smart-contract/#prepare-your-metamask-wallet-for-obscuro-testnet) to prepare your MetaMask wallet for Obscuro Testnet.
