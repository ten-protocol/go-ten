create table if not exists historical_contract_count
(
    id          int  NOT NULL PRIMARY KEY,
    total       int  NOT NULL
);

insert into historical_contract_count (id, total)
values (1, 0) on CONFLICT (id) DO NOTHING;