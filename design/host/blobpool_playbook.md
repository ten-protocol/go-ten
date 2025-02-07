# Geth BlobPool Playbook


BlobPool is the transaction pool dedicated to EIP-4844 blob transactions. This document
was taken from the Geth source code [here](https://github.com/ethereum/go-ethereum/blob/93c541ad563124e81d125c7ebe78938175229b2e/core/txpool/blobpool/blobpool.go#L132-L293)

Blob transactions are special snowflakes that are designed for a very specific
purpose (rollups) and are expected to adhere to that specific use case. These
behavioural expectations allow us to design a transaction pool that is more robust
(i.e. resending issues) and more resilient to DoS attacks (e.g. replace-flush
attacks) than the generic tx pool. These improvements will also mean, however,
that we enforce a significantly more aggressive strategy on entering and exiting
the pool:
    - Blob transactions are large. With the initial design aiming for 128KB blobs,
    we must ensure that these only traverse the network the absolute minimum
    number of times. Broadcasting to sqrt(peers) is out of the question, rather
    these should only ever be announced and the remote side should request it if
    it wants to.
    
  - Block blob-space is limited. With blocks being capped to a few blob txs, we
    can make use of the very low expected churn rate within the pool. Notably,
    we should be able to use a persistent disk backend for the pool, solving
    the tx resend issue that plagues the generic tx pool, as long as there's no
    artificial churn (i.e. pool wars).
    
  - Purpose of blobs are layer-2s. Layer-2s are meant to use blob transactions to
    commit to their own current state, which is independent of Ethereum mainnet
    (state, txs). This means that there's no reason for blob tx cancellation or
    replacement, apart from a potential basefee / miner tip adjustment.
    
  - Replacements are expensive. Given their size, propagating a replacement
    blob transaction to an existing one should be aggressively discouraged.
    Whilst generic transactions can start at 1 Wei gas cost and require a 10%
    fee bump to replace, we suggest requiring a higher min cost (e.g. 1 gwei)
    and a more aggressive bump (100%).
    
  - Cancellation is prohibitive. Evicting an already propagated blob tx is a huge
    DoS vector. As such, a) replacement (higher-fee) blob txs mustn't invalidate
    already propagated (future) blob txs (cumulative fee); b) nonce-gapped blob
    txs are disallowed; c) the presence of blob transactions exclude non-blob
    transactions.
    
  - Malicious cancellations are possible. Although the pool might prevent txs
    that cancel blobs, blocks might contain such transaction (malicious miner
    or flashbotter). The pool should cap the total number of blob transactions
    per account as to prevent propagating too much data before cancelling it
    via a normal transaction. It should nonetheless be high enough to support
    resurrecting reorged transactions. Perhaps 4-16.
    
  - Local txs are meaningless. Mining pools historically used local transactions
    for payouts or for backdoor deals. With 1559 in place, the basefee usually
    dominates the final price, so 0 or non-0 tip doesn't change much. Blob txs
    retain the 1559 2D gas pricing (and introduce on top a dynamic blob gas fee),
    so locality is moot. With a disk backed blob pool avoiding the resend issue,
    there's also no need to save own transactions for later.
    
  - No-blob blob-txs are bad. Theoretically there's no strong reason to disallow
    blob txs containing 0 blobs. In practice, admitting such txs into the pool
    breaks the low-churn invariant as blob constraints don't apply anymore. Even
    though we could accept blocks containing such txs, a reorg would require moving
    them back into the blob pool, which can break invariants.
    
  - Dropping blobs needs delay. When normal transactions are included, they
    are immediately evicted from the pool since they are contained in the
    including block. Blobs however are not included in the execution chain,
    so a mini reorg cannot re-pool "lost" blob transactions. To support reorgs,
    blobs are retained on disk until they are finalised.
    
  - Blobs can arrive via flashbots. Blocks might contain blob transactions we
    have never seen on the network. Since we cannot recover them from blocks
    either, the engine_newPayload needs to give them to us, and we cache them
    until finality to support reorgs without tx losses.
    
Whilst some constraints above might sound overly aggressive, the general idea is
that the blob pool should work robustly for its intended use case and whilst
anyone is free to use blob transactions for arbitrary non-rollup use cases,
they should not be allowed to run amok the network.

Implementation wise there are a few interesting design choices:
  - Adding a transaction to the pool blocks until persisted to disk. This is
    viable because TPS is low (2-4 blobs per block initially, maybe 8-16 at
    peak), so natural churn is a couple MB per block. Replacements doing O(n)
    updates are forbidden and transaction propagation is pull based (i.e. no
    pileup of pending data).
  - When transactions are chosen for inclusion, the primary criteria is the
    signer tip (and having a basefee/data fee high enough of course). However,
    same-tip transactions will be split by their basefee/datafee, preferring
    those that are closer to the current network limits. The idea being that
    very relaxed ones can be included even if the fees go up, when the closer
    ones could already be invalid.
When the pool eventually reaches saturation, some old transactions - that may
never execute - will need to be evicted in favor of newer ones. The eviction
strategy is quite complex:
  - Exceeding capacity evicts the highest-nonce of the account with the lowest
    paying blob transaction anywhere in the pooled nonce-sequence, as that tx
    would be executed the furthest in the future and is thus blocking anything
    after it. The smallest is deliberately not evicted to avoid a nonce-gap.
  - Analogously, if the pool is full, the consideration price of a new tx for
    evicting an old one is the smallest price in the entire nonce-sequence of
    the account. This avoids malicious users DoSing the pool with seemingly
    high paying transactions hidden behind a low-paying blocked one.
  - Since blob transactions have 3 price parameters: execution tip, execution
    fee cap and data fee cap, there's no singular parameter to create a total
    price ordering on. What's more, since the base fee and blob fee can move
    independently of one another, there's no pre-defined way to combine them
    into a stable order either. This leads to a multi-dimensional problem to
    solve after every block.
  - The first observation is that comparing 1559 base fees or 4844 blob fees
    needs to happen in the context of their dynamism. Since these fees jump
    up or down in ~1.125 multipliers (at max) across blocks, comparing fees
    in two transactions should be based on log1.125(fee) to eliminate noise.
  - The second observation is that the basefee and blobfee move independently,
    so there's no way to split mixed txs on their own (A has higher base fee,
    B has higher blob fee). Rather than look at the absolute fees, the useful
    metric is the max time it can take to exceed the transaction's fee caps.
    Specifically, we're interested in the number of jumps needed to go from
    the current fee to the transaction's cap:
    jumps = log1.125(txfee) - log1.125(basefee)
  - The third observation is that the base fee tends to hover around rather
    than swing wildly. The number of jumps needed from the current fee starts
    to get less relevant the higher it is. To remove the noise here too, the
    pool will use log(jumps) as the delta for comparing transactions.
    delta = sign(jumps) * log(abs(jumps))
  - To establish a total order, we need to reduce the dimensionality of the
    two base fees (log jumps) to a single value. The interesting aspect from
    the pool's perspective is how fast will a tx get executable (fees going
    down, crossing the smaller negative jump counter) or non-executable (fees
    going up, crossing the smaller positive jump counter). As such, the pool
    cares only about the min of the two delta values for eviction priority.
    priority = min(deltaBasefee, deltaBlobfee)
  - The above very aggressive dimensionality and noise reduction should result
    in transaction being grouped into a small number of buckets, the further
    the fees the larger the buckets. This is good because it allows us to use
    the miner tip meaningfully as a splitter.
  - For the scenario where the pool does not contain non-executable blob txs
    anymore, it does not make sense to grant a later eviction priority to txs
    with high fee caps since it could enable pool wars. As such, any positive
    priority will be grouped together.
    priority = min(deltaBasefee, deltaBlobfee, 0)
Optimisation tradeoffs:
  - Eviction relies on 3 fee minimums per account (exec tip, exec cap and blob
    cap). Maintaining these values across all transactions from the account is
    problematic as each transaction replacement or inclusion would require a
    rescan of all other transactions to recalculate the minimum. Instead, the
    pool maintains a rolling minimum across the nonce range. Updating all the
    minimums will need to be done only starting at the swapped in/out nonce
    and leading up to the first no-change.