CREATE TABLE IF NOT EXISTS historical_contract_count
(
    id          SERIAL PRIMARY KEY,
    total       INT  NOT NULL
);

INSERT INTO historical_contract_count (id, total)
VALUES (1, 0)
    ON CONFLICT (id)
DO NOTHING;
