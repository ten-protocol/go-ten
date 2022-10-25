# Fast finality design

## Scope

The introduction of a single sequencer to improve cost-per-transaction during the bootstrapping phase, as 
described in the [Bootstrapping Strategy design doc](./Bootstrapping_strategy.md).

## Requirements

* Finality
  * Transaction *soft* finality (finality guaranteed by the sequencer) is achieved in less than two seconds
  * There is an eventual transaction *hard* finality (finality guaranteed by the L1)
* Cost
  * L1 transaction costs can be driven arbitrarily low, at the expense of extending the hard-finality window
* User/dev experience
  * The responses to RPC calls reflect the soft-finalised transactions, and not just the hard-finalised transactions
* Operations
  * The sequencer is highly available
* Security
  * The sequencer is not able to "rewrite history" (or is strongly disincentivised from doing so), even for soft-final 
    transactions (e.g. to perform front-running)

## Assumptions

* There is a single aggregator that is also the genesis node. This node is known as the *sequencer*
* The enclave has a start-up delay that exceeds the production time for light batches (see below)

## Constraints

* It must be possible to only publish rollups to the L1 every `x` blocks, where `x` >> 1
* It must be possible to produce light batches (see below) at a higher frequency than L1 blocks

## Design

### Node start-up

At start-up, the host checks if one of the following applies:

* They are the genesis node and are an aggregator
* They are *not* the genesis node and are *not* an aggregator

If neither of these conditions is met, the host shuts down.

The identity of the sequencer is listed in the management contract on the L1. This allows other nodes to verify that 
the rollups are created by the sequencer.

### Creation of light batches

On each block, the sequencer's host feeds a set of transactions to the enclave. The enclave responds by creating a 
signed and encrypted *light batch*. This light batch is formally identical to the rollup of the final design, including 
a list of the provided transactions and a header including information on the current light batch and the hash of the 
"parent" light batch.

The sequencer's host immediately distributes the light batch to all other nodes. Unlike in the final design, these 
light batches are not sent to be included on the L1.

The linkage of each light batch to its parent ensures that the sequencer's host cannot feed the enclave a light batch, 
use RPC requests to gain information about the contents of the corresponding transactions, then feed the enclave a 
different light batch (e.g. where the sequencer performs front-running) to be shared with peers.

From the user's perspective, the transactions in the light batch are considered final (e.g. responses to RPC calls from 
the client behave as if the transactions were completely final).

### Creation of rollups

Every `x` blocks, the sequencer's host requests the creation of a rollup. This rollup contains all the light batches 
created since the last rollup, in a Merkle tree structure. This rollup is sent to be included on the L1.  The rollup is 
signed and encrypted by the sequencer.

A rollup is produced whenever one of the following conditions is met:

* The number of transactions across all light-batches since the last rollup exceeds `y`
* The value of all transactions across all light-batches since the last rollup exceeds `z`
* The time since the last rollup was produced exceeds `w`

`y`, `z` and `w` are configurable per network.

The nodes scan incoming blocks to retrieve this rollup. They validate the received rollup by:

* Checking that it is produced by the designated sequencer
* Checking that it contains all the light batches since the last rollup

They then persist the rollup, so that they have a record of which light batches have been confirmed on the L1.

## Future work

* Allowing nodes to challenge the sequencer's rollups (e.g. if the light batches are missing transactions, or if a 
  certain light batch is not included in the rollup)
* Creation of an inbox to allow transactions to be "forced through" if the sequencer is excluding them
* High-availability of the sequencer

## Unresolved issues

* How do we prevent the sequencer from running `n` enclaves and using `n-1` of them to test the impact of various 
  transaction sets (e.g. to identify front-running opportunities)?
