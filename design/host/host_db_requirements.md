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
create table if not exists rollup_host
(
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    hash              binary(16) NOT NULL UNIQUE,
    start_seq         int        NOT NULL,
    end_seq           int        NOT NULL,
    time_stamp        int        NOT NULL,
    ext_rollup        blob       NOT NULL,
    compression_block binary(32) NOT NULL
);

create index IDX_ROLLUP_HASH_HOST on rollup_host (hash);
create index IDX_ROLLUP_PROOF_HOST on rollup_host (compression_block);
create index IDX_ROLLUP_SEQ_HOST on rollup_host (start_seq, end_seq);
```

Calculating the `L1BlockHeight` as done in `calculateL1HeightsFromDeltas` will be quite computationally expensive so we 
can just order them by `end_seq`.

### Batch 
Storing the encoded ext batch so that we can provide rich data to the UI including gas, receipt, cross-chain hash etc.
```sql
create table if not exists batch_host
(
    sequence       int primary key,
    full_hash      binary(32) NOT NULL,
    hash           binary(16) NOT NULL unique,
    height         int        NOT NULL,
    ext_batch      mediumblob NOT NULL
    );
    
create index IDX_BATCH_HEIGHT_HOST on batch_host (height);

```

### Transactions 

We need to store these separately for efficient lookup of the batch by tx hash and vice versa. 

Because we are able to decrypt the encrypted blob on testnet we are able to retrieve the number of transactions that way
but on mainnet this won't be possible, so we need to store the `tx_count` in this table. There is a plan to remove
`ExtBatch.TxHashes` and expose a new Enclave API to retrieve this.

```sql
create table if not exists transactions_host
(
    hash           binary(32) primary key,
    b_sequence     int REFERENCES batch_host
);

create table if not exists transaction_count
(
    id          int  NOT NULL primary key,
    total       int  NOT NULL
);

```

## Database Choice

The obvious choice is MariaDB as this is what is used by the gateway so we would have consistency across the stack. It 
would make deployment simpler as the scripts are already there. Main benefits of MariaDB:

* Offer performance improvements through the use of aria storage engine which is not available through MySQL
* Strong security focus with RBAC and data-at-rest encryption 
* Supports a large number of concurrent connections 

Postgres would be the obvious other alternative but given it is favoured for advanced data types, complex queries and 
geospatial capabilities, it doesn't offer us any benefit for this use case over MariaDB.

## Cross Chain Messages

We want to display L2 > L1 and L1 > L2 transaction data. We will expose an API to retrieve these and the implementation
for retrieving the data will either be via subscriptions to the events API or we will store them in the database. TBC   