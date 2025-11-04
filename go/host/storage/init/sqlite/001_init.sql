create table if not exists config
(
    ky  varchar(64) primary key,
    val mediumblob NOT NULL
    );

insert into config
values ('CURRENT_SEQ', -1);

create table if not exists block_host
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    hash        binary(32)      NOT NULL UNIQUE,
    header      blob            NOT NULL
);

create index IDX_BLOCK_HASH_HOST on block_host (hash);

create table if not exists rollup_host
(
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    hash              binary(32) NOT NULL UNIQUE,
    start_seq         int        NOT NULL,
    end_seq           int        NOT NULL,
    time_stamp        int        NOT NULL,
    ext_rollup        blob       NOT NULL,
    compression_block int NOT NULL references block_host
);

create table if not exists cross_chain_message_host
(
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    message_hash      binary(32) NOT NULL UNIQUE,
    rollup_id         int NOT NULL references rollup_host,
    message_type      char(1) NOT NULL CHECK (message_type IN ('m', 'v'))
);

create index IDX_CCM_HASH_HOST on cross_chain_message_host (message_hash);
create index IDX_CCM_ROLLUP_HOST on cross_chain_message_host (rollup_id);


create index IDX_ROLLUP_HASH_HOST on rollup_host (hash);
create index IDX_ROLLUP_PROOF_HOST on rollup_host (compression_block);
create index IDX_ROLLUP_SEQ_HOST on rollup_host (start_seq, end_seq);

create table if not exists batch_host
(
    sequence       int primary key,
    hash           binary(32) NOT NULL,
    height         int        NOT NULL,
    ext_batch      mediumblob NOT NULL,
    txs_size       int        NOT NULL
);
create index IDX_BATCH_HASH_HOST on batch_host (hash);
create index IDX_BATCH_HEIGHT_HOST on batch_host (height);

create table if not exists transaction_host
(
    id             int   PRIMARY KEY,
    hash           binary(32) ,
    b_sequence     int REFERENCES batch_host
);
create index TX_HASH_HOST on transaction_host (hash);

create table if not exists transaction_count
(
    id          int  NOT NULL PRIMARY KEY,
    total       int  NOT NULL
);

insert into transaction_count (id, total)
values (1, 0) on CONFLICT (id) DO NOTHING;

create table if not exists sequencer_attestation_host
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    enclave_id binary(20) UNIQUE NOT NULL,
    is_active  boolean           NOT NULL DEFAULT 1
);

create index IDX_SEQ_ATT_ENCLAVE_ID_HOST on sequencer_attestation_host (enclave_id);