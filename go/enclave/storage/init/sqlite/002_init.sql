drop table rollup;

create table rollup
(
    hash              binary(32) primary key,
    start_seq         int  NOT NULL,
    end_seq           int  NOT NULL,
    header            blob NOT NULL,
    compression_block binary(32) NOT NULL
);