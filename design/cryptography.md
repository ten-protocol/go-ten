# Obscuro Cryptography


Privacy in Obscuro is achieved by:
1. using local databases that run inside TEEs and write only encrypted data to disk
2. the ledger (list of user transactions) is stored encrypted as rollups in Ethereum
3. users submitting transactions encrypted with a well known key (the Obscuro Key)


The first element can be achieved by using existing solutions like [edgelessDB](https://www.edgeless.systems/products/edgelessdb/).

To minimise the complexity for the second and third, we propose to start with a securely generated piece of entropy called: "The Master Seed"(MS),
and then use that to deterministically derive the keys necessary to achieve the other goals.

An additional complexity is that Obscuro has the requirement to support temporary privacy, and reveal transactions after a period of time has elapsed.


## Scope

1. Master Seed Generation
2. Enclave Encrypted Communication
3. Transaction Encryption per Rollup with Revelation Period
4. Reveal the rollup transactions

## 1. Master Seed Generation
- Master Seed is a strongly generated entropy data that serves as the encryption base of the network
- Use the random number generation available inside the TEEs
- Generation happens when the central sequencer bootstraps the network
- The central sequencer generates the Master seed and initializes the Management Contract
- Upon successful attestation the Central Sequencer uses the attested enclave public key to share the Master Seed


## 2. Ensure Enclave Encrypted Communication
- Client-Enclave communication is encrypted using the Obscuro Key
- Clients encrypt using the Obscuro Key - Public Key
- Enclave responds encrypting with the Obscuro Key - Private Key
- The Obscuro Key is a Key-Pair derived from the Master Seed + Genesis Rollup
- The Public key will be published to the Management contract and used by all obscuro tooling ( like the wallet extension ) to encrypt transactions
- The Private key will be only available inside the enclave through key derivation
- Enclaves will determine the Private key when deriving the Master Seed + Genesis Rollup 
  - Enclaves have the Master Seed through the attestation process
  - Enclaves fetch Genesis Rollup through the Layer 1 blocks

- HKDF (HMAC-based Key Derivation Function) is used to derive keys
  - Given the high entropy of the Master Seed no need for PBKDF Family key stretching
  - Derivations use public rollup metadata such as height ( or some other shared field ) in the Info component of HKDF and use a fixed size salt as per https://soatok.blog/2021/11/17/understanding-hkdf/


## 3. Transaction Encryption per Rollup with Revelation Period
- Transactions are encrypted in Obscuro to provide confidentiality
- To avoid base-key compromise, transaction encryption keys are derived twice
  - Each rollup has a Rollup Encryption Key derived from the Master Seed + Rollup
  - Each transaction is encrypted with a Revelation Period Key that is derived from the Rollup Encryption Key + Revelation Period
- HKDF (HMAC-based Key Derivation Function) is used to derive keys
  - Given the high entropy of the Master Seed no need for PBKDF Family key stretching
  - Derivations use public rollup metadata such as height ( or some other shared field ) in the Info component of HKDF and use a fixed size salt as per https://soatok.blog/2021/11/17/understanding-hkdf/
- EVM-Compatible Transactions using [AccessList](https://eips.ethereum.org/EIPS/eip-2930) property are able to specify the desired revelation period
  - Using `AccessList` Address _Null Address_ combined with Storage Key _hexadecimal 1-5_ determines the revelation period
  - Transactions without Revelation Period specified are encrypted by default using the smallest unit of time revelation period


The [specification](https://eips.ethereum.org/EIPS/eip-2930) of the `AccessList` Field is as follows:

```
For the transaction to be valid, accessList must be of type [[{20 bytes}, [{32 bytes}...]]...], 
where ... means “zero or more of the thing to the left”.
For example, the following is a valid access list (all hex strings would in reality be in byte representation):

[
    [
        "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",
        [
            "0x0000000000000000000000000000000000000000000000000000000000000003",
            "0x0000000000000000000000000000000000000000000000000000000000000007"
        ]
    ],
    [
        "0xbb9bc244d798123fde783fcc1c72d3bb8c189413",
        []
    ]
]
```

Obscuro uses the `AccessList` object to determine the revelation period.
A special address and special storage keys combination is checked for internally.

To simplify the protocol the Null address (0x0000000000000000000000000000000000000000) will be used and the hexadecimal representation of 1 to 5 is used to specify the revelation period.

An example would be:
```
[
    [
        "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",                             <- Canonical AccessList usage
        [
            "0x0000000000000000000000000000000000000000000000000000000000000003",
            "0x0000000000000000000000000000000000000000000000000000000000000007"
        ]
    ],
    [
        "0x0000000000000000000000000000000000000000",                             <- Obscuro Revelation Period
         [
            "0x0000000000000000000000000000000000000000000000000000000000000003", <- 1 day revelation Period
            "0x0000000000000000000000000000000000000000000000000000000000000007"  <- ignored
        ]
    ]
]
```


### 4. Reveal the rollup transactions
- L1 blocks are used as a clock mechanism
  - There are 5 revelation periods
    * XS - 12 seconds (1 block)
    * S - 1 hour (3600/12 = 300 blocks)
    * M - 1 day (24 * 3600/12 = 7200 blocks)
    * L - 1 month (30 * 24 * 3600/12 = 216000 blocks)
    * XL - 1 year (12 * 30 * 24 * 3600/12 = 2592000 blocks)
- Validators do not reveal keys so that keys are not susceptible to attacks
- Central Sequencer stores the keys in a rollup, decrypt time, key tuple in a database.
- When a Light Batch is created the revealed keys are appended to it and removed from the database.


### Problems

1. Symmetric vs Asymmetric encryption for rollup tx encryption.
- symmetric is more space efficient
- AES 256 is the battle tested standard
- AES-256-GCM allows for stream encryption ( no block padding + more efficient ) and authentication (check for tampering before decrypting)

2.  How does a developer specify the reveal period?
- Transactions use the AccessList field to specify for how long they should be encrypted

3. Ensure a node operator can't fast-forward the clock.
- Validators do not release keys, only the centralized sequencer releases them

4. In the event of a db failure how does the central sequencer know which keys to recalculate
- Traversing the existing chain and recalculating keys might be too expensive
- Sequencer can reveal keys on-demand with recalculating the whole chain

5. Do validators have a validation period ? How can a validation period be enforced if they are given a non-mutable Master Seed ?
6. Do validators need the Obscuro Private Key ?