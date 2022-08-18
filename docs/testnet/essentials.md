# Essential Testnet Information
An easy-to-read list of essential parameters and configuration settings for Obscuro Testnet.

## Connection to an Obscuro Node
- **RPC http address:** `testnet.obscu.ro:13000`
- **RPC websocket address:** `testnet.obscu.ro:13001`

## Custom Network for MetaMask
- **Network Name:** `Obscuro Testnet`
- **New RPC URL:** `http://127.0.0.1:3000/`
- **Chain ID:** `777`
- **Currency Symbol:** `OBX`

## ObscuroScan
- **URL:** [http://testnet.obscuroscan.io/](http://testnet.obscuroscan.io/)

## Rollup Encryption/Decryption Key
The symmetric key used to encrypt and decrypt transaction blobs in rollups on the Obscuro Testnet:

```
bddbc0d46a0666ce57a466168d99c1830b0c65e052d77188f2cbfc3f6486588c
```

N.B. Decrypting transaction blobs is only possible on testnet, where the rollup encryption key is long-lived and 
well-known. On mainnet, rollups will use rotating keys that are not known to anyone - or anything - other than the 
Obscuro enclaves.

## ERC20 Contracts
We have a couple of "official" testnet ERC20 contracts that are automatically deployed with a static address every time 
testnet is restarted. Both of these can be deposited from and withdrawn to their equivalent ERC20 contracts on the L1.

Please contact us in the discord if you'd like some tokens to test with these contracts (faucet coming soon).

WETH (official ERC20 for wrapped ETH on testnet):

```
0x9802F661d17c65527D7ABB59DAAD5439cb125a67
```

WOBX (official ERC20 for wrapped OBX on testnet):

```
0xf3a8bd422097bFdd9B3519Eaeb533393a1c561aC
```

N.B. We also have a native OBX token on the obscuro network, wrapping/unwrapping is not available yet
