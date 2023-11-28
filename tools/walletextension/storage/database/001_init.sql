CREATE DATABASE ogdb;

USE ogdb;

GRANT SELECT, INSERT, UPDATE, DELETE ON ogdb.* TO 'obscurouser';

CREATE TABLE IF NOT EXISTS ogdb.users (
    user_id varbinary(20) PRIMARY KEY,
    private_key varbinary(32)
    );
CREATE TABLE IF NOT EXISTS ogdb.accounts (
    user_id varbinary(20),
    account_address varbinary(20),
    signature varbinary(65),
    FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
    );