# Deploying a Smart Contract to Obscuro Testnet
Using the steps below you will add an extension to your MetaMask wallet so it can connect to Obscuro Testnet, then using Remix you will deploy your smart contract to Testnet.

## Prerequisites
* MetaMask installed in your browser.
* A local copy of the Obscuro MetaMask wallet extension downloaded from the [Obscuro repo on GitHub](https://github.com/obscuronet/go-obscuro/tree/main/tools/walletextension).

## Prepare Your MetaMask Wallet for Obscuro Testnet
An essential part of how Obscuro provides full privacy is the encryption of communication between an Obscuro application and Obscuro nodes on the network.

Encyption is achieved using the public key of the Trusted Execution Environment (a secure area of the computer's CPU and memory) of the Obscuro Testnet node you are connected to. So you can check the details of a transaction on Obscuro Testnet, a viewing key which is unique to you decrypts the information to make it readable.

The wallet extension should be run locally so that no sensitive data leaves the your machine unencrypted.

1. Open a command line or terminal window on your local machine.
1. Run `walletextension/main/main()` with the following flags to start the wallet extension:

   ```--nodeHost=<Obscuro node's host> --nodePortHTTP=<Obscuro node's HTTP RPC address> --nodePortWS=<Obscuro node's websockets RPC address>```

   The wallet extension is now listening on `http://localhost:3000/`

1. In MetaMask choose to "Add Network"

    1. Network Name = Obscuro Testnet
    1. New RPC URL = http://127.0.0.1:3000/
    1. Chain ID  = 777
    1. Currency Symbol = OBX
    1. Block Explorer URl = _leave blank_

    Requests and responses for Obscuro Testnet will now pass through the local wallet extension in an encrypted format.

1. Browse to `http://localhost:3000/viewingkeys/` to generate a new viewing key unique to you. Sign the viewing key when prompted by MetaMask.

    Communications to the Obscuro Testnet will be now be encrypted with the viewing key unique to you, and decrypted automatically.

    Do remember a new viewing key must be created each time the wallet extension is started.

Now your wallet is configured to make your encrypted traffic viewable you can go ahead and cdeploy your smart contract to the Obscuro Testnet.

## Deploy Your Smart Contract Using an IDE
1. Open your favourite Solidity-compatible Integrated Development Environment. The steps below will assume you are using [Remix](https://github.com/ethereum/remix-ide).

1. Open your solidity smart contract in Remix.

1. Open MetaMask and confirm you are connected to Obscuro Testnet network.

1. In the _Deploy & Run Transactions_ section of Remix change the Environment to _Injected Web3_. This tells Remix to use the network settings currently configured in your MetaMask wallet, which in this case is the Obscuro Testnet.

1. Click the _Deploy_ button to deploy your smart contract to the Obscuro Testnet.