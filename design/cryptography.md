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

1. Generating the Master seed
2. Derive the Obscuro Key
3. Encrypting transactions such that they can be revealed
4. Reveal the rollup transactions

## 1. Generating the Master seed
- Use the random number generation available inside the TEEs
- Generation happens when the central sequencer bootstraps the network


## 2. Derive the Obscuro Key
- The Obscuro Key is a Key-Pair derived from the Master Seed + Genesis Rollup
- The Public key will be published to the Management contract and used by all obscuro tooling ( like the wallet extension ) to encrypt transactions
- The Private key will be determined by other enclaves by deriving the Master Seed + Genesis Rollup 
  - The Master Seed will be shared using the attestation process
  - The Genesis Rollup will be shared when reading the L1 blocks
- HKDF (HMAC-based Key Derivation Function) is used to derive keys
  - Given the high entropy of the Master Seed no need for PBKDF Family key stretching
  - Derivations use public rollup metadata such as height ( or some other shared field ) in the Info component of HKDF and use a fixed size salt as per https://soatok.blog/2021/11/17/understanding-hkdf/


## 3. Encrypting transactions such that they can be revealed
- Transactions are encrypted by default using the smallest unit of time revelation period
- Transactions using AccessList can specify the desired revelation period
- Using Null Address combined with hexadecimal 1-5 level determines the revelation period
- Obscuro interprets the revelation period and encrypts the transaction with the corresponding revelation period key

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
- If the database is lost for some reason, the Central Sequencer can recalculate the keys on-demand.


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
- Traversing the existing chain might be too expensive
- Sequencer can provide and on-demand key reveal service 
  - Backed by enough horizontal scalability
  - Providing the rollup height and revelation period any service can easily determine if the key is to be released

