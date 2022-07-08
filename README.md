# Go Obscuro

This repository contains the reference implementation of the [Obscuro Protocol](https://whitepaper.obscu.ro/).

*Note that this is still very much a work in progress, so there are many rough edges and unfinished components.*

## High level overview

The typical blockchain node runs multiple services in a single process. For example:
- P2P Service
- RPC Service
- Data storage
- Transaction execution
- Mempool
- .. 

Obscuro uses Trusted Execution Environments (TEE), like Intel SGX, to execute transactions in a confidential environment, which means we diverge from the typical architecture. 
There are three main components of the architecture, each running as a separate process: the Enclave, the Host and the Wallet Extension.

[//]: # (Todo - @joel - change the diagram to represent the Wallet Extension)
![Architecture](design/obscuro_arch.jpeg)

### I. The Enclave

This is the core component of Obscuro which runs inside the TEE. 
See [go/enclave](go/enclave)

We use [EGo](https://www.edgeless.systems/products/ego/), an open source SDK for developing this confidential component.

The Enclave exposes an [interface](go/common/enclave.go) over RPC which attempts to minimise the Trusted computing base(TCB).

The Enclave component has these main responsibilities:

#### 1. Execute EVM transactions
Obscuro has the goal to be fully compatible with the EVM, so smart contracts can be ported freely from other EVM compatible
chains. To achieve this and minimise the effort and incompatibilities, we depend on [go-ethereum](https://github.com/ethereum/go-ethereum).

The dependency on go-ethereum is not straight forward, since transaction execution is coupled with Ethereum specific consensus rules,
which had to be mocked out.

See [go/enclave/evm](go/enclave/evm)


#### 2. Store the state
The dependency on `go-ethereum` for transaction execution, means that we use the same storage interfaces.

In the current iteration we use [EdglessDB](https://www.edgeless.systems/products/edgelessdb/), an open source database tailor-made 
for confidential computing.

`go-ethereum` uses a key-value store interface, which we implement on top of the SQL database. The reason for this odd choice 
is that the data needs to be indexed and encrypted befored being sent for storage. And the operation of storing data needs 
to be resistant to side-channel analysis which would allow an attacker to infer information on what calculations are being made
based on what data the Enclave is requesting or storing.

In a future iteration, we'll look at alternatives to connect to a performant key-value store designed or modified for confidential computing.

See [go/enclave/db](go/enclave/db)


#### 3. Consume Ethereum blocks 
The Enclave is fed Ethereum blocks through the RPC interface. These blocks are used as the "Source of Truth", and the Enclave 
extracts useful information from them, such as published rollups, deposits to the bridge, etc. Ethereum re-orgs have to be detected
at this level to rollback the Obscuro state accordingly.

To avoid the risk of the Enclave being fed invalid blocks which an attacker can use to probe for information, or to shorten the 
[revelation period](https://whitepaper.obscu.ro/obscuro-whitepaper/detailed-design.html#revelation-mechanism), the blocks have to be checked for validity, which includes checking that enough "work" went into them.
To achieve this we depend on the [Blockchain](https://github.com/ethereum/go-ethereum/blob/e6fa102eb08c2b83ab75e85ca7860eea3a10dab0/core/blockchain.go) 
logic.


#### 4. Bridge to Ethereum 
One of the key aspects of Ethereum Layer 2(L2) solutions is to feature a decentralised bridge that is resistant to 51% attacks.

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

The big advantage of running the mempool inside the secure Enclave is that the ordering of transactions cannot be gamed by the aggregator, 
which makes MEV much more difficult.

See [go/enclave/mempool](go/enclave/mempool)

*Note that the current mempool implementation is very primitive. It always includes all received transactions that were not already
included in a rollup.*


#### 6. The rollups and the PoBI protocol

Like in any blockchain the unit of the protocol is the batch of transactions organized in a chain. 
The Obscuro blocks have an encrypted payload, which is only visible inside the secure Enclave.
All of the logic of maintaining the current state based on incoming data and of producing new rollups is found in the
[go/enclave/rollupchain](go/enclave/rollupchain) package.


#### 7. Cryptography

This is where the Obsuro specific cryptography is implemented. 

- Master seed generation
- Key derivation
- Payload encryption/decryption


See [go/enclave/crypto](go/enclave/crypto)

*Note: The current implementation is still very incipient. There are hardcoded keys, no key derivation, etc.*


### II. The Host

The Host service is the software that is under the control of the operator. It does not run inside a secure Enclave, and there is no attestation on it.

From a threat model point of view, the Host service is seen as an adversary by an Enclave. Any data that it feeds into the Enclave
will be verified and considered malicious.

A secure solution that uses confidential computing will generally try to minimize the TCB, and run as much as possible outside the secure Enclave,
while still achieving the same security goals.

The Host service is the equivalent of a typical blockchain node, and is responsible for:
 
- P2P messaging: Gossiping of encrypted transactions and rollups
- RPC: Exposing an RPC interface similar to the one exposed by normal Ethereum nodes
- Communicating with an Ethereum node for retrieving blocks and for submitting transactions with data that was generated inside the Enclave.
  "This means an Ethereum wallet and the control keys to accounts with enough ETH to publish transactions is required.


See [go/host](go/host)

*Note: The code for host is currently in an incipient phase. The focus of the first phase of development was on the main 
building blocks of the Enclave* 


### III. The Wallet Extension

The missing link to achieving fully private transactions, while allowing end users to continue using their favourite wallets
(like MetaMask) is a very thin component that has to be installed on the user machine.
This component is responsible for the encryption/decryption of the traffic originating from an Obscuro node. It does that 
by generating Viewing Keys behind the scenes.

[//]: # (TODO - Joel: want to add anything here)

See: [tools/walletextension](tools/walletextension)


## Repository Structure
```
|
| - contracts: This folder contains the source code for the solidity Management contract, which will be deployed on Ethereum.  
| - go: The "Golang" implementation of the Obscuro protocol.
|.   |
|.   | - common: This is a somewhat unstructured package containing base types and utils. Note: It will be cleaned up once more patterns emerged.
|.   | - config: A place where the default configurations are found.
|.   | - enclave: This is the component that is loaded up inside SGX.
|.   |.   | - bridge: 
|.   |.   | - core: Data
|.   |.   | - crypto: 
|.   |.   | - db: 
|.   |.   | - enclaverunner: 
|.   |.   | - evm: 
|.   |.   | - main: 
|.   |.   | - mempool: 
|.   |.   | - rollupchain: 
|.   |.   | - rpcencryptionmanager: 
|.   |.   |  
|.   | - ethadapter: 
|.   | - host: The component that performs the networking duties.
|.   | - rpcclientlib: Basic RPC library client code.
|.   | - wallet: 
| 
| - integration: The end to end integration tests.
|
| - testnet: source code for end to end integration tests. 
|.   | - azuredeployer: 
|.   | - contractdeployer: 
|.   | - networkmanager: 
|.   | - obscuroscan: 
|.   | - walletextension: 
| 
| - tools:  

```

## Usage
Todo

