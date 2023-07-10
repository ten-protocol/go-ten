create table if not exists keyvalue
(
    ky  varbinary(64) primary key,
    val mediumblob
);

create table if not exists config
(
    key  varchar(64) primary key,
    val mediumblob
);

create table if not exists attestation_key
(
--     party  binary(20) primary key, // todo -pk
    party  binary(20),
    key binary(33)
);

create table if not exists block
(
    hash         binary(32) primary key,
    parent       binary(32) REFERENCES block,
    is_canonical boolean,
    header       text,
    height       int
);
create index IDX_BLOCK_HEIGHT on block (height);

create table if not exists l1_msg
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    message varbinary(1024),
    block   binary(32) REFERENCES block
);

create table if not exists rollup
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    start_seq int,
    end_seq   int,
    block     binary(32) REFERENCES block
);

create table if not exists batch_body
(
    hash    binary(32) primary key,
    content mediumblob
);

create table if not exists batch
(
    hash         binary(32) primary key,
    parent       binary(32) REFERENCES batch,
    sequence     int,
    height       int,
    is_canonical boolean,
    header       text,
    body         binary(32) REFERENCES batch_body,
    l1Proof      binary(32) REFERENCES block,
    source       text
);
create index IDX_BATCH_HEIGHT on batch (height);
create index IDX_BATCH_SEQ on batch (sequence);

create table if not exists tx
(
    hash          binary(32) primary key,
    content       mediumblob,
    senderAddress binary(20),
    nonce         int,
    idx           int,
    body          binary(32) REFERENCES batch_body
);

create table if not exists exec_tx
(
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    createdContractAddress binary(32),
    receipt                mediumblob,
    tx                     binary(32) REFERENCES tx,
    batch                  binary(32) REFERENCES batch
);

-- todo - remove some fields
create table if not exists events
(
    topic0         binary(32),
    topic1         binary(32),
    topic2         binary(32),
    topic3         binary(32),
    topic4         binary(32),
    datablob       mediumblob,
    blockHash      binary(32),
    blockNumber    int,
    txHash         binary(32),
    txIdx          int,
    logIdx         int,
    address        binary(32),
    lifecycleEvent boolean,
    relAddress1    binary(20),
    relAddress2    binary(20),
    relAddress3    binary(20),
    relAddress4    binary(20)
);
create index IX_AD on events (address);
create index IX_RAD1 on events (relAddress1);
create index IX_RAD2 on events (relAddress2);
create index IX_RAD3 on events (relAddress3);
create index IX_RAD4 on events (relAddress4);
create index IX_BLH on events (blockHash);
create index IX_BLN on events (blockNumber);
create index IX_TXH on events (txHash);
create index IX_T0 on events (topic0);
create index IX_T1 on events (topic1);
create index IX_T2 on events (topic2);
create index IX_T3 on events (topic3);
create index IX_T4 on events (topic4);