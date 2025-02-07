# Failure Mode when L1 Transactions get stuck

## Current Problem

There is hardcoded limit of 16 txs in the blobpool. The blobpool is grouped by account and if enough of our rollups fail
during L1 publishing then this pool gets full and we need to evict some of the old transactions so we can continue 
publishing rollups. This would be a pretty catastrophic scenario to get into on mainnet so we need to try and add some safeguards
to prevent that from happening and make it easier to fix if it does happen. 

The current way of resolving is incredibly tedious and time-consuming which will mean we're stuck unable to post rollups 
for a while on mainnet. It would require us finding the tx hash from the logs, finding the corresponding rollup hash 
and crafting a transaction manually with a higher gas price than the lowest one stuck in the pool.

There are many reasons why a rollup tx may fail but another issue we have is we will just keep retrying it and if there's 
there's a gas anomaly we could easily fill up the blobpool and be forced to remediate.  

## Proposed Solution
* Had a max retries limit for publishing a tx 
* Persist any tx details when the retry limit is hit so any stuck txs can be easily found/ retried manually
