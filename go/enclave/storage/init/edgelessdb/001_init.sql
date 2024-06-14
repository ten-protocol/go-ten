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
    INDEX USING HASH (hash)
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
    INDEX USING HASH (hash),
    primary key (id)
);
GRANT ALL ON obsdb.rollup TO obscuro;

create table if not exists obsdb.batch
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
    INDEX (height)
);
GRANT ALL ON obsdb.batch TO obscuro;

create table if not exists obsdb.tx
(
    id             INTEGER AUTO_INCREMENT,
    hash           binary(32) NOT NULL,
    content        mediumblob NOT NULL,
    sender_address int NOT NULL,
    idx            int        NOT NULL,
    batch_height   int        NOT NULL,
    INDEX USING HASH (hash),
    INDEX (sender_address),
    INDEX (batch_height, idx),
    primary key (id)
);
GRANT ALL ON obsdb.tx TO obscuro;

create table if not exists obsdb.exec_tx
(
    id                       INTEGER AUTO_INCREMENT,
    created_contract_address int,
    receipt                  mediumblob,
    tx                       int,
    batch                    int NOT NULL,
    INDEX (batch),
    INDEX (created_contract_address, tx),
    primary key (id)
);
GRANT ALL ON obsdb.exec_tx TO obscuro;

create table if not exists obsdb.contract
(
    id      INTEGER AUTO_INCREMENT,
    address binary(20) NOT NULL,
    owner   int        NOT NULL,
    primary key (id),
    INDEX USING HASH (address)
);
GRANT ALL ON obsdb.contract TO obscuro;

create table if not exists obsdb.externally_owned_account
(
    id      INTEGER AUTO_INCREMENT,
    address binary(20) NOT NULL,
    primary key (id),
    INDEX USING HASH (address)
);
GRANT ALL ON obsdb.externally_owned_account TO obscuro;

create table if not exists obsdb.event_type
(
    id              INTEGER AUTO_INCREMENT,
    contract        int        NOT NULL,
    event_sig       binary(32) NOT NULL,
    lifecycle_event boolean    NOT NULL,
    primary key (id),
    INDEX USING HASH (contract, event_sig)
);
GRANT ALL ON obsdb.event_type TO obscuro;

create table if not exists obsdb.event_topic
(
    id          INTEGER AUTO_INCREMENT,
    topic       binary(32) NOT NULL,
    rel_address INTEGER,
    primary key (id),
    INDEX USING HASH (topic),
    INDEX (rel_address)
);
GRANT ALL ON obsdb.event_topic TO obscuro;

create table if not exists obsdb.event_log
(
    id         INTEGER AUTO_INCREMENT,
    event_type INTEGER NOT NULL,
    topic1     INTEGER,
    topic2     INTEGER,
    topic3     INTEGER,
    datablob   mediumblob,
    log_idx    INTEGER NOT NULL,
    exec_tx    INTEGER NOT NULL,
    primary key (id),
    INDEX (exec_tx, event_type, topic1, topic2, topic3)
);
GRANT ALL ON obsdb.event_log TO obscuro;