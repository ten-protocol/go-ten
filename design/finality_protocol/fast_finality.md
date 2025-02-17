# Fast finality design

## Scope

The introduction of a single sequencer to improve cost-per-transaction during the bootstrapping phase, as 
described in the [Bootstrapping Strategy design doc](./sequencer_bootstrapping_strategy.md).

## Glossary

* *Batch*: A set of transactions produced by the designated sequencer, considered to be soft-finalised in that order
* *Sequencer*: The network's single aggregator that is responsible for disseminating batches and posting rollups

## Assumptions

* Enclaves can be hacked; a signature from an attested enclave cannot be taken as an absolute proof that an attacker 
  does not have control over the enclave's execution
* The market does not sufficiently trust the TEN foundation to allow it to run the sequencer unchecked; we 
  therefore need some validation of its execution to provide assurance
* The adoption of the sequencer model may be long-lived, and must therefore be fully production ready
* The enclave has a start-up delay that exceeds the production time for batches (see below)

## Requirements

* Decentralisation
  * The network has a single sequencer that is also the genesis node
  * The entity operating the sequencer gains no power over the protocol's execution, other than the ability to 
    temporarily halt posting of rollups
  * Validators can easily verify the attestation over the sequencer's enclave
* Finality
  * Transaction *soft* finality (finality guaranteed by the sequencer) has a median duration of one second, from when 
    the user submits the transaction to when the user receives the transaction receipts and events
  * There is eventual transaction *hard* finality (finality guaranteed by the L1)
  * The sequencer is strongly incentivised to hard-finalise transactions in the same order they are soft-finalised, and 
    on the agreed cadence
  * Nodes can easily prove whether a given soft-finalised transaction was hard-finalised, and in the correct order
  * The sequencer is not able to "rewrite history" (or is strongly disincentivised from doing so), even for soft-final
    transactions
  * Soft finality does not break the mechanism used by some smart contracts (e.g. SushiSwap) of using `block.number` as 
    a rudimentary clock
* Censorship resistance
  * End-users have a mechanism to force the sequencer to include their transactions (possibly at higher cost and slower 
    finality)
  * The sequencer distributes all soft-finalised transactions to all nodes
* Value-extraction resistance
  * The sequencer cannot precompute the effects of running a given set of transactions without committing to that set 
    of transactions (e.g. by running a single transaction, then `eth_call`ing to see how it has affected a smart 
    contract's state)
* Cost
  * L1 transaction costs can be driven lower at the expense of extending the hard-finality window
  * The network can be configured to only publish rollups to the L1 every `x` blocks, where `x` >> 1
  * The network can be configured to produce batches (see below) at a higher frequency than L1 blocks
* User/dev experience
  * The responses to RPC calls reflect the soft-finalised transactions, and not just the hard-finalised transactions
* Resilience
  * Failover of the sequencer does not require a governance action
  * In the case of an L1 reorg, the sequencer must republish the rollups in the same order

## Non-requirements

* Censorship resistance
  * If censorship-resistance is achieved through an L1 inbox mechanism, it is acceptable for the costs of the 
    associated L1 transactions (which will be greater than the TEN costs) to fall on the transaction submitter
* Resilience
  * During failover, upgrade or planned maintenance of the sequencer, it is acceptable to break the one-second 
    soft-finality guarantee, and drop transactions that have not been soft-finalised (see e.g. 
    https://status.optimism.io/history/2, where Optimism regular has five-minute upgrade windows)

## Design

Here is an overview of the sequencer design:

![architecture diagram](./resources/fast_finality_arch.jpeg)

### Sequencer identity

The sequencer's identity is given in the management contract on the L1 as a set of enclave attestations. This allows 
other nodes to verify that the batches and rollups are created by the sequencer. We use a set instead of a single 
attestation to allow faster failover in the case of a sequencer node crashing (see the section `Resilience`, below).

Each attestation matches one of the sequencer's enclaves, and contains the hash of that enclave's key. Since the
attestations are unique per machine, the enclaves cannot be impersonated. The foundation admin will then whitelist
these attestations. If one of the sequencer's enclaves goes down irrecoverably, the old attestation can no longer be 
used (since it is tied to the machine). A replacement attestation can then be added and whitelisted to induct a 
replacement enclave into the sequencer's cluster, following the same process.

### Production of batches

A batch is produced on the required cadence to meet the network's soft-finality window of one second. Only the 
sequencer produces batches, with the enclave code enforcing that only whitelisted enclaves can produce batches.

To produce a batch, the sequencer's host feeds a set of transactions to the enclave. The enclave responds by creating a 
signed and encrypted *batch*. This batch is formally identical to the rollup of the final design, including a list of 
the provided transactions and a header including information on the current batch and the hash of the "parent" batch. 
This header includes a monotonically increasing counter, in order to support smart contracts that use 
`block.number` as a rudimentary clock (e.g. SushiSwap). The header also includes the height of the associated L1 block
(rather than the L1 block hash, to make it easier to handle reorgs).

The sequencer's host immediately distributes the batch to all other nodes, who gossip it onwards to other nodes
(ensuring the sequencer cannot restrict the distribution of batches to specific nodes, provided one of the nodes who 
received the batch is honest). These batches are not sent to be included on the L1.

When a node receives a batch, it first checks whether the batch is valid. If not, it rejects the batch as not part of 
the canonical chain. It then checks that it has also stored the batch's parent. If not, it walks the chain backwards, 
requesting any batches it is missing until it hits a stored batch. In the current design, it requests these batches 
from random nodes; once a gossip protocol is implemented, it will request the batches from its known peers.

The linkage of each batch to its parent also ensures that the sequencer's host cannot feed the enclave a batch, use RPC 
requests to gain information about the contents of the corresponding transactions, then feed the enclave a different 
batch (e.g. one where the sequencer performs front-running) to be shared with peers. The enclave will automatically 
reject a second batch with the same parent.

From the user's perspective, the transactions in the batch are considered final (e.g. responses to RPC calls from the 
client behave as if the transactions were completely final).

### Production of rollups

Only the sequencer produces rollups, with the enclave code enforcing that only whitelisted enclaves can produce rollups.

A rollup is produced whenever one of the following conditions is met:

1. The number of transactions across all batches since the last rollup exceeds `x`
2. The total size of all transactions across all batches since the last rollup exceeds `y`
3. The number of blocks since the last rollup was produced exceeds `z`

`x`, `y` and `z` are configurable per network. Rules (1) and (2) ensure that the rollup can fit within the Ethereum gas 
limit. Rule (3) reduces the risk associated with the sequencer failing.

To produce a rollup, the sequencer's host requests the creation of a rollup. The sequencer's enclave produces a rollup 
containing all the batches created since the last rollup, in a sparse Merkle tree structure. This rollup is encrypted 
and signed, then sent to be included on the L1.

The management contract on the L1 verifies that the rollup is produced by a designated sequencer.

### Discovery of rollups

Nodes scan incoming L1 blocks for new rollups. For each new rollup, the node checks that it contains all the batches 
produced since the last rollup. Each batch contains the number of the rollup that will contain it. Since the rollup is 
a sparse Merkle tree, proving non-inclusion of a given batch is straightforward.

If they discover an issue, they can mount a challenge, as per the "Staking and slashing" section, below. They also 
reject the rollup, waiting for a valid rollup at the same height to be produced.

If the rollup is valid, the node persists the rollup, so that they have a record of which batches have been confirmed 
on the L1.

### Staking and slashing

The sequencer must put up a stake in the management contract.

If a node finds that the contents of a given batch do not match the rollup published on the L1 (e.g. transactions 
missing, transactions in the wrong order), it can post a challenge including the batch and offending rollup to the L1. 
The management contract will inspect this challenge. If successful, the sequencer will be slashed, with their entire 
stake split between the foundation and the challenger as a reward. This reward is greater than the cost of posting the 
challenge, to incentivise prompt discovery of issues, but less than the total stake, to reduce the incentive to mount 
an attack to win the stake.

Initially, there will be no slashing mechanism to ensure batches are produced at the correct cadence. It is expected to 
be a sufficient incentive that the sequencer is operated by the foundation and has an interest in the well-running of 
the network, including a reliable batch cadence. Any failure to preserve the correct batch cadence will be immediately 
visible to all observers.

### Resilience

#### Goals

The sequencer holds three important types of data:

1. Transactions submitted but not yet included in a batch
2. Batches (including their transactions)
3. Rollups

In the case of failover, (1) can be dropped if needed, while (3) can be recovered from the L1 chain (once submitted) or 
recreated (provided the batches are available). Thus, during failover, the key concern in terms of data resiliency is 
(2).

In addition, while it is acceptable to break the one-second soft-finality guarantee during failover, we should still 
seek to minimise the recovery time. Solutions that require, for example, the full reingestion of the L1 chain are 
unworkable.

#### Cluster configuration

To achieve the desired data resiliency and recovery times, we must achieve resiliency of both the sequencer's host and 
the sequencer's enclave.

On the host side, the sequencer's host can use a database with resiliency configured. When they receive a batch, they 
must store it and wait for the confirmation before proceeding. If the host crashes, we simply restart it and restore it 
from its resilient database.

On the enclave side, the sequencer can run a cluster of `n` enclaves, each backed by a separate database. All the 
enclaves are active at once. A leader enclave is selected to be the sole producer of batches and rollups, while the 
follower enclaves behave like regular validator enclaves, receiving the batches via a gossiping process and retrieving 
the rollups from the L1.

The cluster's leader is selected via an RPC operation on the host. It is the responsibility of the sequencer's operator 
to monitor the healthiness of the host and enclaves. In the event that a follower crashes, it can be restarted and 
recover data from the host, just like a regular node. If the event that the leader crashes, the sequencer operator 
must select a new leader.

The key risk during failover to a new leader is that a single batch (the latest) may be lost. There are two specific 
issues that must be handled:

1. A new leader may come online and be missing the latest batch
2. The original leader comes back online, but their database contains a batch that was never distributed, and now 
   represents a fork

To avert (1), the sequencer's host must carefully update the new leader before it starts execution.

To avert (2), we need to be able to overwrite the state of the recovered leader. However, this must be handled 
carefully - if a node's state can be overwritten arbitrarily, the node can be used to front-run by repeatedly writing 
new batch chains and inspecting the results. To address this, an enclave's batch chain can only be overwritten 
immediately after start-up. The recovered leader will poll its former followers for conflicting batches, and overwrite 
any batches at the same height. Once the leader starts normal execution, this overwriting mechanism will be disabled 
(and thus an enclave restart, with the attendant start-up delay, must be incurred to overwrite again and inspect the 
results).

#### Alternatives considered

This approach was selected over a number of alternatives.

##### As above, but with `n` hosts all speaking to `m` enclaves

The selected approach is simpler and more closely aligned to our current implementation (which assumes one enclave per 
host, and vice-versa).

##### Having a single node that is restored from backup

This approach has several downsides:

* Recovery would be much slower in this approach, as a governance action would be required to whitelist the new 
  sequencer attestation in the management contract (this could be mitigated by preallocating some future enclave 
  machines to use for failover, and whitelisting attestations for them in the management contract)
* If the database is lost and has to be restored from backup, the latest batches (those not contained in the backup) 
  would have to be recovered from network peers, which would be more complicated than recovering them from specific, 
  sequencer-operator controlled nodes. In particular, we'd to handle the case of a node receiving a batch from the 
  crashed leader that they then fail to gossip out correctly; some mechanism would have to be provided to allow them to 
  return to the non-forked batch chain (e.g. overwriting the batch chain if they receive another with greater height)
* This approach is more difficult operationally (creation, storage and recreation from backups)

## Future work

* Creation of an inbox to allow transactions to be "forced through" if the sequencer is excluding them (the rough idea 
  is that there will be an inbox for transactions in the management contract, and validators will reject batches that 
  do not contain any (valid) transactions that have sat in the inbox for a certain amount of time)
* Implement a mechanism to prevent value extraction (see the section 
  `Possible designs for preventing value-extraction`, below)
* Implement an incentive mechanism to ensure we achieve the desired cadence of batches without tying ourselves to the 
  L1 block cadence, in light of the absence of time within an enclave. Can the host request batches on the correct 
  cadence, with some incentive mechanism to prevent deviations from the "correct" cadence?

## Unresolved issues

* How do we prevent denial-of-service attacks on the sequencer?
* How should the counter in batch headers increase? Should it increase once per L1 block, or once per batch?
* How do we recover from a malicious sequencer, post a successful challenge? Is it sufficient for validators to just 
  wait for a successful rollup to be produced? Is there some unpicking or cleanup that needs to be done?
* How do we prevent deposits from the L1 having to wait the full rollup duration to land on the L2? Do we include the 
  withdrawals instructions in the batch header?

## Appendices

### Possible designs for preventing value-extraction

In this section, we investigate various designs to prevent value-extraction. The specific attack we have in mind is one 
where the sequencer runs `n` enclaves, and uses `n-1` of them to test the impact of various transaction sets to 
identify value-extraction opportunities.

#### Do nothing

In this approach, we rely on trust in the sequencer and the fact that value-extraction opportunities are reduced in
TEN due to its data-visibility rules.

#### Disable `eth_call` on sequencer enclaves

By disabling `eth_call` on the sequencer enclaves, we prevent the operator from extracting information about the impact
of a given batch.

This approach is unworkable. The operator can run a separate, non-sequencer enclave, feed the batches to that enclave, 
and use `eth_call`s on that enclave to determine the impact of a given batch.

#### Run a single sequencer enclave

If there is a single sequencer enclave across the entire network, there are no other enclaves to use to identify 
value-extraction opportunities. The single sequencer cannot be restarted to identify value-extraction opportunities, 
since the enclave start-up delay will then prevent the sequencer from reaching its block production target.

This approach is unworkable. We cannot achieve the desired recovery times with a single sequencer enclave.

#### Detect restarts on sequencer enclaves

We allow validators to detect how often sequencer enclaves are restarted, incentivising good behaviour on behalf of the 
sequencer, and allowing the issue of a malicious (or incompetent) sequencer to be detected and handled as a governance 
action. For this to work, validators would have to actively assess whether the sequencer is doing an adequate job.

There are two flavours of this.

##### Have sequencer enclaves produce lifetime proofs

Every `x`th batch includes a proof of how long (e.g. in terms of batches or L1 blocks) each sequencer enclave has been 
up. This creates a history of when each sequencer enclave was restarted. `x` can be arbitrarily high, since you can 
work backwards from this proof and the previous proof to check the enclave has been up the entire time.

An alternative model is for each sequencer enclave to post a restart proof to the L1 or batch chain whenever it starts, 
and wait for that proof to be included in the chain before continuing execution (thus forcing the proof to be posted on 
each restart).

This history can be queried via RPC from validators.

In this model, the sequencer operator would get `n-1` shots at front-running without being detected before having to 
restart one or more sequencer enclaves. This front-running would come at the expense of making the enclave useless 
until it's restarted.

##### Include proofs from _all_ sequencer enclaves in each batch

Every batch contains signatures from all sequencer enclaves that are currently up. Because enclaves have a start-up 
delay, restarting the enclave will cause it to miss one or more batch signatures. A single signature is sufficient for 
the batch to be accepted, but this creates a history of which sequencer enclaves were up at each point in time.

This history can be queried via RPC from validators.

In this model, the sequencer operator would have to restart a sequencer enclave after every shot at front-running.
