create table if not exists contract_host
(
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    address             binary(20) NOT NULL UNIQUE,
    creator             binary(20) NOT NULL,
    transparent      boolean    NOT NULL,
    custom_config   boolean    NOT NULL,
    batch_seq  int        NOT NULL,
    height     int        NOT NULL,
    time       int        NOT NULL
);

create index IDX_CONTRACT_ADDR_HOST on contract_host (address);
create index IDX_CONTRACT_CREATOR_HOST on contract_host (creator);
create index IDX_CONTRACT_DEPLOYED_HOST on contract_host (time DESC);
create index IDX_CONTRACT_TRANSPARENT_HOST on contract_host (transparent, custom_config);

create table if not exists contract_count
(
    id          int  NOT NULL PRIMARY KEY,
    total       int  NOT NULL
);

insert into contract_count (id, total)
values (1, 0) on CONFLICT (id) DO NOTHING;

