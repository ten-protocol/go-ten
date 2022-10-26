# Fast finality design

## Scope

The introduction of a single sequencer to improve cost-per-transaction during the bootstrapping phase, as 
described in the [Bootstrapping Strategy design doc](./Bootstrapping_strategy.md).

## Requirements

* Finality
  * Transaction *soft* finality (finality guaranteed by the sequencer) is achieved in under one second
  * There is eventual transaction *hard* finality (finality guaranteed by the L1)
  * The sequencer achieves hard finality on an agreed cadence
  * Transactions are hard-finalised in the same order they are soft-finalised
  * The sequencer is not able to "rewrite history" (or is strongly disincentivised from doing so), even for soft-final
    transactions
* Censorship resistance
  * Nodes can bypass the sequencer to include transactions on the L1 directly (possibly at higher cost and slower 
    finality)
  * The sequencer distributes all soft-finalised transactions to all nodes
  * Nodes can easily prove whether a given soft-finalised transaction was hard-finalised, and in the correct order
* Value-extraction resistance
  * The sequencer cannot precompute the effects of running a given set of transactions without committing to that set 
    of transactions (e.g. by running a single transaction, then `eth_call`ing to see how it has affected a smart 
    contract's state)
* Cost
  * L1 transaction costs can be driven lower at the expense of extending the hard-finality window
* User/dev experience
  * The responses to RPC calls reflect the soft-finalised transactions, and not just the hard-finalised transactions
* Operations
  * The sequencer is highly available; the failure of a single component does not impact its ability to deliver on any 
    of the requirements above

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

The sequencer's identity is listed in the management contract on the L1. This serves two purposes:

* It allows other nodes to verify that the light batches and rollups are created by the sequencer
* It prevents non-sequencer nodes from entering "sequencer mode" to evaluate the impact of a given light batch

In the management contract, the sequencer's identity is given as a set of enclave attestations. Each attestation 
matches one of the sequencer's enclaves, and contains the hash of that enclave's key. Since the attestations are unique 
per machine, the enclaves cannot be impersonated. The foundation admin will then whitelist these attestations. If one 
of the sequencer's enclaves goes down irrecoverably, a replacement attestation is added and whitelisted, following the 
same process.

### Production of light batches

A light batch is produced on the required cadence to meet the network's soft-finality window of one second.

To produce a light batch, the sequencer's host feeds a set of transactions to the enclave. The enclave responds by 
creating a signed and encrypted *light batch*. This light batch is formally identical to the rollup of the final 
design, including a list of the provided transactions and a header including information on the current light batch and 
the hash of the "parent" light batch.

The sequencer's host immediately distributes the light batch to all other nodes. Unlike in the final design, these 
light batches are not sent to be included on the L1.

The linkage of each light batch to its parent ensures that the sequencer's host cannot feed the enclave a light batch, 
use RPC requests to gain information about the contents of the corresponding transactions, then feed the enclave a 
different light batch (e.g. one where the sequencer performs front-running) to be shared with peers. The enclave will 
automatically reject a second light batch with the same parent.

From the user's perspective, the transactions in the light batch are considered final (e.g. responses to RPC calls from 
the client behave as if the transactions were completely final).

### Production of rollups

A rollup is produced whenever one of the following conditions is met:

1. The number of transactions across all light-batches since the last rollup exceeds `x`
2. The total size of all transactions across all light-batches since the last rollup exceeds `y`
3. The number of blocks since the last rollup was produced exceeds `z`

`x`, `y` and `z` are configurable per network. Rules (1) and (2) ensure that the rollup can fit within the Ethereum gas 
limit. Rule (3) reduces the risk associated with the sequencer failing.

To produce a rollup, the sequencer's host requests the creation of a rollup. The sequencer's enclave produces a rollup 
containing all the light batches created since the last rollup, in a sparse Merkle tree structure. This rollup is 
encrypted and signed, then sent to be included on the L1.

### Discovery of rollups

Nodes scan incoming L1 blocks for new rollups. They validate each new rollup by:

* Checking that it is produced by the designated sequencer, based on the sequencer listed in the management contract
* Checking that it contains all the light batches produced since the last rollup. Each light batch contains the number 
  of the rollup that will contain it. Since the rollup is a sparse Merkle tree, proving non-inclusion of a given light 
  batch is straightforward

They then persist the rollup, so that they have a record of which light batches have been confirmed on the L1.

## Future work

* Allowing nodes to challenge the sequencer's rollups (e.g. if the light batches are missing transactions, or if a 
  certain light batch is not included in the rollup)
* Creation of an inbox to allow transactions to be "forced through" if the sequencer is excluding them
* High-availability of the sequencer

## Unresolved issues

* Select a design for preventing value extraction (see the section 
  `Designs considering for preventing value-extraction`, below)
* How do we ensure the sequencer distributes all light batches to all nodes?
* Do the light batches need to be linked to the latest block that was fed into the enclave?
* How do we achieve high-availability of the sequencer, both for the enclave and for the host? How do we prevent 
  denial-of-service attacks on both?
* How do we prevent attackers from "gumming up" the sequencer with random, fake light-batches?

## Appendices

### Designs considering for preventing value-extraction

In this section, we investigate various designs to prevent value-extraction. The specific attack we have in mind is one 
where the sequencer runs `n` enclaves, and uses `n-1` of them to test the impact of various transaction sets to 
identify value-extraction opportunities.

#### Do nothing

In this approach, we rely on trust in the sequencer and the fact that value-extraction opportunities are reduced in
Obscuro due to its data-visibility rules.

#### Disable `eth_call` on sequencer enclaves

By disabling `eth_call` on the sequencer enclaves, we prevent the operator from extracting information about the impact
of a given light batch.

This approach is unworkable. The operator can run a separate, non-sequencer enclave, feed the light batches to that
enclave, and use `eth_call`s on that enclave to determine the impact of a given light batch.

#### Run a single sequencer enclave

If there is a single sequencer enclave, there are no other enclaves to use to identify value-extraction opportunities. 
The single sequencer cannot be restarted to identify value-extraction opportunities, since the enclave start-up delay 
will then prevent the sequencer from reaching its block production target.

This approach is unworkable. We cannot achieve the desired high-availability with a single sequencer enclave.

#### Detect restarts on sequencer enclaves

We allow validators to detect how often sequencer enclaves are restarted, incentivising good behaviour on behalf of the 
sequencer, and allowing the issue of a malicious (or incompetent) sequencer to be handled as a governance action. For 
this to work, validators would have to actively assess whether the sequencer is doing an adequate job.

There are two flavours of this.

##### Have sequencer enclaves produce lifetime proofs

Every `x`th light batch includes a proof of how long (e.g. in terms of light batches or L1 blocks) each sequencer 
enclave has been up. This creates a history of when each sequencer enclave was restarted. `x` can be arbitrarily high, 
since you can work backwards from this proof and the previous proof to check the enclave has been up the entire time.

This history can be queried via RPC from validators.

In this model, the sequencer operator would get `n-1` shots at front-running before having to restart one or more 
sequencer enclaves.

##### Include proofs from _all_ sequencer enclaves in each light batch

Every light batch contains signatures from all sequencer enclaves that are currently up. Because enclaves have a 
start-up delay, restarting the enclave will cause it to miss one or more light-batch signatures. A single signature is 
sufficient for the light batch to be accepted, but this creates a history of which sequencer enclaves were up at each 
point in time.

This history can be queried via RPC from validators.

In this model, the sequencer operator would have to restart a sequencer enclave after every shot at front-running.
