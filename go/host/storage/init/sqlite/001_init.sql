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

create index IDX_ROLLUP_HASH_HOST on rollup_host (hash);
create index IDX_ROLLUP_PROOF_HOST on rollup_host (compression_block);
create index IDX_ROLLUP_SEQ_HOST on rollup_host (start_seq, end_seq);

create table if not exists batch_host
(
    sequence       int primary key,
    hash           binary(32) NOT NULL,
    height         int        NOT NULL,
    ext_batch      mediumblob NOT NULL
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