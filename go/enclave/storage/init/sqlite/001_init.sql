create table if not exists keyvalue
(
    ky  varbinary(64) primary key,
    val mediumblob NOT NULL
);

create table if not exists config
(
    ky  varchar(64) primary key,
    val mediumblob NOT NULL
);

insert into config
values ('CURRENT_SEQ', -1);

create table if not exists attestation_key
(
--     party  binary(20) primary key, // todo -pk
    party binary(20),
    ky    binary(33) NOT NULL
);

create table if not exists block
(
    hash         binary(32) primary key,
    parent       binary(32),
    is_canonical boolean NOT NULL,
    header       blob    NOT NULL,
    height       int     NOT NULL
--   the unique constraint is commented for now because there might be multiple non-canonical blocks for the same height
--     unique (height, is_canonical)
);
create index IDX_BLOCK_HEIGHT on block (height);

create table if not exists l1_msg
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    message varbinary(1024) NOT NULL,
    block   binary(32)      NOT NULL REFERENCES block,
    is_transfer boolean
);

create table if not exists rollup
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    start_seq int        NOT NULL,
    end_seq   int        NOT NULL,
    header    blob       NOT NULL,
    block     binary(32) NOT NULL REFERENCES block
);

create table if not exists batch_body
(
    hash    binary(32) primary key,
    content mediumblob NOT NULL
);

create table if not exists batch
(
    hash         binary(32) primary key,
    parent       binary(32),-- REFERENCES batch,
    sequence     int     NOT NULL unique,
    height int NOT NULL,
    is_canonical boolean NOT NULL,
    header blob NOT NULL,
    body binary(32) NOT NULL REFERENCES batch_body,
    l1_proof binary(32) NOT NULL, -- normally this would be a FK, but there is a weird edge case where an L2 node might not have the block used to create this batch
    is_executed  boolean NOT NULL
--   the unique constraint is commented for now because there might be multiple non-canonical batches for the same height
--   unique (height, is_canonical, is_executed)
);
create index IDX_BATCH_HEIGHT on batch (height);
create index IDX_BATCH_SEQ on batch (sequence);
create index IDX_BATCH_Block on batch (l1_proof);

create table if not exists tx
(
    hash           binary(32) primary key,
    content        mediumblob NOT NULL,
    sender_address binary(20) NOT NULL,
    nonce          int        NOT NULL,
    idx            int        NOT NULL,
    body           binary(32) REFERENCES batch_body
);

create table if not exists exec_tx
(
    id                       binary(64) PRIMARY KEY, -- batch_hash||tx_hash
    created_contract_address binary(20),
    receipt                  mediumblob,
--     commenting out the fk until synthetic transactions are also stored
--     tx                       binary(32) REFERENCES tx,
    tx                       binary(32) NOT NULL,
    batch                    binary(32) NOT NULL REFERENCES batch
);
create index IX_EX_TX1 on exec_tx (tx);

create table if not exists events
(
    topic0          binary(32) NOT NULL,
    topic1          binary(32),
    topic2          binary(32),
    topic3          binary(32),
    topic4          binary(32),
    datablob        mediumblob,
    log_idx         int        NOT NULL,
    address         binary(20) NOT NULL,
    lifecycle_event boolean    NOT NULL,
    rel_address1    binary(20),
    rel_address2    binary(20),
    rel_address3    binary(20),
    rel_address4    binary(20),
    exec_tx_id      binary(64) REFERENCES exec_tx
);
create index IX_AD on events (address);
create index IX_RAD1 on events (rel_address1);
create index IX_RAD2 on events (rel_address2);
create index IX_RAD3 on events (rel_address3);
create index IX_RAD4 on events (rel_address4);
create index IX_T0 on events (topic0);
create index IX_T1 on events (topic1);
create index IX_T2 on events (topic2);
create index IX_T3 on events (topic3);
create index IX_T4 on events (topic4);