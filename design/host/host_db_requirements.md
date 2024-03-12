# Moving Host DB to SQL 

The current implementation uses the `ethdb.KeyValueStore` which provides fast access but is not sufficient for the 
querying capabilities required by Tenscan. We want to move to an SQL implementation similar to what the Enclave uses.

## Current Storage 
### Schema Keys
```go
var (
	blockHeaderPrefix       = []byte("b")
	blockNumberHeaderPrefix = []byte("bnh")
	batchHeaderPrefix       = []byte("ba")
	batchHashPrefix         = []byte("bh")
	batchNumberPrefix       = []byte("bn")
	batchPrefix             = []byte("bp")
	batchHashForSeqNoPrefix = []byte("bs")
	batchTxHashesPrefix     = []byte("bt")
	headBatch               = []byte("hb")
	totalTransactionsKey    = []byte("t")
	rollupHeaderPrefix      = []byte("rh")
	rollupHeaderBlockPrefix = []byte("rhb")
	tipRollupHash           = []byte("tr")
	blockHeadedAtTip        = []byte("bht")
)
```
Some of the schema keys are dummy keys for entries where we only have one entry that is updated such as totals or tip 
data. The rest of the schema keys are used as prefixes appended with the `byte[]` representation of the key.

| Data Type        | Description                     | Schema | Key                          | Value (Encoded)    |
|------------------|---------------------------------|--------|------------------------------|--------------------|
| **Batch**        | Batch hash to headers           | ba     | BatchHeader.Hash()           | BatchHeader        |
| **Batch**        | Batch hash to ExtBatch          | bp     | ExtBatch.Hash()              | ExtBatch           |
| **Batch**        | Batch hash to TX hashes         | bt     | ExtBatch.Hash()              | ExtBatch.TxHashes  |
| **Batch**        | Batch number to batch hash      | bh     | BatchHeader.Number           | BatchHeader.Hash() |
| **Batch**        | Batch seq no to batch hash      | bs     | BatchHeader.SequencerOrderNo | BatchHeader.Hash() |
| **Batch**        | TX hash to batch number         | bn     | ExtBatch.TxHashes[i]         | BatchHeader.Number |
| **Batch**        | Head Batch                      | hb     | "hb"                         | ExtBatch.Hash()    |
| **Block**        | L1 Block hash to block header   | b      | Header.Hash()                | Header             |
| **Block**        | L1 Block height to block header | bnh    | Header.Number                | Header             |
| **Block**        | Latest Block                    | bht    | "bht"                        | Header.Hash()      |
| **Rollup**       | Rollup hash to header           | rh     | RollupHeader.Hash()          | RollupHeader       |
| **Rollup**       | L1 Block hash to rollup header  | rhb    | L1Block.Hash()               | RollupHeader       |
| **Rollup**       | Tip rollup header               | tr     | "tr"                         | RollupHeader       |
| **Transactions** | Total number of transactions    | t      | "t"                          | Int                |

## Tenscan Functionality Requirements

### Mainnet Features 
#### Currently supported 
* Return the list of batches in descending order 
* View details within the batch (BatchHeader and ExtBatch)
* Return the number of transactions within the batch
* Return the list of transactions in descending order

### Not currently supported
* Return a list of rollups in descending order 
* View details of the rollup (probably needs to be ExtBatch for user )
* Navigate to the L1 block on etherscan from the rollup
* Return the list of batches within the rollup 
* Navigate from the transaction to the batch it was included in
* Navigate from the batch to the rollup that it was included in
* TODO Cross chain messaging - Arbiscan shows L1>L2 and L2>L1 

### Testnet-Only  Features
#### Currently supported
* Copy the encrypted TX blob to a new page and decrypt there

#### Not currently supported
* From the batch you should be able to optionally decrypt the transactions within the batch 
* Navigate into the transaction details from the decrypted transaction  
* We want to be able to navigate up the chain from TX to batch to rollup

## SQL Schema

There are some considerations here around the behaviour of tenscan for testnet vs mainnet. Because we are able to decrypt 
the encrypted blob on testnet we are able to retrieve the number of transactions but on mainnet this wont be possible so 
we need to store the TxCount in 

### Rollup
```sql
create table if not exists rollup
(
    hash              binary(16) primary key,
    start_seq         int        NOT NULL,
    end_seq           int        NOT NULL,
    started_at        int        NOT NULL,
    compression_block binary(16) NOT NULL
);

create index IDX_ROLLUP_PROOF on rollup (compression_block);
create index IDX_ROLLUP_SEQ on rollup (start_seq, end_seq);
```

Calculating the `L1BlockHeight` as done in `calculateL1HeightsFromDeltas` will be quite computationally expensive so we 
can just order them by `end_seq`.

### Batch 
```sql
create table if not exists batch
(
    sequence       int primary key,
    hash           binary(16) NOT NULL unique,
    height         int        NOT NULL,
    tx_count       int        NOT NULL,
    header         blob       NOT NULL,
    body           int        NOT NULL REFERENCES batch_body
    );
create index IDX_BATCH_HASH on batch (hash);
create index IDX_BATCH_HEIGHT on batch (height);
```

Because we are able to decrypt the encrypted blob on testnet we are able to retrieve the number of transactions that way
 but on mainnet this won't be possible, so we need to store the `tx_count` in this table. There is a plan to remove 
 `ExtBatch.TxHashes` and expose a new Enclave API to retrieve this. 

We don't need to store the TX data since it's provided by the `EnclaveClient` which returns a `PublicTransaction`
containing the `BatchHeight` which can be used to navigate to the batch. The `height` is stored outside the header to 
provide an efficient lookup otherwise we'd need to decode the batch header and find it that way. 

Storing the encoded batch header so that we can provide rich data to the UI including gas, receipt, cross-chain hash etc.   

### BatchBody
```sql
create table if not exists batch_body
(
    id          int        NOT NULL primary key,
    content     mediumblob NOT NULL
);

```

Splitting this out means we only fetch when a user inspects a batch which will reduce the overhead on the table pages. 


### Queries
These will be different for the actual DB implementation but just as a representation for how the functional requirements 
can be solved.  

List all batches allowing for pagination 
```sql 
SELECT * FROM batch ORDER BY height DESC
``` 

Return the list of transactions within the batch (FE can decrypt and count the results)
```sql 
SELECT content FROM batch_body WHERE id =?
```

Return a list of rollups in descending order
```sql 
SELECT * FROM rollup ORDER BY end_seq DESC
```

Return list of batches in the rollup
```sql
  SELECT b.* 
  FROM batch b JOIN rollup r ON b.height BETWEEN r.start_seq AND r.end_seq
  WHERE r.hash = ?
```
Navigate from to the rollup from the batch
```sql
SELECT r.hash
FROM rollup r
WHERE ? BETWEEN r.start_seq AND r.end_seq;
```

## Database Choice

The obvious choice is MariaDB as this is what is used by the gateway so we would have consistency across the stack. It 
would make deployment simpler as the scripts are already there. Main benefits of MariaDB:

* Offer performance improvements through the use of aria storage engine which is not available through MySQL
* Strong security focus with RBAC and data-at-rest encryption 
* Supports a large number of concurrent connections 

Postgres would be the obvious other alternative but given it is favoured for advanced data types, complex queries and 
geospatial capabilities, it doesn't offer us any benefit for this use case over MariaDB. 