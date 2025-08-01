create table if not exists tendb.network_upgrades
(
    id                   INTEGER AUTO_INCREMENT,
    feature_name         varchar(64)  NOT NULL,
    feature_data         mediumblob,
    applied_at_l1_height int          NOT NULL,
    applied_at_l1_hash   varchar(66)  NOT NULL,
    tx_hash              varchar(66)  NOT NULL,
    status               varchar(20)  NOT NULL DEFAULT 'pending', -- 'pending' or 'finalized'
    finalized_at_height  int,
    finalized_at_hash    varchar(66),
    created_at           timestamp default CURRENT_TIMESTAMP,
    primary key (id),
    INDEX (feature_name),
    INDEX (applied_at_l1_height),
    INDEX (applied_at_l1_hash),
    INDEX (status),
    INDEX (tx_hash)
); 