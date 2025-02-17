# Design for the Ten escape hatch

*Note: this document is still a WIP*

This document proposes a design for an "Escape Hatch". 
It describes the ultimate way by which users can exit Ten.


## High level requirement

The "escape hatch" is the mechanism that kicks in when **all hope is lost**, the last resort. 
Something that happens just before the network goes down permanently for some unforeseen reason.

For example, wen the central sequencer is no longer able to produce blocks and the Ten foundation is unable to replace it with something working in due time.

The high level requirement is that even in this situation, users should be able to exit the funds they have locked up in the bridge.

Note that the escape hatch should not be used against censorship resistance. There is another mechanism in place.


### Condition to trigger the "Escape Hatch"

The "Escape mode" will be a flag on the Management contract, that can be set under certain conditions.

- when there were no rollups published for the past x hours ( where x > 48hrs)
- when the foundation with a majority of y ( where y >> 66.6%) votes for it
- ...


### Assumptions

There will be at least one node in possession of the master seed that is able to publish it when this event is triggered.


## High level overview of the solution

Each rollup will contain an additional calldata field.
This field will be a list of hashes, one hash for each individual bridged asset.

This hash is the root hash of the Patricia Merkle State of that asset. 

*Note that this means that there needs to be a standard for how a bridgable ERC20 contract store the data*
*Note that also, the platform needs to add some entropy to each of these trees to avoid leaking information about what changed*

When the escape mode is triggered, the bridge contract will load up the roots from the latest rollup, and will start accepting user requests.

Each user request must be a Merkle proof where the leafs are ``(account,balance)`` with a valid chain to the root of tha asset.
The bridge will be happy to release that amount to the same account after verifying the merkle proof.

*Note that this means that the request doesn't have to be signed because the balance is credited to the same account as in the L2*

Any user who knows the acount and the balance for a certain asset, should be able to just create a transaction (maybe with the help of some tooling), and 
retrieve their funds.

*Note that the state root hashes should reflect the state after the withdrawal instructions from the rollup header.*

#### Revelation

The escape hatch is an apocalyptic event when everything is revealed.
Setting the flag in the management contract is visible in any enclave who will return the keys to decrypt all blocks.

It is assumed that when that happens, participants from the community will make all data public.

This means that any user is able to lookup their own balance, and create the proof required to exit it from the bridge.


## Assumptions and Problems

1. It is possible to use the (account,balance) pair and generate an easily verify-able proof to a root hash. (see https://github.com/zmitton/eth-proof)
2. What happens if there was a challenge on the latest rollup?
3. It is possible to accurately understand the validity of a L1 block in the new PoS world. The foundation should not be able to trigger the revelation mode on a fork of ethereum.
