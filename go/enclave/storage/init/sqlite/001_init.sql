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
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    hash         binary(4),
    full_hash    binary(32),
    is_canonical boolean NOT NULL,
    header       blob    NOT NULL,
    height       int     NOT NULL
    --   the unique constraint is commented for now because there might be multiple non-canonical blocks for the same height
--     unique (height, is_canonical)
);
create index IDX_BLOCK_HEIGHT on block (height);
create index IDX_BLOCK_HASH on block (hash);

create table if not exists l1_msg
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    message     varbinary(1024) NOT NULL,
    block       INTEGER         NOT NULL REFERENCES block,
    is_transfer boolean
);

create table if not exists rollup
(
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    hash              binary(4),
    full_hash         binary(32),
    start_seq         int     NOT NULL,
    end_seq           int     NOT NULL,
    time_stamp        int     NOT NULL,
    header            blob    NOT NULL,
    compression_block INTEGER NOT NULL REFERENCES block
);

create table if not exists batch_body
(
    id      int        NOT NULL primary key,
    content mediumblob NOT NULL
);

create table if not exists batch
(
    sequence       int primary key,
    full_hash      binary(32),
    converted_hash binary(32),
    hash           binary(4) NOT NULL,
    height         int       NOT NULL,
    is_canonical   boolean   NOT NULL,
    header         blob      NOT NULL,
    body           int       NOT NULL REFERENCES batch_body,
    l1_proof       INTEGER   NOT NULL, -- normally this would be a FK, but there is a weird edge case where an L2 node might not have the block used to create this batch
    is_executed    boolean   NOT NULL
    --   the unique constraint is commented for now because there might be multiple non-canonical batches for the same height
--   unique (height, is_canonical, is_executed)
);
create index IDX_BATCH_HASH on batch (hash);
create index IDX_BATCH_HEIGHT on batch (height, is_canonical);
create index IDX_BATCH_Block on batch (l1_proof);

create table if not exists tx
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    hash           binary(4),
    full_hash      binary(32) NOT NULL,
    content        mediumblob NOT NULL,
    sender_address binary(20) NOT NULL,
    nonce          int        NOT NULL,
    idx            int        NOT NULL,
    body           int REFERENCES batch_body
);
create index IDX_Tx_HASH on tx (hash);

create table if not exists exec_tx
(
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    created_contract_address binary(20),
    receipt                  mediumblob,
    --     commenting out the fk until synthetic transactions are also stored
--     tx                       binary(16) REFERENCES tx,
    tx                       INTEGER NOT NULL,
    batch                    INTEGER NOT NULL REFERENCES batch
);
create index IX_EX_TX1 on exec_tx (tx);

-- todo denormalize. Extract contract and user table and point topic0 and rel_addreses to it
create table if not exists events
(
    topic0            binary(4)  NOT NULL,
    topic1            binary(4),
    topic2            binary(4),
    topic3            binary(4),
    topic4            binary(4),
    topic0_full       binary(32) NOT NULL,
    topic1_full       binary(32),
    topic2_full       binary(32),
    topic3_full       binary(32),
    topic4_full       binary(32),
    datablob          mediumblob,
    log_idx           int        NOT NULL,
    address           binary(4)  NOT NULL,
    address_full      binary(20) NOT NULL,
    lifecycle_event   boolean    NOT NULL,
    rel_address1      binary(4),
    rel_address2      binary(4),
    rel_address3      binary(4),
    rel_address4      binary(4),
    rel_address1_full binary(20),
    rel_address2_full binary(20),
    rel_address3_full binary(20),
    rel_address4_full binary(20),
    tx                INTEGER    NOT NULL,
    batch             INTEGER    NOT NULL REFERENCES batch
);
create index IDX_AD on events (address);
create index IDX_RAD1 on events (rel_address1);
create index IDX_RAD2 on events (rel_address2);
create index IDX_RAD3 on events (rel_address3);
create index IDX_RAD4 on events (rel_address4);
create index IDX_T0 on events (topic0);
create index IDX_T1 on events (topic1);
create index IDX_T2 on events (topic2);
create index IDX_T3 on events (topic3);
create index IDX_T4 on events (topic4);