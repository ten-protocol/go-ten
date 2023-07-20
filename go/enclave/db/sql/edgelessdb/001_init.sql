CREATE USER "obscuro" REQUIRE ISSUER "/CN=obscuroCA" SUBJECT "/CN=obscuroUser";
CREATE DATABASE obsdb;

create table if not exists obsdb.keyvalue
(
    ky  varbinary(64),
    val mediumblob,
    primary key (ky)
);
GRANT ALL ON obsdb.keyvalue TO obscuro;

create table if not exists obsdb.config
(
    ky  varchar(64),
    val mediumblob,
    primary key (ky)
);
GRANT ALL ON obsdb.config TO obscuro;

insert into obsdb.config
values ('CURRENT_SEQ', -1);

create table if not exists obsdb.attestation_key
(
    party binary(20),
    ky    binary(33)
    -- primary key (party)
);
GRANT ALL ON obsdb.attestation_key TO obscuro;

create table if not exists obsdb.block
(
    hash         binary(32),
    parent       binary(32),
    is_canonical boolean,
    header       text,
    height       int,
    FOREIGN KEY (parent) REFERENCES obsdb.block (hash),
    primary key (hash),
    INDEX (height)
);
GRANT ALL ON obsdb.block TO obscuro;

create table if not exists obsdb.l1_msg
(
    id      INTEGER AUTO_INCREMENT,
    message varbinary(1024),
    block   binary(32),
    FOREIGN KEY (block) REFERENCES obsdb.block (hash),
    primary key (id)
);
GRANT ALL ON obsdb.l1_msg TO obscuro;

create table if not exists obsdb.rollup
(
    id        INTEGER AUTO_INCREMENT,
    start_seq int,
    end_seq   int,
    header    text,
    block     binary(32),
    FOREIGN KEY (block) REFERENCES obsdb.block (hash),
    primary key (id)
);
GRANT ALL ON obsdb.rollup TO obscuro;

create table if not exists obsdb.batch_body
(
    hash    binary(32),
    content mediumblob,
    primary key (hash)
);
GRANT ALL ON obsdb.batch_body TO obscuro;

create table if not exists obsdb.batch
(
    hash         binary(32),
    parent       binary(32),
    sequence     int,
    height       int,
    is_canonical boolean,
    header       text,
    body         binary(32),
    l1_proof     binary(32),
    source       text,
    FOREIGN KEY (parent) REFERENCES obsdb.batch (hash),
    FOREIGN KEY (body) REFERENCES obsdb.batch_body (hash),
    FOREIGN KEY (l1_proof) REFERENCES obsdb.block (hash),
    INDEX (height),
    INDEX (sequence),
    primary key (hash)
);
GRANT ALL ON obsdb.batch TO obscuro;

create table if not exists obsdb.tx
(
    hash           binary(32),
    content        mediumblob,
    sender_address binary(20),
    nonce          int,
    idx            int,
    body           binary(32),
    FOREIGN KEY (body) REFERENCES obsdb.batch_body (hash),
    primary key (hash)
);
GRANT ALL ON obsdb.tx TO obscuro;

create table if not exists obsdb.exec_tx
(
    id                       binary(64),
    created_contract_address binary(32),
    receipt                  mediumblob,
    tx                       binary(32),
    batch                    binary(32),
    FOREIGN KEY (batch) REFERENCES obsdb.batch (hash),
--     commenting out the fk until synthetic transactions are also stored
--    FOREIGN KEY (tx) REFERENCES obsdb.tx (hash),
    INDEX (tx),
    primary key (id)
);
GRANT ALL ON obsdb.exec_tx TO obscuro;

create table if not exists obsdb.events
(
    topic0          binary(32),
    topic1          binary(32),
    topic2          binary(32),
    topic3          binary(32),
    topic4          binary(32),
    datablob        mediumblob,
    log_idx         int,
    address         binary(32),
    lifecycle_event boolean,
    rel_address1    binary(20),
    rel_address2    binary(20),
    rel_address3    binary(20),
    rel_address4    binary(20),
    exec_tx_id      binary(64),
    FOREIGN KEY (exec_tx_id) REFERENCES obsdb.exec_tx (id),
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