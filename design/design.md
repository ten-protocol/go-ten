# Obscuro design

## Scope

The purpose of this document is to describe aspects of Obscuro's technical design that are not addressed in the 
[Obscuro whitepaper](https://whitepaper.obscu.ro/).

## Overview

The following diagram shows the key components of an Obscuro deployment:

![architecture diagram](./obscuro_arch.jpeg)

The Ethereum node and Ethereum chain components shown in this diagram are developed and maintained by third-parties. 
The following additional components must be developed:

* **The enclave:** The trusted part of the Obscuro node that runs inside a trusted execution environment (TEE)
* **The host:** The remainder of the Obscuro node that runs outside the TEE
* **The Obscuro management contract:** The Ethereum mainnet contracts required by the Obscuro protocol, described 
  [here](https://whitepaper.obscu.ro/obscuro-whitepaper/l1-contracts)
* **Client apps:** Applications that interact with the Obscuro node (e.g. Obscuro wallets)

## Host/enclave split

The node is divided into two components, the host and the enclave. Wherever reasonable, node logic should be part of 
the host rather than the enclave. This has two benefits:

* It minimises the amount of code in the 
  [trusted computing base (TCB)](https://en.wikipedia.org/wiki/Trusted_computing_base)
* It reduces churn in the TCB, reducing the frequency of re-attestations

The host and the enclave are two separate OS processes, rather than separate threads in a single process. This is 
because our initial target TEE, [Intel SGX](https://en.wikipedia.org/wiki/Software_Guard_Extensions), requires the 
TEE to execute as a separate process.

The host and the enclave communicate via RPC, using the [gRPC](https://grpc.io/) library. gRPC was selected as it is
open-source (Apache 2.0) and has broad adoption.

For simplicity, this transport is not authenticated (e.g. using TLS or credentials). One possible attack vector is for
a _parasite_ aggregator to only run the host software, and connect to another aggregator's enclave to submit
transactions, in order to economise on operating costs. To avoid this scenario, the enclave is designed to have full
control over which account receives the rollup rewards, meaning that a would-be parasite aggregator does not receive
any rewards for acting in this manner.

To reduce coupling, the enclave process will be monitored and managed by a supervisor, and not by the host process.

## Enclave datastore

The enclave is backed by a datastore. This datastore stores seven maps:

1. L1 block hash -> the state after ingesting the block
2. L1 block hash -> the corresponding block and its height
3. L2 transaction hash -> the corresponding mempool transaction
4. Rollup hash -> the state after adding the rollup
5. Rollup hash -> the corresponding rollup
6. Rollup hash -> the transactions in the corresponding rollup
7. int -> the proposed rollups with height <int>

The files backing the datastore are stored outside the enclave. To ensure the datastore contents remain confidential, 
the values in the datastore maps are stored in encrypted form [WITH WHAT KEY?]. The hash/int keys are not considered 
sensitive and are not encrypted.

The datastore used is LevelDB. We use LevelDB because:

* It is used by Go Ethereum. This shows that it is suitable for workloads of this type. It also means that we already 
  have an indirect dependency on LevelDB
* It is a library, meaning that we do not introduce one or more additional datastore components to manage
  * An additional benefit is that since it is not exposed to the user as an additional component to manage, we can 
    switch it out later at lower cost
* It has reasonable adoption and is well-maintained
* There is a Go implementation (https://github.com/syndtr/goleveldb)

### Can enclave data be ephemeral?

An alternative approach to enclave data would be to make it ephemeral. State could be recovered through a combination 
of requesting data from peer nodes (e.g. blocks, transactions), requesting resubmission from clients (e.g. mempool 
transactions) and recreating the data (e.g. candidate rollups).

However, this has several downsides:

* It is less efficient
* Is provides a worse user experience (transactions have to be resubmitted on occasion)
* It forces the enclave to hold the full state inside the enclave memory
