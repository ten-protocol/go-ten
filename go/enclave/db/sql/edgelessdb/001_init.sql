CREATE USER "obscuro" REQUIRE ISSUER "/CN=obscuroCA" SUBJECT "/CN=obscuroUser";
CREATE DATABASE obsdb;

CREATE TABLE obsdb.keyvalue
(
    ky  varbinary(64) primary key,
    val mediumblob
);

create table obsdb.events
(
    topic0         binary(32),
    topic1         binary(32),
    topic2         binary(32),
    topic3         binary(32),
    topic4         binary(32),
    datablob       mediumblob,
    blockHash      binary(32),
    blockNumber    int,
    txHash         binary(32),
    txIdx          int,
    logIdx         int,
    address        binary(20),
    lifecycleEvent boolean,
    relAddress1    binary(20),
    relAddress2    binary(20),
    relAddress3    binary(20),
    relAddress4    binary(20)
);

create index IX_RAD1 on obsdb.events (relAddress1);
create index IX_RAD2 on obsdb.events (relAddress2);
create index IX_RAD3 on obsdb.events (relAddress3);
create index IX_RAD4 on obsdb.events (relAddress4);
create index IX_AD on obsdb.events (address);
create index IX_BLH on obsdb.events (blockHash);
create index IX_BLN on obsdb.events (blockNumber);
create index IX_TXH on obsdb.events (txHash);
create index IX_T0 on obsdb.events (topic0);
create index IX_T1 on obsdb.events (topic1);
create index IX_T2 on obsdb.events (topic2);
create index IX_T3 on obsdb.events (topic3);
create index IX_T4 on obsdb.events (topic4);
GRANT ALL ON obsdb.keyvalue TO obscuro;
GRANT ALL ON obsdb.events TO obscuro;