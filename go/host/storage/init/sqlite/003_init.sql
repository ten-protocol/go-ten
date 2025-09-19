create table if not exists historical_transaction_count
(
    id          int  NOT NULL PRIMARY KEY,
    total       int  NOT NULL
);

insert into historical_transaction_count (id, total)
values (1, 0) on CONFLICT (id) DO NOTHING;