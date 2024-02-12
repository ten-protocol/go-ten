# Moving Host DB to SQL 

The current implementation uses the `ethdb.KeyValueStore` which provides fast access but is not sufficient for the 
querying capabilities required by Tenscan. We want to move to an SQL implementation similar to what the Enclave uses.

## Current Storage 
| Data Type        | Key             | Value                                     |
|------------------|-----------------|-------------------------------------------|
| **Batch**        | Batch Hash      | Batch Header                              |
| **Batch**        | Batch Hash      | Transaction Hashes                        |
| **Batch**        | Batch Hash      | Batch Data                                |
| **Batch**        | Batch Hash      | Transaction hashes belonging to the Batch |
| **Block**        | L1 Block Hash   | L1 Block Header                           |
| **Block**        | L1 Block Number | L1 Block Header                           |
| **Rollup**       | L1 Block Hash   | Rollup Header                             |
| **Rollup**       | Rollup Hash     | Rollup Header                             |
| **Transactions** | Transactions    | Total number of transactions              |

## Tenscan Functionality Requirements

### Supported
* Return the list of batches in descending order 
* Return the list of transactions within the batch (once decrypted)
* Return the list of transactions in descending order 
* Decrypt the encrypted TX blob

### Not supported
* Return a list of rollups in descending order
* Navigate to the L1 block on etherscan from the rollup
* Return the list of batches within the rollup
* Navigate from the transaction to the batch it was included in
* Navigate from the batch to the rollup that it was included in

## SQL Schema

There are some considerations here around the behaviour of tenscan for testnet vs mainnet. Because we are able to decrypt 
the encrypted blob on testnet we are able to retrieve the number of transactions but on mainnet this wont be possible so 
we need to store the TxCount in 

### Rollup 
* `rollupHash` PK
* `FirstBatchSeqNo`
* `LastBatchSeqNo`
* `L1BlockHash` 

Calculating the `L1BlockHeight` as done in `calculateL1HeightsFromDeltas` will be quite computationally expensive so we 
can just order them by `LastBatchSeqNo`. I don't see any requirements for us to store the encoded headers since this is 
the only information we actually need to link to the batches/ L1 blocks. 

### Batch 
* `SequencerOrderNo` PK
* `batchHash`
* `Number`
* `TxCount`
* `BatchHeader`

Because we are able to decrypt the encrypted blob on testnet we are able to retrieve the number of transactions that way
 but on mainnet this won't be possible so we need to store the TxCount in this table. There is a plan to remove 
 `ExtBatch.TxHashes` and expose a new Enclave API to retrieve this. 

We don't need to store the TX data since it's provided by the `EnclaveClient` which returns a `PublicTransaction`
containing the `BatchHeight` which can be used to navigate to the batch. The `Number` is stored to provide an efficient 
lookup otherwise we'd need to decode the batch header and find it that way. 

### BatchBody
* `batchHash` PK
* `EncryptedTxBlob`

Splitting this out means we only fetch when a user inspects a batch which will reduce the overhead on the table pages. 


### Queries
These will be different for the actual DB implementation but just as a representation for how the functional requirements 
can be solved.  

List all batches allowing for pagination 
```sql 
SELECT * FROM batch ORDER BY SequencerOrderNo DESC
``` 

Return the list of transactions within the batch (FE can decrypt and count the results)
```sql 
SELECT EncryptedTxBlob FROM batch_body WHERE batchHash =?
```

Return a list of rollups in descending order
```sql 
SELECT * FROM rollup ORDER BY LastBatchSeqNo DESC
```

Return list of batches in the rollup
```sql
  SELECT b.* 
  FROM batch b JOIN rollup r ON b.SequencerOrderNo BETWEEN r.FirstBatchSeqNo AND r.LastBatchSeqNo
  WHERE r.rollupHash = ?
```
Navigate from the batch to the rollup it was included in
```sql
SELECT r.rollupHash
FROM rollup r
WHERE ? BETWEEN r.FirstBatchSeqNo AND r.LastBatchSeqNo;
```
