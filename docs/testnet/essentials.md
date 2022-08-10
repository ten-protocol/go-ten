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

