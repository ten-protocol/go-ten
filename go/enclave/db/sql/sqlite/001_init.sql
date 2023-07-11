create table if not exists keyvalue
(
    ky  varbinary(64) primary key,
    val mediumblob
);

create table if not exists events
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
    address        binary(32),
    lifecycleEvent boolean,
    relAddress1    binary(20),
    relAddress2    binary(20),
    relAddress3    binary(20),
    relAddress4    binary(20)
);

create index IX_AD on events (address);
create index IX_RAD1 on events (relAddress1);
create index IX_RAD2 on events (relAddress2);
create index IX_RAD3 on events (relAddress3);
create index IX_RAD4 on events (relAddress4);
create index IX_BLH on events (blockHash);
create index IX_BLN on events (blockNumber);
create index IX_TXH on events (txHash);
create index IX_T0 on events (topic0);
create index IX_T1 on events (topic1);
create index IX_T2 on events (topic2);
create index IX_T3 on events (topic3);
create index IX_T4 on events (topic4);
