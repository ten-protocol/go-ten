CREATE TABLE IF NOT EXISTS block_host
(
    id          SERIAL PRIMARY KEY,
    hash        BYTEA          NOT NULL,
    header      BYTEA          NOT NULL
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
    compression_block INT       NOT NULL,
    FOREIGN KEY (compression_block) REFERENCES block_host(id)
    );

CREATE INDEX IF NOT EXISTS IDX_ROLLUP_HASH_HOST ON rollup_host USING HASH (hash);
CREATE INDEX IF NOT EXISTS IDX_ROLLUP_PROOF_HOST ON rollup_host (compression_block);
CREATE INDEX IF NOT EXISTS IDX_ROLLUP_SEQ_HOST ON rollup_host (start_seq, end_seq);

CREATE TABLE IF NOT EXISTS cross_chain_message_host
(
    id                SERIAL PRIMARY KEY,
    message_hash      BYTEA       NOT NULL UNIQUE,
    rollup_id         INT         NOT NULL,
    message_type      CHAR(1)     NOT NULL CHECK (message_type IN ('m', 'v')),
    FOREIGN KEY (rollup_id) REFERENCES rollup_host(id)
);

CREATE INDEX IF NOT EXISTS IDX_CCM_HASH_HOST ON cross_chain_message_host USING HASH (message_hash);
CREATE INDEX IF NOT EXISTS IDX_CCM_ROLLUP_HOST ON cross_chain_message_host (rollup_id);

CREATE TABLE IF NOT EXISTS batch_host
(
    sequence    INT PRIMARY KEY,
    hash        BYTEA         NOT NULL ,
    height      INT           NOT NULL,
    ext_batch   BYTEA         NOT NULL,
    txs_size    INT           NOT NULL
);

CREATE INDEX IF NOT EXISTS IDX_BATCH_HASH_HOST ON batch_host USING HASH (hash);
CREATE INDEX IF NOT EXISTS IDX_BATCH_HEIGHT_HOST ON batch_host (height);

CREATE TABLE IF NOT EXISTS transaction_host
(
    id             SERIAL PRIMARY KEY,
    hash           BYTEA,
    b_sequence     INT,
    FOREIGN KEY (b_sequence) REFERENCES batch_host(sequence)
);

CREATE INDEX IF NOT EXISTS IDX_TX_HASH_HOST ON transaction_host USING HASH (hash);
CREATE INDEX IF NOT EXISTS IDX_TX_SEQ_HOST ON transaction_host (b_sequence);

CREATE TABLE IF NOT EXISTS transaction_count
(
    id          SERIAL PRIMARY KEY,
    total       INT  NOT NULL
);

INSERT INTO transaction_count (id, total)
VALUES (1, 0)
    ON CONFLICT (id)
DO NOTHING;

CREATE TABLE IF NOT EXISTS block_count
(
    id          SERIAL PRIMARY KEY,
    total       INT  NOT NULL
);

INSERT INTO block_count (id, total)
VALUES (1, 0)
    ON CONFLICT (id)
DO NOTHING;
