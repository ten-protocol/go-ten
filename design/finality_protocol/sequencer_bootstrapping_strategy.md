# Sequencer bootstrapping strategy

*Note: this document is still a WIP*

As we are getting closer to production and have already designed implemented the key components of Ten, it's time to think about our bootstrapping strategy.

The POBI protocol described in the whitepaper assumes a network with significant traction. 

To get there, we estimate TEN will need at least one year, most likely more.

Other things we have to take into consideration during bootstrapping is code maturity and security. In the industry, this is called a period where the "training wheels" are on.

The other major L2 networks (Arbitrum and Optimism), opted for a pragmatic approach where they started out with a centralised sequencer.

We propose that TEN starts out similarly to the L2s. Centralised block production and decentralised validation.


## Single block producer

This is in essence a very simplified "POBI" with a single aggregator (SA).

The SA is operated by the TEN Foundation, and is configured as a variable in the TEN Management Contract (MC) on Ethereum. 
Only the foundation has the power to set the designated SA.

Note that this means that the "Consensus problem" becomes relatively simple in this first stage. 
The SA unilaterally decides the ordering of txs and when to publish a rollup.

Clearly the SA needs to be monitored and reasonably HA, since it is a single point of failure.


### Decentralised solution

For the solution to still be overall decentralised, we need to add a couple of mechanisms.


#### Invalid rollup challenge mechanism

Given there is a single aggregator that can publish rollups, we need to introduce a challenge mechanism in case it 
malfunctions or turns malicious.

In POBI, there was no need for a challenge mechanism. In case of an invalid rollup, everyone would just ignore it and continue publishing rollups on a different fork.

Tis is no longer possible because in this simplified model only the SA can submit rollups. (Unless we add significant complexity) 

To keep it simple, the challenge mechanism will be disruptive in the first stage. Basically withdrawals will be paused, and need to be manually re-enable by the Foundation. A challenge can only be produced as a result of a hack, same as the SA attempting malicious behaviour. 

A challenge is a payload that an SGX enclave produces when a published rollup fails processing.
The node that produced this payload can submit it to the Management Contract (MC), which will accept it if the signature belongs
to one of the approved enclaves, and the challenged rollup matches one of the latest n ones.

```go
    type Challenge struct {
        ChallengedRollup    Hash
		Error               String
		R,S                 big.Int // Signature
    }
```

The MC has no way to evaluate the merits of challenge, so it will enter a "hacked mode", where it disables withdrawals, until this is sorted out.

*Note that to produce an invalid challenge, a hacker must break into the TEE* 

TBD - how to exit this state? Can we use an auction staking mechanism like the one we proposed in the WP?


#### Censorship resistance

A centralised block producer means we no longer have full censorship resistance on the l2 layer.

Note: There is still a lot of censorship resistance since the SA can't see the transactions, and they can be submitted through proxies to completely hide the source.

The following mechanism is an extra, in case that is not enough. 

On the MC on the L1, there will be an "inbox" method.
Users will submit encrypted l2 transactions to that inbox.

The SA has to include those transactions forced by the code running in SGX.
The enclave is fed all the l1 blocks, so there is no way to not include those transactions, short of hacking sgx or halting.

Not including those transactions by hacking sgx is discovered by the validators, who will issue a challenge.

Note that this solution is superior to the one the other centralised L2s provide, because our nodes ingest L1 blocks, so transactions have to be included quicker.


#### Escape hatch

See the  "Escape hatch" design



### Fast Finality

By publishing rollups less frequently we solved the cost problem (during bootstrapping).

The challenge now is to achieve practical fast finality.

What does "finality" mean from the pov of an end-user?

It's when the UI can confidently display a green tick after a submitted transaction, and can display the status of their account.

It's worth splitting finality into "soft finality" and "hard finality". They can be treated identically by UIs if the source of finality is trusted or the incentives are very good.
Ultimate "hard" finality is guaranteed by the finality of the rollup including a transaction in the L1.

"Soft finality" is always an approximation of that.
In the first stage we'll rely on trust that the SA will operate correctly, mostly to protect its reputation, with the plan of replacing it with incentives in the next stage.


### Requirements of a Protocol for fast finality

If we wouldn't have to worry about scalability, or if TEN functioned like transparent chains, then the SA could just emit events and receipts when it processes a tx, 
and send them to a caching layer to be consumed by clients.
Also, the state could be cached, so that "eth_call" requests could be handled from the caching layer.

Unfortunately, things are more complicated, because the result to events, receipts and rpc eth_call, depend on who is asking.
Only logic run inside an approved SGX is able to return a result.

The big technical challenge for this approach is scalability and HA.

One super simple option would be for all clients to just connect to the SA. The problem is that no matter how much hardware we throw at it, this will only scale up to a limit. 

*Note that We can start with this approach to prove the concept.*

The solution is to delegate the work of replying to users to the entire network of l2 nodes. These nodes are already incentivised read-only participants.

#### Incentives for L2 node operators

See the incentives design in go-ten.

#### Proposal for a protocol to keep L2 nodes in sync

The problem to be solved is how to keep the L2 network in sync with the real-time ordering performed on the SA.

Our business requirement is to respond to users with a "soft finality" in no more than 2s.
Note: That might be a bit ambitious given there are network latencies involved.

The other requirement is that all L2 nodes must act as read-only replicas of the SA, able to serve data requests on behalf of the SA.
Given that they are owned by byzantine hosts, the L2 node must only execute the transaction batches that are originating from the SA. 
If they were able to execute anything else, they could maliciously produce wrong results.

Users submit transactions to any L2 node, who gossip them among themselves, in a group that also includes the SA. (same as today)

Every second the enclave of the SA produces a signed "batch", which is a list of tx hashes, with a header pointing to the previous batch.
This batch is encrypted, and gossiped to all the other nodes.
Note: The enclave can only produce it once it calculated a chain of N hashes (roughly a second worth). This is necessary to prevent frontrunning by the SA. Basically the timer is not the clock of the host, but a mechanism based on a light proof of work.

Once the batch reaches the enclave of any other node, and they check it came from the SA, it will be similarly processed as a canonical rollup in POBI. 
It will be considered as soft-final, and it will generate receipts and events, which can be requested by the users who are connected to those nodes.

The SA itself will not respond to any user requests.
After N such light rounds, when the Aggregator has gathered enough txs, it submits a Rollup to the L1.


#### Implementation considerations

##### The batch as the canonical unit

In the scenario, from the point of view of the TenVM, each batch is the equivalent to an Ethereum block.

The rollups published to the L1 function more like logical checkpoints.

In case there is a discrepancy between the distributed signed batches, and the published rollup, it's a big problem.

One big problem in this approach is how to handle catching up based on the data published in the L1.
The node catching up needs to recreate the exact chain that everyone else that was live has.
This means that the published rollup needs to contain enough information to allow reconstruction.
Given that calldata space will be the most expensive resource, this has to be optimised.


##### The published rollup (PR) as the canonical unit

This is more like the current POBI approach.

The complexity here is that we need to implement a new special mechanism for batches. Basically maintain two parallel databases.
A temporary one, where the batch is the unit, to be used between PRs, and then the main database where the PR is the unit.

When a new node catches up, it will only create the PR database.



#### Sequence:

1. Txs are received to non-Agg nodes and gossiped around, all nodes store them by their hashes in mempool cache.
2. SA sequences them into a signed batch pointing to a previous batch ( or to a PR - for the second option)
3. SA performs some super light pow to get a timer which tells the enclave it's time to release a new batch, and broadcasts it to the other nodes.
4. Nodes receive batch, they execute the same txes by hash from their mempool against their internal state and check the final header hash matches the batch.
5. Nodes now have the full state to serve soft-final receipts and events, and respond to eth_call



### Failure mechanisms

This protocol introduces a few more failures.

#### Invalid soft finality challenge mechanism

It is possible that users receive a soft finality confirmation that proves to be incorrect.

This may be due to:
- malfunction of the protocol
- SA signed multiple competing batches
- L2 node operator found a hack and was able to generate invalid data

TBD How do we handle this?
