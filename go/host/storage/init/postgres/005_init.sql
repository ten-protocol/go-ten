CREATE TABLE IF NOT EXISTS contract_host
(
    id                  SERIAL PRIMARY KEY,
    address             BYTEA      NOT NULL UNIQUE,
    creator             BYTEA      NOT NULL,
    transparent      BOOLEAN    NOT NULL,
    custom_config   BOOLEAN    NOT NULL,
    batch_seq  BIGINT     NOT NULL,
    height     BIGINT     NOT NULL,
    time       INT     NOT NULL
);

CREATE INDEX IF NOT EXISTS IDX_CONTRACT_ADDR_HOST ON contract_host USING HASH (address);
CREATE INDEX IF NOT EXISTS IDX_CONTRACT_CREATOR_HOST ON contract_host (creator);
CREATE INDEX IF NOT EXISTS IDX_CONTRACT_DEPLOYED_HOST ON contract_host (time DESC);
CREATE INDEX IF NOT EXISTS IDX_CONTRACT_TRANSPARENT_HOST ON contract_host (transparent, custom_config);

CREATE TABLE IF NOT EXISTS contract_count
(
    id          SERIAL PRIMARY KEY,
    total       INT  NOT NULL
);

INSERT INTO contract_count (id, total)
VALUES (1, 0)
    ON CONFLICT (id)
DO NOTHING;

