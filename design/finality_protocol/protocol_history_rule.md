# L2 History Rule

### Context

- The Sequencer is linked to an Ethereum node and has a real-time view of the "canonical L1 chain". This view might be different from the view of the other Obscuro nodes.
- The result of processing an L2 block (batch) depends on an L1 block because of cross-chain messages.
- As a consequence, to process an L2 block, an Obscuro node needs the L1 block (and its entire history) that was linked with the batch.


### Problem

Given that an L1 block is not final in real-time, it is theoretically possible that the sequencer produces batches linked to non-canonical L1 blocks.
In the current protocol, the sequencer will produce a sibling batch linked to the new canonical L1 block when this happens.
This is fine if the validators have received all the l1 blocks because they can execute and validate each batch. 
Based on what is canonical, the enclave will use the right calculated results.

There is an extreme edge case (which could be used to attack Obscuro) where an L1 block was visible only to the sequencer and not to any of the validators. 
If the sequencer has produced batches using this block as a reference, then the result of this batch can't be calculated by any validator.

Note: this attack is not performed by the sequencer. The attacker is an Ethereum staker that somehow sends a block to the Obscuro Sequencer that it doesn't broadcast to the network.


### Properties of the ideal solution

- the requirement is that the network must progress, even in this case.
- given this is an extreme edge case, the solution doesn't need to be optimised for it.

Note: Without this "data availability" problem, an Obscuro node can just wait to receive all the L1 blocks it needs before processing an L2 batch that it received via p2p.

To restate, some of the blocks might never arrive so this approach could freeze the network.


## The Solution

The solution has two elements:
1. a rule for producing batches
2. a rule for synchronising a validator with the sequencer

### Rule 1

**If multiple batches have the same height (number), the one with the highest sequence number is considered canonical.**

There are multiple consequences.

The sequencer cannot build on top of a batch chain that is non-canonical. This means it must create a "twin" of a non-canonical batch if that batch has suddenly become canonical.

### Rule 2

An Obscuro node receiving batches via p2p must wait for a "canonical" batch for which it has the l1 block.

This generalises the simpler rule that doesn't handle the DA issue.

The significance is that the node will wait until it has all the data required to process that sequence of batches.

Another way to put it is that, in case the Validator is out-of-sync with the sequencer ( because the validator does not have an l1 block that the sequencer considers canonical), 
the validator will wait until it is in sync, before moving forward.


## Solution that uses this rule






