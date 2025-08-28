CREATE TABLE IF NOT EXISTS block_count
(
    id          SERIAL PRIMARY KEY,
    total       INT  NOT NULL
);

INSERT INTO block_count (id, total)
VALUES (1, 0)
    ON CONFLICT (id)
DO NOTHING;