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

The Enclave exposes an [interface](go/common/enclave.go) over RPC which attempts to minimise the "trusted computing base"(TCB).

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
is that the data needs to be indexed and encrypted before being sent for storage. And the operation of storing data needs 
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
One of the key aspects of Ethereum Layer 2 (L2) solutions is to feature a decentralised bridge that is resistant to 51% attacks.

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
See [go/enclave/crypto](go/enclave/crypto)

These are the main components:

##### a) RPC encryption/decryption

The Obscuro "Wallet extension" encrypts all requests from users (transactions or smart contract method calls) with the "Obscuro public key", which is a key derived from the 
master seed. 
The response is in turn encrypted with the "Viewing Key" of the requesting user.

This component manages viewing keys and handles the encryption and decryption.

The transactions received from users are gossiped with the other aggregators encrypted with the "Obscuro Public Key".

*Note: In the current implementation, this key is hardcoded.*

See: [go/enclave/rpcencryptionmanager](go/enclave/rpcencryptionmanager)

##### b) Transaction blob encryption

Transactions are stored as `calldata` blobs in ethereum transactions.
The component that creates these blobs has to encrypt them with derived keys according to their revelation period.

*Note: In the current implementation the payload is encrypted with a single hardcoded key*

See: [go/enclave/crypto/transaction_blob_crypto.go](go/enclave/crypto/transaction_blob_crypto.go)


#### 8. RPC

The enclave exposes an RPC interface generated with [proto-buf](https://developers.google.com/protocol-buffers).

The interface is described in [enclave.proto](go/common/rpc/generated/enclave.proto).

See [go/common/rpc](go/common/rpc)

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
  This means an Ethereum wallet and the control keys to accounts with enough ETH to publish transactions is required.


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

<pre>
root
├── <a href="./contracts">contracts</a>: : Solidity Management contract, which will be deployed on Ethereum
├── <a href="./go">go</a>
│   ├── <a href="./go/common">common</a>: Unstructured package containing base types and utils. Note: It will be cleaned up once more patterns emerged.
│   ├── <a href="./go/config">config</a>: A place where the default configurations are found.
│   ├── <a href="./go/enclave">enclave</a>: The component that is loaded up inside SGX.
│   │   ├── <a href="./go/enclave/bridge">bridge</a>: The platform side of the decentralised bridge logic.
│   │   ├── <a href="./go/enclave/core">core</a>: Base data structures used only inside the enclave. 
│   │   ├── <a href="./go/enclave/crypto">crypto</a>: Implementation of the Obscuro cryptography.
│   │   ├── <a href="./go/enclave/db">db</a>: The database implementations. 
│   │   ├── <a href="./go/enclave/enclaverunner">enclaverunner</a>: The entry point to the standalone enclave process. 
│   │   ├── <a href="./go/enclave/evm">evm</a>: Obscuro transaction execution on top of the EVM.
│   │   ├── <a href="./go/enclave/main">main</a>: Main
│   │   ├── <a href="./go/enclave/mempool">mempool</a>: The mempool living inside the enclave
│   │   ├── <a href="./go/enclave/rollupchain">rollupchain</a>: The main logic for calculating the state and the POBI protocol.
│   │   └── <a href="./go/enclave/rpcencryptionmanager">rpcencryptionmanager</a>: Responsible for encrypting the communication with the wallet extension.
│   ├── <a href="./go/ethadapter">ethadapter</a>: Responsible for interpreting L1 transactions 
│   │   ├── <a href="./go/ethadapter/erc20contractlib">erc20contractlib</a>: Understand ERC20 transactions.
│   │   └── <a href="./go/ethadapter/mgmtcontractlib">mgmtcontractlib</a>: Understand Obscuro Management contrract transactions. 
│   ├── <a href="./go/host">host</a>: The standalone host process 
│   │   ├── <a href="./go/host/hostrunner">hostrunner</a>: The entry point.
│   │   └── <a href="./go/host/p2p">p2p</a>: The P2P communication implementation. 
│   ├── <a href="./go/rpcclientlib">rpcclientlib</a>: Library to allow go applications to connect to a host via RPC.
│   └── <a href="./go/wallet">wallet</a>: Logic around wallets. Used both by the node, which is an ethereum wallet, and by the tests
├── <a href="./integration">integration</a>: Integration tests that spin up Obscuro networks.
│   ├── <a href="./integration/simulation">simulation</a>: A series of tests that simulate running networks with different setups.
├── <a href="./testnet">testnet</a>: Utilities for deploying a testnet.
└── <a href="./tools">tools</a>: Peripheral tooling. 
│   ├── <a href="./tools/azuredeployer">azuredeployer</a>: Help with deploying obscuro nodes on SGX enabled azure VMs.
│   ├── <a href="./tools/contractdeployer">contractdeployer</a>: todo - Joel 
│   ├── <a href="./tools/networkmanager">networkmanager</a>: todo - Joel
│   ├── <a href="./tools/obscuroscan">obscuroscan</a>: todo - Joel
│   └── <a href="./tools/walletextension">walletextension</a>: todo - Joel

</pre>


## Testing

The Obscuro integration tests are found in: [integration/simulation](integration/simulation).

The main tests are "simulations", which means they spin up both an L1 network and an L2 network, and then inject random transactions.
Due to the non-determinism of both the "mining" protocol in the L1 network and the nondeterminism of POBI, coupled with the random traffic,
it allows the tests to capture many corner cases without having to explicitly write individual tests for them. 

The first [simulation_in_mem_test](integration/simulation/simulation_in_mem_test.go) runs fully in one single process on top of a 
mocked L1 network and with the networking components of the Obscuro node swapped out, and is just focused on producing 
random L1 blocks at very short intervals.  The [ethereummock](integration/ethereummock) implementation is based on the ethereum protocol with the individual nodes 
gossiping with each other with random latencies, producing blocks at a random interval distributed 
around a configured ``AvgBlockDuration``, and making decisions about the canonical head based on the longest chain.
The L2 nodes are each connected to one of these mocked L1 nodes, and receive a slightly different view.
If this test is run long enough, it verifies the POBI protocol.

There are a number of simulations that gradually become more realistic, but at the cost of a reduction in the number of 
blocks that can be generated.

The [simulation_geth_in_mem_test](integration/simulation/simulation_geth_in_mem_test.go) replaces the mocked ethereum nodes with a 
network of geth nodes started in clique mode. The lowest unit of time of producing blocks in that mode is `1 second`.

The [simulation_full_network_test](integration/simulation/simulation_full_network_test.go) starts standalone local processes for
both the enclave and the obscuro node connected to real geth nodes.

The [simulation_docker_test](integration/simulation/simulation_docker_test.go) goes a step further and runs the enclave in "Simulation mode" 
in a docker container with the "EGo" library. 

The [simulation_azure_enclaves_test](integration/simulation/simulation_azure_enclaves_test.go) is the ultimate test where the enclaves are deployed 
in "Real mode" on SGX enabled VMs on Azure.


A [transaction injector](integration/simulation/transaction_injector.go) is able to create and inject random transactions in any 
of these setups by receiving RPC handles to the nodes.

## Usage
Todo

### Compiling

- Install go version > 1.8

### Compiling


## Community 

Discussions around development 