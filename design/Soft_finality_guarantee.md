# The soft-finality guarantee 

The decentralised POBI protocol links the rollups to the L1 block used as proof, which means that any L1 reorg
results in an implicit L2 reorg. The finality of an L2 rollup is fully linked to the finality of an L1 block.

During bootstrapping, we want a central sequencer to guarantee finality in around ~1 second and for the bridge to be decentralised 
and resistant to reorgs.

Note: See bridge and finality designs for more details


## Problem

The central sequencer can only guarantee a consistent ordering of the transactions it receives from the Obscuro users.
It can't offer guarantees over the finality of cross-chain transactions. E.g., deposits to the bridge.

Given that deposits affect the balances of accounts, which can be used in Obscuro transactions, there is a dependence 
between the ordering of the L1 messages and the result of executing L2 transactions.

In the "fast finality" design, the "soft finality" of an L2 transaction is defined as the moment when the sequencer produces
a Light Batch(LB) containing that transaction.

There is a tension between the guarantees of "soft-finality" and the re-org resistance of the bridge.

What does it mean that a transaction is soft-finalised:

- Executing this transaction can never result in a different outcome.
- Or that the execution result might be different under certain factors unlikely to happen unless there is a re-org attack.


## Requirements

- In the context of fast finality and bridge designs, we need to define "Soft finality".
- In what conditions can the sequencer generate competing light batches without being challenged?
- Define the structure of a Light batch.

## Assumption

Messages from the L1 are processed as they are found in the L1 blocks, and synthetic transactions are generated based on them.
These synthetic transactions are implicitly included in the light batches.


## Solutions

Since the problem is expressed in terms of "tension", we'll first try to explore the different options, and
understand where they sit on the tradeoff spectrum.


### Option 1 - Adding delays

The straightforward solution is to wait a while before processing deposits so that these transactions are "final" on the L1.
Eventually, the cost of mounting a re-org attack on the L1 becomes prohibitively high, ensuring security.

There are two main disadvantages:

Users must wait before bridge messages are processed, which is a bad UX and disables certain applications.
 The solution is only as secure as the delay.

The subtlety of this approach is that it can transform Obscuro into a side-chain even though we roll up to Ethereum.

This depends on the definition of an L2, which is a moving target. If the consensus is that an L2 needs a re-org resistant
bridge, then this option will not be favourable

On the tradeoff spectrum, this solution lies towards the: "The result of executing this transaction can never result in a different outcome" end,
which is good, but it has some severe disadvantages.


### Option 2 - Link light batches to blocks

Like POBI, each LB is generated from a single block and links to it. (Actually, multiple LBs will be linked to a single block)

The sequencer monitors the L1, and the moment a reorg is happening, the sequencer will produce new LBs linked to the blocks from 
the new branch, which will then be distributed.

The difference is that LBs are not published to the L1, so they are accepted by a "source of truth" authority the moment 
they are no longer linked to a block found on the canonical chain. After that, they are given to the other L2 nodes.

This can lead to chaotic results and hard-to-enforce behaviour. The sequencer could use the re-org reason to generate competing LBs,
and attempt some MEV.

Under this option, the result of executing a transaction will entirely depend on what happens on the L1.

This could be better because the re-org on the L1 might have no logical relationship whatsoever with the L2 transaction.


### Option 3 - Link light batches to blocks, but guarantee the L2 ordering

This option attempts to improve on the previous one by reducing the impact of L1 reorgs.

The critical insight is that the sequencer *can* offer guarantees over the L2 transactions' order.

Instead of being challenged over the result of executing transactions, the sequencer can be challenged over the root hash 
of the L2 transactions alone.

A sequencer will be able to generate sibling LBs in case an L1 fork is happening, but only as long as the order of the L2 
transactions are preserved.


```go
type LightBatchHeader struct{
	ParentLB        Hash    // the id of the parent Light Batch
	L1BlockHash     Hash    // the head l1 block at the time of this LB
	L2TxsHash       Hash    // the MTree root of the list of transactions
	Number          int     // the height of this LB - since when? the latest L1 block, the latest rollup, the beginning?
	StateRoot       Hash    // the state root hash
	R, S            *big.Int // signature of enclave
}
```

There can be sibling LBs (with the same `Number` and `ParentLB`) with the same `L2TxsHash` but different `L1BlockHash`. 

This option offers a stronger finality guarantee. A transaction can never be replaced with a competing one, but it might result
in a different outcome in the case, a deposit transaction was processed or not before it.

The sibling LBs will result in the same execution result for the vast majority of transactions. 


### Option 4 - Guarantee the L2 ordering, but link batches to deposit messages

To improve the guarantee even further, we can reduce the reasons for creating sibling LBs to only the situation where
there is an actual change in the ordering of L1 messages (deposits).

The insight is that instead of linking the batch to an L1 block, we can link a batch to a list of L1 messages.
If the order of L1 messages is preserved during an L1 reorg ( which is very likely unless the sender actively
tries an attack), then there is no reason for generating a sibling LB.

```go
type LightBatchHeader struct{
	ParentLB        Hash    // the id of the parent Light Batch
	L1Number        int     // the height of the l1 block used 
	L1MsgsHash      Hash    // the list of messages from the L1.
	L2TxsHash       Hash    // the MTree root of the list of transactions
	Number          int     // the height of this LB - since when? the latest L1 block, the latest rollup, the beginning?
	StateRoot       Hash    // the state root hash
	R, S            *big.Int // signature of enclave
}
```

The drawback of this solution is some complexity. As the validators process blocks and batches, they might reach different L1MsgsHash, which needs handling.

A transaction can never be replaced with a competing one, but it might result in a different outcome if there is a deliberate
double-spend attack on the L1 for a deposit transaction. The other accounts should be fine.


## Conclusion

The last option has the best tradeoffs.
The sequencer guarantees that it will consistently execute L2 transactions in the same order and batches, even if there are reorgs on the L1.
If there is an L1 re-org, but there are no L2 messages, or the L2 messages are included in the same order, then the finality guarantee is full.

The only case where finality can change is when there is a deliberate double-spend attack on the L1, and this will only affect
the user who is directly involved in that attack.


## TODO - Reason about sealing an LB chain in a rollup and reorgs


