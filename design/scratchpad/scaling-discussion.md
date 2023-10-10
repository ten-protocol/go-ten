# Scaling Obscuro enclave

Some notes for discussion and comment about how we ensure the enclave database is performant for the long-term.

### Problem to solve:
- Enclave database performance is a crucial bottleneck for the smooth running of the network
- It is worse than typical database scaling problems because of the constraints of running a database in an SGX enclave
- In particular the batch table is currently an issue, we create 1 row/second and the table is on the critical path for batch/rollup production
- We can use indexes and clever tricks (like skipping batch production when no activity) but eventually we'll start bumping up against limits again

Approaches considered:
1. Indexes and query optimisation
2. Avoid producing empty batches
3. Discarding historic data from enclave


### 1. Indexes and query optimisation
- I'm not an expert and not up to speed with what we've done on this, but I assume eventually the reads and or writes will start to become a serious bottleneck regardless of optimisations and hardware

### 2. Avoid producing empty batches
Only produce a batch if there will be a tx in it or if 1 minute has passed since prev batch.
- This could reduce number of rows in the batch table drastically while we're in a bootstrapping phase
- This is only helpful in the short term, once we get close to 1 tx/sec with traction or automated processes then this no longer helps
- Complexity to keep a sensible DX:
    * some contracts use block height as a clock, expecting them to be fairly regular
    * if we fudge the batch height to respect the 1 batch/sec then we break the assumption that `batch(height - 1) == parentBatch` which could be an even more painful gotcha than breaking the regular interval

### 3. Discarding historic data from enclave
- We use ethereum calldata to store all batch and transaction data, the enclave DB is not the source of truth once batches have been rolled up
- AFAICT there's no reason we can't start deleting data from the enclave after some time has passed to keep the db snappy
- For the sequencer this seems extremely sensible, it needs to be very performant to produce batches and rollups reliably and we don't need it to service historic data requests
- For validators this is a little more nuanced, we need to ensure we can always service historic requests:
    * if the enclave doesn't have the data it could replay rollup transactions against a statedDB snapshot from just before the rollup (VERY expensive but worst case scenario)
    * it would be possible to have enclaves that retained different periods of recent history so the OG could send requests to the appropriate ones (with expensive replay as the worst case scenario)


#### Some concerns:
- how does this change when EIP-4844 (blobspace/proto-danksharding) goes live? We'll need another data storage solution (it could just be the HA databases on the hosts I guess?)
- is the stuff in geth around snapshots and archiving re-usable or not really relevant? Don't want to naively reimplement stuff.

#### Some gaps in my knowledge:
- How urgent is this discussion and work? Have we staved off perf issues for months with the latest fixes so have bigger fish to fry?
- How are the stateDB caches stored and managed?
    * are they relevant to enclave performance? Do they take a large amount of enclave memory or a lot of DB persistence?
    * can we control stateDB snapshots for important batches (like the one before each rollup) while purging other historical entries?