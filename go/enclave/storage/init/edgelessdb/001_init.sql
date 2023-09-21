CREATE USER "obscuro" REQUIRE ISSUER "/CN=obscuroCA" SUBJECT "/CN=obscuroUser";
CREATE DATABASE obsdb;

create table if not exists obsdb.keyvalue
(
    ky  varbinary(64),
    val mediumblob NOT NULL,
    primary key (ky)
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
    hash         binary(32),
    parent       binary(32),
    is_canonical boolean NOT NULL,
    header       blob    NOT NULL,
    height       int     NOT NULL,
    INDEX (parent),
    primary key (hash),
    INDEX (is_canonical),
    INDEX (height),
    INDEX (is_canonical, height)
);
GRANT ALL ON obsdb.block TO obscuro;

create table if not exists obsdb.l1_msg
(
    id      INTEGER AUTO_INCREMENT,
    message varbinary(1024) NOT NULL,
    block   binary(32)      NOT NULL,
    is_transfer boolean NOT NULL,
    INDEX (block),
    primary key (id)
);
GRANT ALL ON obsdb.l1_msg TO obscuro;

create table if not exists obsdb.rollup
(
    hash              binary(32),
    start_seq         int        NOT NULL,
    end_seq           int        NOT NULL,
    header            blob       NOT NULL,
    compression_block binary(32) NOT NULL,
    INDEX (compression_block),
    primary key (hash)
    );
GRANT ALL ON obsdb.rollup TO obscuro;

create table if not exists obsdb.batch_body
(
    hash    binary(32),
    content mediumblob NOT NULL,
    primary key (hash)
);
GRANT ALL ON obsdb.batch_body TO obscuro;

create table if not exists obsdb.batch
(
    hash         binary(32),
    parent       binary(32),
    sequence     int        NOT NULL,
    height       int        NOT NULL,
    is_canonical boolean    NOT NULL,
    header       blob       NOT NULL,
    body         binary(32) NOT NULL,
    l1_proof     binary(32) NOT NULL,
    is_executed  boolean    NOT NULL,
    INDEX (parent),
    INDEX (body),
    INDEX (l1_proof),
    INDEX (height),
    INDEX (sequence),
    INDEX (is_canonical),
    INDEX (is_executed),
    INDEX (is_canonical, is_executed),
    INDEX (is_canonical, is_executed, height),
    primary key (hash)
);
GRANT ALL ON obsdb.batch TO obscuro;

create table if not exists obsdb.tx
(
    hash           binary(32),
    content        mediumblob NOT NULL,
    sender_address binary(20) NOT NULL,
    nonce          int        NOT NULL,
    idx            int        NOT NULL,
    body           binary(32) NOT NULL,
    INDEX (body),
    primary key (hash)
);
GRANT ALL ON obsdb.tx TO obscuro;

create table if not exists obsdb.exec_tx
(
    id                       binary(64),
    created_contract_address binary(20),
    receipt                  mediumblob,
    tx                       binary(32) NOT NULL,
    batch                    binary(32) NOT NULL,
    INDEX (batch),
    INDEX (tx),
    primary key (id)
);
GRANT ALL ON obsdb.exec_tx TO obscuro;

create table if not exists obsdb.events
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
    exec_tx_id      binary(64) NOT NULL,
    INDEX (exec_tx_id),
    INDEX (address),
    INDEX (rel_address1),
    INDEX (rel_address2),
    INDEX (rel_address3),
    INDEX (rel_address4),
    INDEX (topic0),
    INDEX (topic1),
    INDEX (topic2),
    INDEX (topic3),
    INDEX (topic4)
);
GRANT ALL ON obsdb.events TO obscuro;