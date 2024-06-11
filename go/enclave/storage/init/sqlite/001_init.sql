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

create table if not exists batch_body
(
    id      int        NOT NULL primary key,
    content mediumblob NOT NULL
);

create table if not exists batch
(
    sequence       int primary key,
    converted_hash binary(32) NOT NULL,
    hash           binary(32) NOT NULL,
    height         int        NOT NULL,
    is_canonical   boolean    NOT NULL,
    header         blob       NOT NULL,
    body           int        NOT NULL REFERENCES batch_body,
    l1_proof_hash  binary(32) NOT NULL,
    l1_proof       INTEGER, -- normally this would be a FK, but there is a weird edge case where an L2 node might not have the block used to create this batch
    is_executed    boolean    NOT NULL
    --   the unique constraint is commented for now because there might be multiple non-canonical batches for the same height
--   unique (height, is_canonical, is_executed)
);
create index IDX_BATCH_HASH on batch (hash);
create index IDX_BATCH_BLOCK on batch (l1_proof_hash);
create index IDX_BATCH_BODY on batch (body);
create index IDX_BATCH_L1 on batch (l1_proof);
create index IDX_BATCH_HEIGHT on batch (height);

create table if not exists tx
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    hash           binary(32) NOT NULL,
    content        mediumblob NOT NULL,
    sender_address binary(20) NOT NULL,
    nonce          int        NOT NULL,
    idx            int        NOT NULL,
    body           int        NOT NULL REFERENCES batch_body
);
create index IDX_TX_HASH on tx (hash);
create index IDX_TX_SENDER_ADDRESS on tx (sender_address);

create table if not exists exec_tx
(
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    created_contract_address INTEGER REFERENCES contract,
    receipt                  mediumblob,
    --     commenting out the fk until synthetic transactions are also stored
    tx                       INTEGER,
    batch                    INTEGER NOT NULL REFERENCES batch
);
create index IDX_EX_TX_BATCH on exec_tx (batch);
create index IDX_EX_TX_CCA on exec_tx (created_contract_address, tx);

create table if not exists contract
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    address binary(20) NOT NULL
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
    event_sig       binary(32) NOT NULL, -- no need to index because there are only a few events for an address
    lifecycle_event boolean    NOT NULL  -- set based on the first event, and then updated to false if it turns out it is true
);
create index IDX_EV_CONTRACT on event_type (contract, event_sig);

--  very large table with user values
create table if not exists event_topic
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    topic       binary(32) NOT NULL,
    rel_address INTEGER references externally_owned_account
--    pos         INTEGER    NOT NULL -- todo
);
-- create index IDX_TOP on event_topic (topic, pos);
create index IDX_TOP on event_topic (topic);

create table if not exists event_log
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    event_type INTEGER NOT NULL references event_type,
    topic1     INTEGER references event_topic,
    topic2     INTEGER references event_topic,
    topic3     INTEGER references event_topic,
    datablob   mediumblob,
    log_idx    INTEGER NOT NULL,
    exec_tx    INTEGER NOT NULL references exec_tx
);
-- create index IDX_BATCH_TX on event_log (exec_tx);
create index IDX_EV on event_log (exec_tx, event_type, topic1, topic2, topic3);

-- requester - address
-- exec_tx - range of batch heights or a single batch
-- address []list of contract addresses
-- topic0 - event sig   []list
-- topic1    []list
-- topic2    []list
-- topic3    []list


-- select * from event_log
--          join exec_tx on exec_tx
--              join batch on exec_tx.batch -- to get the batch height range
--          join event_type ec on event_type
--              join contract c  on
--          left join event_topic t1 on topic1
--              left join externally_owned_account eoa1 on t1.rel_address
--          left join event_topic t2 on topic2
--              left join externally_owned_account eoa2 on t2.rel_address
--          left join event_topic t3 on topic3
--              left join externally_owned_account eoa3 on t3.rel_address
-- where
--  exec_tx.
--  c.address in [address..] AND
--  ec.event_sig in [topic0..] AND
--  t1.topic in [topic1..] AND
--  t2.topic in [topic2..] AND
--  t3.topic in [topic3..] AND
--  b.height in [] and b.is_canonical=true
--  (ec.lifecycle_event OR eoa1.address=requester OR eoa2.address=requester OR eoa3.address=requester)
