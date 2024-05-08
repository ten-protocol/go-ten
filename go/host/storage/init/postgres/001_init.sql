CREATE TABLE IF NOT EXISTS block_host
(
    id          SERIAL PRIMARY KEY,
    hash        BYTEA          NOT NULL UNIQUE,
    header      BYTEA          NOT NULL,
    rollup_hash BYTEA          NOT NULL
);

CREATE INDEX IF NOT EXISTS IDX_BLOCK_HASH_HOST ON block_host USING HASH (hash);

CREATE TABLE IF NOT EXISTS rollup_host
(
    id                SERIAL PRIMARY KEY,
    hash              BYTEA       NOT NULL UNIQUE,
    start_seq         INT         NOT NULL,
    end_seq           INT         NOT NULL,
    time_stamp        INT         NOT NULL,
    ext_rollup        BYTEA       NOT NULL,
    compression_block BYTEA       NOT NULL
);

CREATE INDEX IF NOT EXISTS IDX_ROLLUP_HASH_HOST ON rollup_host USING HASH (hash);
CREATE INDEX IF NOT EXISTS IDX_ROLLUP_PROOF_HOST ON rollup_host (compression_block);
CREATE INDEX IF NOT EXISTS IDX_ROLLUP_SEQ_HOST ON rollup_host (start_seq, end_seq);

CREATE TABLE IF NOT EXISTS batch_host
(
    sequence    INT PRIMARY KEY,
    full_hash   BYTEA         NOT NULL,
    hash        BYTEA         NOT NULL UNIQUE,
    height      INT           NOT NULL,
    ext_batch   BYTEA         NOT NULL
);

CREATE INDEX IF NOT EXISTS IDX_BATCH_HEIGHT_HOST ON batch_host (height);

CREATE TABLE IF NOT EXISTS transaction_host
(
    hash           BYTEA PRIMARY KEY,
    full_hash      BYTEA NOT NULL UNIQUE,
    b_sequence     INT,
    FOREIGN KEY (b_sequence) REFERENCES batch_host(sequence)
);

CREATE TABLE IF NOT EXISTS transaction_count
(
    id          SERIAL PRIMARY KEY,
    total       INT  NOT NULL
);

INSERT INTO transaction_count (id, total)
VALUES (1, 0)
    ON CONFLICT (id)
DO NOTHING;
