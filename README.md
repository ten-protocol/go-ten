# Go Obscuro

This repository contains the reference implementation of the [Obscuro Protocol](https://whitepaper.obscu.ro/).


## High level overview

The typical blockchain node runs multiple services in a single process. For example:
- P2P Service
- RPC Service
- Data storage
- Transaction execution
- Mempool

Obscuro uses TEEs (Intel SGX) to execute transactions in a confidential environment, which means we diverge from the typical architecture. 
There are three main components of the architecture, each running as a separate process: the enclave, the host and the wallet extension.

[//]: # (Todo - @joel - change the diagram to represent the wallet extension)
![Architechture](design/obscuro_arch.jpeg)

### I. The Enclave

This is the core component of Obscuro which runs inside the TEE. 
See [go/enclave](go/enclave)

We use [EGo](https://www.edgeless.systems/products/ego/), an open source SDK for developing this confidential component.

The Enclave exposes an [interface](go/common/enclave.go) over RPC which attempts to minimise the TCB (Trusted computing base).

The enclave component has these main responsibilities:

#### 1. Execute EVM transactions
Obscuro has the goal to be fully compatible with the EVM, so smart contracts can be ported freely from other EVM compatible
chains. To achieve this and minimise the effort and incompatibilities, we depend on [go-ethereum](https://github.com/ethereum/go-ethereum).

The dependency on go-ethereum is not straight forward, since transaction execution is coupled with ethereum specific consensus rules,
which had to be mocked out.

See [go/enclave/evm](go/enclave/evm)


#### 2. Store the state
The dependency on go-ethereum for transaction execution, means that we have to use the same storage interfaces.

In the current iteration we use [EdglessDB](https://www.edgeless.systems/products/edgelessdb/), an open source database tailor-made 
for confidential computing.

Go-ethereum uses a key-value store interface, which we implemented on top of the SQL database. The reason for this odd choice 
is that the data needs to be indexed and encrypted befored being sent for storage. And the operation of storing data needs 
to be resistant to side-channel analysis which would allow an attacker to infer information on what calculations are being made
based on what data the enclave is requesting or storing.

In a future iteration, we'll look at alternatives to connect to a performant k/v store designed or modified for confidential computing.

See [go/enclave/db](go/enclave/db)


#### 3. Consume Ethereum blocks 
The enclave is fed ethereum blocks through the rpc interface. These blocks are used as the "Source of Truth", and the enclave 
extracts useful information from them, such as published rollups, deposits to the bridge, etc. Ethereum Re-orgs have to be detected
at this level to rollback the Obscuro state accordingly.

To avoid the risk of the enclave being fed invalid blocks which an attacker can use to probe for information, or to shorten the 
revelation period, the blocks have to be checked for validity, which includes checking that enough "work" went into them.
To achieve this we depend on the [Blockchain](https://github.com/ethereum/go-ethereum/blob/e6fa102eb08c2b83ab75e85ca7860eea3a10dab0/core/blockchain.go) 
logic.


#### 4. Bridge to Ethereum 
Obscuro is an Ethereum Layer 2, and one of the key aspects of layer2s is to feature a re-org resistent decentralised bridge.

Obscuro features a L2 side of the bridge that is completely under the control of the platform.

During processing of the Ethereum blocks, the platform generates synthetic L2 transactions based on every relevant transaction found there.
For example when Alice deposits 10ABC from her account to the L1 bridge, Obscuro will execute a synthetic L2 transaction (that it deterministically
generated from the L1 transaction), which moves 10WABC from the L2 bridge to Alice's address on Obscuro. 

This logic is part of the consensus of Obscuro, every node receiving the same block containing the rollup and the deposits, will generate the exact same synthetic transaction.

See [go/enclave/bridge](go/enclave/bridge)

*Note that the current bridge implementation is very primitive and only features two supported hardcoded ERC20 tokens to demonstrate
the mechanics.*


#### 5. Mempool 

The mempool is the component which handles the incoming transactions and is responsible for selecting which transactions 
to include in the current batch and pick the order.

The big advantage of running the mempool inside the secure enclave is that the ordering of transactions cannot be gamed by the aggregator, 
which makes MEV much more difficult.

See [go/enclave/mempool](go/enclave/mempool)

*Note that the current mempool implementation is very primitive. It always includes all received transactions that were not already
included in a rollup.*


#### 6. The rollups and the PoBI protocol

Like in any blockchain the unit of the protocol is the batch of transactions organized in a chain. 
The Obscuro blocks have an encrypted payload, which is only visible inside the secure enclave.
All the logic of maintaining the current state based on incoming data and of producing new rollups is found in the
[go/enclave/rollupchain](go/enclave/rollupchain) package.


#### 7. Cryptography

This is where the Obsuro specific cryptography is implemented. 

- Master seed generation
- Key derivation
- Payload encryption/decryption


See [go/enclave/crypto](go/enclave/crypto)

*Note: The current implementation is still very incipient. There are hardcoded keys, no key derivation, etc.*


### II. The Host

The host service is the software that is under the control of the operator. It does not run inside a secure enclave, and there is no attestation on it.

From a threat model point of view, the host service is seen as an adversary by an enclave. Any data that it feeds into the enclave
will be verified and considered malicious.

A secure solution that uses confidential computing will generally try to minimize the TCB, and run as much as possible outside the secure enclave,
while still achieving the same security goals.

The host service is the equivalent of a typical blockchain node, and is responsible for:
 
- P2P messaging: Gossiping of encrypted transactions and rollups
- RPC: Exposing an RPC interface similar to the one exposed by normal ethereum nodes
- Communicating with an Ethereum node for retrieving blocks and for submitting transactions with data that was generated inside the enclave. 
This means it has to be an ethereum wallet and control keys to accounts with enough Eth to publish transactions.  

See [go/host](go/host)

Note: 

### III. The Wallet Extension

---

Besides the main components, there are a number of tools ..



## Repository Structure
```
|
| - contracts: This folder contains the source code for the solidity Management contract, which will be deployed on Ethereum.  
|               For now it's a very basic version that just stores rollups.
| - go: The obscuro protocol code.
|.   |
|.   | - obscuronode: source code for the obscuro node. Note that the node is composed of two executables: "enclave" and "host".
|.   |      | 
|.   |      | - enclave: This is the component that is loaded up inside SGX.
|.   |      | - host: The component that performs the networking duties.
|.   |      | - obscuroclient: basic RPC library client code.
|.   |      | - nodecommon:
|
| 
| - integration: source code for end to end integration tests. 
|

```

## Usage
Todo

## High level overview
Todo


