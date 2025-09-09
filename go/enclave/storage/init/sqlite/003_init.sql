-- Create table for storing network upgrades
create table if not exists network_upgrade
(
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    feature_name        varchar(255)   NOT NULL,
    feature_data        mediumblob     NOT NULL,
    block_hash          binary(32)     NOT NULL,
    block_height_final  bigint,
    block_height_active bigint,
    created_at          timestamp      DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for efficient queries
create index IDX_NETWORK_UPGRADE_FEATURE on network_upgrade (feature_name);
create index IDX_NETWORK_UPGRADE_BLOCK_HASH on network_upgrade (block_hash);
create index IDX_NETWORK_UPGRADE_HEIGHTS on network_upgrade (block_height_final, block_height_active);
