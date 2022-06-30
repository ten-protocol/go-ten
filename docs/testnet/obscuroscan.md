# ObscuroScan
[ObscuroScan](http://obscuroscan-01.uksouth.azurecontainer.io/) is Obscuro's equivalent of Etherscan. ObscuroScan is a blockchain explorer for the Obscuro Testnet. At the moment it is rudimentary and over time it will be improved to provide more functionality and look better.

ObscuroScan allows you to decrypt rollup transactions blobs on Testnet. You can also monitor in realtime the L1 blocks and the Obscuro rollups via an Obscuro node connected to the Testnet.

## How to Decrypt Transaction Blobs
Decrypting transaction blobs is only possible on Testnet to help you understand how Obscuro works. ObscuroScan on Testnet uses a rollup encryption key which is long-lived and well-known. On Mainnet, rollups will use rotating keys that are not known to anyone, or anything, other than the Obscuro enclaves.

From the [ObscuroScan landing page](http://obscuroscan-01.uksouth.azurecontainer.io/) click _Decrypt rollup transaction blobs (testnet only!)_ to go to the decryption page.

## How to Monitor L1 Blocks and Obscuro Rollups
From the [ObscuroScan landing page](http://obscuroscan-01.uksouth.azurecontainer.io/) click _Monitor L1 blocks and Obscuro rollups via the connected Obscuro node_ to go to the monitoring page.

You will see in realtime the current head L1 block and the current Obscuro rollup.