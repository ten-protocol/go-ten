-- Create table for storing network upgrades
create table if not exists tendb.network_upgrade
(
    id                  int            NOT NULL AUTO_INCREMENT,
    feature_name        varchar(255)   NOT NULL,
    feature_data        mediumblob     NOT NULL,
    block_hash          binary(32)     NOT NULL,
    block_height_final  bigint,
    block_height_active bigint,
    created_at          timestamp      DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- Indexes for efficient queries
CREATE INDEX feature ON tendb.network_upgrade (feature_name);
CREATE INDEX block_hash ON tendb.network_upgrade (block_hash);
CREATE INDEX heights ON tendb.network_upgrade (block_height_final, block_height_active);
