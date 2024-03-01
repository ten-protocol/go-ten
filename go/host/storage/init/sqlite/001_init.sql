create table if not exists config
(
    ky  varchar(64) primary key,
    val mediumblob NOT NULL
);

insert into config
values ('CURRENT_SEQ', -1);

create table if not exists rollup
(
    id                int PRIMARY KEY AUTO_INCREMENT,
    hash              binary(16) NOT NULL UNIQUE,
    start_seq         int        NOT NULL,
    end_seq           int        NOT NULL,
    time_stamp        int        NOT NULL,
    header            blob       NOT NULL,
    compression_block binary(32) NOT NULL
);

create index IDX_ROLLUP_HASH on rollup (hash);
create index IDX_ROLLUP_PROOF on rollup (compression_block);
create index IDX_ROLLUP_SEQ on rollup (start_seq, end_seq);

create table if not exists batch_body
(
    id          int        NOT NULL primary key,
    content     mediumblob NOT NULL
);

create table if not exists batch
(
    sequenceOrder       int primary key,
    full_hash      binary(32) NOT NULL,
    hash           binary(16) NOT NULL unique,
    height         int        NOT NULL,
    tx_count       int        NOT NULL,
    header         blob       NOT NULL,
    body_id        int        NOT NULL REFERENCES batch_body
    );
create index IDX_BATCH_HASH on batch (hash);
create index IDX_BATCH_HEIGHT on batch (height);

create table if not exists transaction_count
(
    id          int  NOT NULL primary key,
    count       int  NOT NULL
);
