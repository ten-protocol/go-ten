# The soft-finality guarantee 

The decentralised POBI protocol links the rollups to the L1 block used as proof, which meant that any L1 reorg
results in an implicit L2 reorg. The finality of a rollup is 100% linked to the finality of an L1 block.

Bootstrapping .. central sequencer .. light batches .. see fast finality design.
Tldr; we want a central sequencer to guarantee finality ~1second, and for the bride to be decentralised and resistant to reorgs.
(See bridge design)


## Problem

The sequencer can only guarantee the ordering of transactions which it receives from the Obscuro users.
It can't offer guarantees over the finality of cross-chain transactions. E.g.: deposits to the bridge.

Given that deposits affect the balances of accounts, which can be used in Obscuro transactions, there is a dependence 
between the ordering of the L1 messages and the result of executing L2 transactions.

This means there is a tension between the guarantees of "soft-finality" and the re-org resistance of the bridge.

In the "fast finality" design, the "soft finality" of an L2 transaction is defined as the moment when the sequencer produces
a LB containing that transaction. 

This design analyses what should happen to a Light Batch and a Rollup when deposit transactions disappear from the L1. 


## Requirements

In the context of the fast finality and the bridge designs, we need to define what "Soft finality" is, and in what conditions
can the sequencer generate competing Light batches.

## Assumption

Messages from the L1 are processed as they are found in the blocks and synthetic transactions are generated based on them.
These synthetic transactions are included in the light batches.


## Solutions

### Option 1 - Adding delays

The straight forward solution is to just wait before processing deposits.
Eventually, the cost of mounting a re-org attack on the L1 become prohibitively high, which ensures the security.

There are two main disadvantages:

- users have to wait before bridge messages are processed
- the solution is only as secure as the delay.

The subtlety of this approach is that it can transform Obscuro into a side-chain even though we rollup to Ethereum.

(Because finality is static.)


### Option 2 - Link light batches to blocks

Similar to POBI, each LB is generated from a single block and links to it.

The difference is that LBs are not published to the L1, so they are not rejected by a source of truth authority the moment 
they are no longer valid.

This can lead to chaotic results, and hard to enforce behaviour.

The sequencer monitors the L1, and the moment a reorg is happening, the sequencer will produce new LBs, which will be distributed.


### Option 3 - Link light batches to blocks, but guarantee the L2 ordering

The key insight for this option is that sequencer *can* offer guarantees over the order of the L2 transactions themselves.

Instead of being challenged over the LB, the sequencer can be challenged over the root hash of the transactions alone.
A sequencer will be able to generate sibling LBs in case an L1 fork is happening.


```go
type LightBatchHeader struct{
	ParentLB        Hash    // the id of the parent Light Batch
	L1BlockHash     Hash    // the head l1 block at the time of this LB
	L2TxsHash       Hash    // the MTree root of the list of transacions
	Number          int     // the height of this LB - since when? the latest L1 block, the latest rollup, the beginning?
	StateRoot       Hash    // the state root hash
	R, S            big.Int // signature of enclave
}
```

In case there is a reorg, the SA can issue a new LB with the same `Number` as a previous one, as long as the `L2TxsHash`
stays the same.


### Option 4 - Guarantee the L2 ordering, but link batches to deposit messages

This is a variation of the above solution, where instead of linking the batch to an L1 block, it links a batch to a list
of l1 messages.

In case the order of L1 messages is preserved during an L1 reorg ( which is very likely unless the message sender actively
tries an attack), then there is no need for generating a sibling LB.


```go
type LightBatchHeader struct{
	ParentLB        Hash    // the id of the parent Light Batch
	L1Number        int     // the height of the l1 block used 
	L1MsgsHash      Hash    // the list of messages from the L1.
	L2TxsHash       Hash    // the MTree root of the list of transactions
	Number          int     // the height of this LB - since when? the latest L1 block, the latest rollup, the beginning?
	StateRoot       Hash    // the state root hash
	R, S            big.Int // signature of enclave
}
```

The drawback of this solution is some complexity.
As the validator process blocks and Light batches, they might not reach the same L1MsgsHash, and that needs handling


## TODO - Sealing a LB chain in a rollup and reorgs


