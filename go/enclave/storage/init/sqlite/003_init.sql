create table if not exists network_upgrades
(
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    feature_name         varchar(64)  NOT NULL,
    feature_data         mediumblob,
    applied_at_l1_height int          NOT NULL,
    applied_at_l1_hash   varchar(66)  NOT NULL,
    tx_hash              varchar(66)  NOT NULL,
    status               varchar(20)  NOT NULL DEFAULT 'pending', -- 'pending' or 'finalized'
    finalized_at_height  int,
    finalized_at_hash    varchar(66),
    created_at           timestamp default CURRENT_TIMESTAMP
);
create index IDX_UPGRADES_FEATURE on network_upgrades (feature_name);
create index IDX_UPGRADES_L1_HEIGHT on network_upgrades (applied_at_l1_height);
create index IDX_UPGRADES_HASH on network_upgrades (applied_at_l1_hash);
create index IDX_UPGRADES_STATUS on network_upgrades (status);
create index IDX_UPGRADES_TX_HASH on network_upgrades (tx_hash); 