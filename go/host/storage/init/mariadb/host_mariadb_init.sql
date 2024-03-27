CREATE TABLE IF NOT EXISTS block_host
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    hash        BINARY(32)      NOT NULL UNIQUE,
    header      BLOB            NOT NULL,
    rollup_hash BINARY(32)      NOT NULL
    );

CREATE INDEX IF NOT EXISTS IDX_BLOCK_HASH_HOST ON block_host (hash);

CREATE TABLE IF NOT EXISTS rollup_host
(
    id                INT AUTO_INCREMENT PRIMARY KEY,
    hash              BINARY(16) NOT NULL UNIQUE,
    start_seq         INT        NOT NULL,
    end_seq           INT        NOT NULL,
    time_stamp        INT        NOT NULL,
    ext_rollup        BLOB       NOT NULL,
    compression_block BINARY(32) NOT NULL
    );

CREATE INDEX IF NOT EXISTS IDX_ROLLUP_HASH_HOST ON rollup_host (hash);
CREATE INDEX IF NOT EXISTS IDX_ROLLUP_PROOF_HOST ON rollup_host (compression_block);
CREATE INDEX IF NOT EXISTS IDX_ROLLUP_SEQ_HOST ON rollup_host (start_seq, end_seq);

CREATE TABLE IF NOT EXISTS batch_host
(
    sequence       INT PRIMARY KEY,
    full_hash      BINARY(32) NOT NULL,
    hash           BINARY(16) NOT NULL UNIQUE,
    height         INT        NOT NULL,
    ext_batch      MEDIUMBLOB NOT NULL
    );

CREATE INDEX IF NOT EXISTS IDX_BATCH_HEIGHT_HOST ON batch_host (height);

CREATE TABLE IF NOT EXISTS transactions_host
(
    hash           BINARY(32) PRIMARY KEY,
    b_sequence     INT,
    FOREIGN KEY (b_sequence) REFERENCES batch_host(sequence)
    );

CREATE TABLE IF NOT EXISTS transaction_count
(
    id          INT  NOT NULL PRIMARY KEY,
    total       INT  NOT NULL
);

INSERT INTO transaction_count (id, total)
VALUES (1, 0) ON DUPLICATE KEY UPDATE id=id;
