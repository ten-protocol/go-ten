---
---
# Using the Testnet Token Faucet
Using the steps below you will request testnet OBX from the faucet available on the Obscuro Discord server.

## Prerequisites
* Access to the [Obscuro Discord server](https://discord.gg/yQfmKeNzNd).
* (Optional) [MetaMask](https://metamask.io/) wallet installed in your browser.
* (Optional) A local copy of the [Obscuro MetaMask wallet extension](https://docs.obscu.ro/wallet-extension/wallet-extension/).

## Requesting Testnet OBX
1. Make a note of your wallet address or copy it to your clipboard.
2. Open the [_faucet-requests_ channel](https://discord.gg/5qyj3qraaH) on Obscuro Discord.
3. Request OBX using the `/faucet` command. The faucet will credit 100,000 OBX by default:
   ![faucet command](../../assets/images/faucet-cmd.png)
4. Provide your wallet address and hit Enter. The faucet will acknowledge your request:
   ![faucet ack](../../assets/images/faucet-ack.png)
5. After a short period of time the faucet will confirm the Testnet OBX have been credited to your wallet:
   ![faucet complete](../../assets/images/faucet-done.png)

## Viewing Your Wallet Balance
To view the balance of your wallet you will need to establish a connection from your wallet to the Obscuro Testnet. An essential part of how Obscuro provides full privacy is the encryption of communication between an Obscuro application and Obscuro nodes on the network. As a result, you will need to use the wallet extension to allow your wallet to communication with the Obscuro Testnet.

Use the steps [here](https://docs.obscu.ro/testnet/deploying-a-smart-contract/#prepare-your-metamask-wallet-for-obscuro-testnet) to prepare your MetaMask wallet for Obscuro Testnet.

## Requesting Testnet OBX directly
In the event that you do not have access to Discord, or the faucet bot is not working, you can request OBX directly from 
the faucet server using the below; 

1. Make a note of your wallet address or copy it to your clipboard.
2. Open a command shell and issue the below command, where `<address>` should be replaced with the value stored in your clipboard (e.g. `0x75Ad715443e1E2EBdaFA33ABB3B08443966019A6`). The faucet server will credit 100,000 OBX by default.
```bash
curl --location --request POST 'http://testnet-faucet.uksouth.azurecontainer.io/fund/eth' --header 'Content-Type: application/json' --data-raw '{ "address":"<your address>" }'
```
3. After a short period of time the curl command will return `{"status":"ok"}` confirming OBX have been credited to your wallet.
