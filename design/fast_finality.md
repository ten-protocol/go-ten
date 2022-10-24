# Block explorer design

## Scope

The introduction of a centralised sequencer to improve cost-per-transaction during the bootstrapping phase, as 
described in the [Bootstrapping Strategy design doc](./Bootstrapping_strategy.md).

## Assumptions

* There is a single aggregator that is also the genesis node. This node assumes the role of sequencer

## Constraints

* Rollups are only published to the L1 every `x` blocks, where `x` >> 1

## Design

### Node start-up

At start-up, the host checks if one of the following applies:

* They are the genesis node and are an aggregator
* They are *not* the genesis node and are *not* an aggregator

If neither of these conditions is met, the host shuts down.

The identity of the sequencer is listed in the management contract on the L1. This allows other nodes to verify that 
the rollups are created by the sequencer.

### Creation of light batches

On each block, a *light batch* is created. This light batch is formally identical to the rollup of the final design, 
including a list of finalised transactions and a header including information on the current light batch and the 
hash of the "parent" light batch. The light batch is signed and encrypted by the sequencer.

As in the final design, each light batch is linked to a specific L1 block, and (at most) a single light batch is 
created per L1 block. The light batch is immediately distributed by the sequencer to all other nodes. Unlike in the 
final design, these light batches are not sent to be included on the L1.

From the user's perspective, the transactions in the light batch are considered final (e.g. responses to RPC calls from 
the client behave as if the transactions were completely final).

### Creation of rollups

Every `x` blocks, the sequencer creates a rollup. This rollup contains all the light batches created since the last 
rollup, in a Merkle tree structure. This rollup is sent to be included on the L1.  The rollup is signed and encrypted 
by the sequencer.

The nodes scan incoming blocks to retrieve this rollup. They validate the received rollup by:

* Checking that it is produced by the designated sequencer
* Checking that it contains all the light batches since the last rollup

They then persist the rollup, so that they have a record of which light batches have been confirmed on the L1.

## Future work

* Allowing nodes to challenge the sequencer's rollups (e.g. if the light batches are missing transactions, or if a 
  certain light batch is not included in the rollup)
* Creation of an inbox to allow transactions to be "forced through" if the sequencer is excluding them

## Unresolved issues

* Where do we store `x`, the frequency with which rollups are produced? Can it be changed at the sequencer's sole 
  discretion?
* Should the creation of light batches be based on a proof-of-work timer, or should they be produced per block?
