create table if not exists keyvalue
(
    id  INTEGER PRIMARY KEY AUTOINCREMENT,
    ky  varbinary(64) UNIQUE NOT NULL,
    val mediumblob           NOT NULL
);
create index IDX_KV on keyvalue (ky);

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
    hash         binary(32) NOT NULL,
    is_canonical boolean    NOT NULL,
    header       blob       NOT NULL,
    height       int        NOT NULL
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
    is_transfer boolean         NOT NULL
);
create index L1_MSG_BLOCK_IDX on l1_msg (block);

create table if not exists rollup
(
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    hash              binary(32) NOT NULL,
    start_seq         int        NOT NULL,
    end_seq           int        NOT NULL,
    time_stamp        int        NOT NULL,
    header            blob       NOT NULL,
    compression_block INTEGER    NOT NULL REFERENCES block
);
create index ROLLUP_COMPRESSION_BLOCK_IDX on rollup (compression_block);
create index ROLLUP_COMPRESSION_HASH_IDX on rollup (hash);

create table if not exists batch
(
    sequence       int primary key,
    converted_hash binary(32) NOT NULL,
    hash           binary(32) NOT NULL,
    height         int        NOT NULL,
    is_canonical   boolean    NOT NULL,
    header         blob       NOT NULL,
    l1_proof_hash  binary(32) NOT NULL,
    l1_proof       INTEGER, -- normally this would be a FK, but there is a weird edge case where an L2 node might not have the block used to create this batch
    is_executed    boolean    NOT NULL
    --   the unique constraint is commented for now because there might be multiple non-canonical batches for the same height
--   unique (height, is_canonical, is_executed)
);
create index IDX_BATCH_HASH on batch (hash);
create index IDX_BATCH_BLOCK on batch (l1_proof_hash);
create index IDX_BATCH_L1 on batch (l1_proof);
create index IDX_BATCH_HEIGHT on batch (height);
create index IDX_BATCH_HEIGHT_COMP on batch (is_canonical, is_executed, height);

create table if not exists tx
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    hash           binary(32) NOT NULL,
    content        mediumblob NOT NULL,
    to_address     int,
    type           int8       NOT NULL,
    sender_address int        NOT NULL REFERENCES externally_owned_account,
    idx            int        NOT NULL,
    batch_height   int        NOT NULL,
    is_synthetic   boolean    NOT NULL
);
create index IDX_TX_HASH on tx (hash);
create index IDX_TX_SENDER_ADDRESS on tx (sender_address);
create index IDX_TX_BATCH_HEIGHT on tx (batch_height, idx);

create table if not exists receipt
(
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    post_state               binary(32),
    status                   int     not null,
    cumulative_gas_used      int     not null,
    effective_gas_price      int,
    created_contract_address binary(20),
    --     commenting out the fk until synthetic transactions are also stored
    tx                       INTEGER,
    batch                    INTEGER NOT NULL REFERENCES batch
);
create index IDX_EX_TX_BATCH on receipt (batch);
create index IDX_EX_TX_CCA on receipt (tx);

create table if not exists contract
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    address         binary(20) NOT NULL,
    --     the tx signer that created the contract
    creator         int        NOT NULL REFERENCES externally_owned_account,
    auto_visibility boolean    NOT NULL,
    transparent     boolean,
    tx              INTEGER    NOT NULL REFERENCES tx
);
create index IDX_CONTRACT_AD on contract (address);

create table if not exists externally_owned_account
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    address binary(20) NOT NULL
);
create index IDX_EOA on externally_owned_account (address);

-- not very large. An entry for every event_type
create table if not exists event_type
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    contract        INTEGER    NOT NULL references contract,
    event_sig       binary(32) NOT NULL,
    auto_visibility boolean    NOT NULL,
    auto_public     boolean,
    config_public   boolean    NOT NULL,
    topic1_can_view boolean,
    topic2_can_view boolean,
    topic3_can_view boolean,
    sender_can_view boolean
);
create index IDX_EV_CONTRACT on event_type (contract, event_sig);

--  very large table with user values
create table if not exists event_topic
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    event_type  INTEGER references event_type,
    topic       binary(32) NOT NULL,
    rel_address INTEGER references externally_owned_account
--    pos         INTEGER    NOT NULL -- todo
);
create index IDX_TOP on event_topic (topic);
create index IDX_REL_A on event_topic (rel_address);

create table if not exists event_log
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    event_type INTEGER NOT NULL references event_type,
    topic1     INTEGER references event_topic,
    topic2     INTEGER references event_topic,
    topic3     INTEGER references event_topic,
    datablob   mediumblob,
    log_idx    INTEGER NOT NULL,
    receipt    INTEGER NOT NULL references receipt
);
create index IDX_EV on event_log (receipt, event_type, topic1, topic2, topic3);

