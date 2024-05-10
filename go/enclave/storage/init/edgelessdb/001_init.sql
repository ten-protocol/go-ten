CREATE USER "obscuro" REQUIRE ISSUER "/CN=obscuroCA" SUBJECT "/CN=obscuroUser";
CREATE DATABASE obsdb;

create table if not exists obsdb.keyvalue
(
    id  INTEGER AUTO_INCREMENT,
    ky  varbinary(64) NOT NULL,
    val mediumblob    NOT NULL,
    primary key (id),
    INDEX USING HASH (ky)
);
GRANT ALL ON obsdb.keyvalue TO obscuro;

create table if not exists obsdb.config
(
    ky  varchar(64),
    val mediumblob NOT NULL,
    primary key (ky)
);
GRANT ALL ON obsdb.config TO obscuro;

insert into obsdb.config
values ('CURRENT_SEQ', -1);

create table if not exists obsdb.attestation_key
(
    party binary(20),
    ky    binary(33) NOT NULL
);
GRANT ALL ON obsdb.attestation_key TO obscuro;

create table if not exists obsdb.block
(
    id           INTEGER AUTO_INCREMENT,
    hash         binary(32) NOT NULL,
    is_canonical boolean    NOT NULL,
    header       blob       NOT NULL,
    height       int        NOT NULL,
    primary key (id),
    INDEX (height),
    INDEX USING HASH (hash(8))
);
GRANT ALL ON obsdb.block TO obscuro;

create table if not exists obsdb.l1_msg
(
    id          INTEGER AUTO_INCREMENT,
    message     varbinary(1024) NOT NULL,
    block       INTEGER         NOT NULL,
    is_transfer boolean         NOT NULL,
    INDEX (block),
    primary key (id)
);
GRANT ALL ON obsdb.l1_msg TO obscuro;

create table if not exists obsdb.rollup
(
    id                INTEGER AUTO_INCREMENT,
    hash              binary(32) NOT NULL,
    start_seq         int        NOT NULL,
    end_seq           int        NOT NULL,
    time_stamp        int        NOT NULL,
    header            blob       NOT NULL,
    compression_block INTEGER    NOT NULL,
    INDEX (compression_block),
    INDEX USING HASH (hash(8)),
    primary key (id)
);
GRANT ALL ON obsdb.rollup TO obscuro;

create table if not exists obsdb.batch_body
(
    id      INTEGER,
    content mediumblob NOT NULL,
    primary key (id)
);
GRANT ALL ON obsdb.batch_body TO obscuro;

create table if not exists obsdb.batch
(
    sequence       INTEGER,
    converted_hash binary(32) NOT NULL,
    hash           binary(32) NOT NULL,
    height         int        NOT NULL,
    is_canonical   boolean    NOT NULL,
    header         blob       NOT NULL,
    body           int        NOT NULL,
    l1_proof_hash  binary(32),
    l1_proof       INTEGER,
    is_executed    boolean    NOT NULL,
    primary key (sequence),
    INDEX USING HASH (hash(8)),
    INDEX (body, l1_proof),
    INDEX (height)
);
GRANT ALL ON obsdb.batch TO obscuro;

create table if not exists obsdb.tx
(
    id             INTEGER AUTO_INCREMENT,
    hash           binary(32) NOT NULL,
    content        mediumblob NOT NULL,
    sender_address binary(20) NOT NULL,
    nonce          int        NOT NULL,
    idx            int        NOT NULL,
    body           int        NOT NULL,
    INDEX USING HASH (hash(8)),
    primary key (id)
);
GRANT ALL ON obsdb.tx TO obscuro;

create table if not exists obsdb.exec_tx
(
    id                       INTEGER AUTO_INCREMENT,
    created_contract_address binary(20),
    receipt                  mediumblob,
    tx                       int,
    batch                    int NOT NULL,
    INDEX (batch),
    INDEX (tx, created_contract_address(4)),
    primary key (id)
);
GRANT ALL ON obsdb.exec_tx TO obscuro;

create table if not exists obsdb.events
(
    id              INTEGER AUTO_INCREMENT,
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
    tx              int        NOT NULL,
    batch           int        NOT NULL,
    primary key (id),
    INDEX (tx, batch),
    INDEX USING HASH (address(8)),
    INDEX USING HASH (rel_address1(8)),
    INDEX USING HASH (rel_address2(8)),
    INDEX USING HASH (rel_address3(8)),
    INDEX USING HASH (rel_address4(8)),
    INDEX USING HASH (topic0(8)),
    INDEX USING HASH (topic1(8)),
    INDEX USING HASH (topic2(8)),
    INDEX USING HASH (topic3(8)),
    INDEX USING HASH (topic4(8))
);
GRANT ALL ON obsdb.events TO obscuro;