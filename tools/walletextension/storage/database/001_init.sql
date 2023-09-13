CREATE DATABASE ogdb;

USE ogdb;

GRANT SELECT, INSERT, UPDATE, DELETE ON ogdb.* TO 'obscurouser';

CREATE TABLE IF NOT EXISTS ogdb.users (
    user_id binary(32) PRIMARY KEY,
    private_key binary(32)
    );
CREATE TABLE IF NOT EXISTS ogdb.accounts (
    user_id binary(32),
    account_address binary(20),
    signature binary(65),
    FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
    );