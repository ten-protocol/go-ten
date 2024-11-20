CREATE USER "ten" REQUIRE ISSUER "/CN=tenCA" SUBJECT "/CN=tenUser";
CREATE DATABASE tendb;
GRANT ALL ON tendb.* TO ten;

create table if not exists tendb.keyvalue
(
    id  INTEGER AUTO_INCREMENT,
    ky  varbinary(64) NOT NULL,
    val mediumblob    NOT NULL,
    primary key (id),
    INDEX USING HASH (ky)
);

create table if not exists tendb.config
(
    ky  varchar(64),
    val mediumblob NOT NULL,
    primary key (ky)
);

insert into tendb.config
values ('CURRENT_SEQ', -1);

create table if not exists tendb.attestation_key
(
    party binary(20),
    ky    binary(33) NOT NULL
);

create table if not exists tendb.block
(
    id           INTEGER AUTO_INCREMENT,
    hash         binary(32) NOT NULL,
    is_canonical boolean    NOT NULL,
    header       blob       NOT NULL,
    height       int        NOT NULL,
    primary key (id),
    INDEX (height),
    INDEX USING HASH (hash)
);

create table if not exists tendb.l1_msg
(
    id          INTEGER AUTO_INCREMENT,
    message     varbinary(1024) NOT NULL,
    block       INTEGER         NOT NULL,
    is_transfer boolean         NOT NULL,
    INDEX (block),
    primary key (id)
);

create table if not exists tendb.rollup
(
    id                INTEGER AUTO_INCREMENT,
    hash              binary(32) NOT NULL,
    start_seq         int        NOT NULL,
    end_seq           int        NOT NULL,
    time_stamp        int        NOT NULL,
    header            blob       NOT NULL,
    compression_block INTEGER    NOT NULL,
    INDEX (compression_block),
    INDEX USING HASH (hash),
    primary key (id)
);

create table if not exists tendb.batch
(
    sequence       INTEGER,
    converted_hash binary(32) NOT NULL,
    hash           binary(32) NOT NULL,
    height         int        NOT NULL,
    is_canonical   boolean    NOT NULL,
    header         blob       NOT NULL,
    l1_proof_hash  binary(32) NOT NULL,
    l1_proof       INTEGER,
    is_executed    boolean    NOT NULL,
    primary key (sequence),
    INDEX USING HASH (hash),
    INDEX USING HASH (l1_proof_hash),
    INDEX (l1_proof),
    INDEX (height),
    INDEX (is_canonical, is_executed, height)
);

create table if not exists tendb.tx
(
    id             INTEGER AUTO_INCREMENT,
    hash           binary(32) NOT NULL,
    content        mediumblob NOT NULL,
    to_address     int,
    type           SMALLINT   NOT NULL,
    sender_address int        NOT NULL,
    idx            int        NOT NULL,
    batch_height   int        NOT NULL,
    is_synthetic   boolean    NOT NULL,
    INDEX USING HASH (hash),
    INDEX (sender_address),
    INDEX (batch_height, idx),
    primary key (id)
);

create table if not exists tendb.receipt
(
    id                       INTEGER AUTO_INCREMENT,
    post_state               binary(32),
    status                   TINYINT not null,
    cumulative_gas_used      BIGINT  not null,
    effective_gas_price      BIGINT,
    created_contract_address binary(20),
    tx                       int,
    batch                    int     NOT NULL,
    INDEX (batch),
    INDEX (tx, batch),
    primary key (id)
);

create table if not exists tendb.contract
(
    id              INTEGER AUTO_INCREMENT,
    address         binary(20) NOT NULL,
    creator         int        NOT NULL,
    auto_visibility boolean    NOT NULL,
    transparent     boolean,
    tx              INTEGER,
    primary key (id),
    INDEX USING HASH (address)
);

create table if not exists tendb.externally_owned_account
(
    id      INTEGER AUTO_INCREMENT,
    address binary(20) NOT NULL,
    primary key (id),
    INDEX USING HASH (address)
);

create table if not exists tendb.event_type
(
    id              INTEGER AUTO_INCREMENT,
    contract        int        NOT NULL,
    event_sig       binary(32) NOT NULL,
    auto_visibility boolean    NOT NULL,
    auto_public     boolean,
    config_public   boolean    NOT NULL,
    topic1_can_view boolean,
    topic2_can_view boolean,
    topic3_can_view boolean,
    sender_can_view boolean,
    primary key (id),
    INDEX USING HASH (contract, event_sig)
);

create table if not exists tendb.event_topic
(
    id          INTEGER AUTO_INCREMENT,
    event_type  INTEGER,
    topic       binary(32) NOT NULL,
    rel_address INTEGER,
    primary key (id),
    INDEX USING HASH (topic),
    INDEX (rel_address)
);

create table if not exists tendb.event_log
(
    id         INTEGER AUTO_INCREMENT,
    event_type INTEGER NOT NULL,
    topic1     INTEGER,
    topic2     INTEGER,
    topic3     INTEGER,
    datablob   mediumblob,
    log_idx    INTEGER NOT NULL,
    receipt    INTEGER NOT NULL,
    primary key (id),
    INDEX (receipt, event_type, topic1, topic2, topic3),
    INDEX (event_type, topic1, topic2, topic3)
);