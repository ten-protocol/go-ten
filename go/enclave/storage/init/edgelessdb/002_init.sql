drop table obsdb.rollup;

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
