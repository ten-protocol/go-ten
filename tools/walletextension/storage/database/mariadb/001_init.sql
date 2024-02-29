/*
    This is a migration file for MariaDB and is executed when the Gateway is started to make sure the database schema is up to date.
*/

CREATE TABLE IF NOT EXISTS ogdb.users (
    user_id varbinary(20) PRIMARY KEY,
    private_key varbinary(32)
);

CREATE TABLE IF NOT EXISTS ogdb.accounts (
    user_id varbinary(20),
    account_address varbinary(20),
    signature varbinary(65),
    FOREIGN KEY(user_id) REFERENCES ogdb.users(user_id) ON DELETE CASCADE
);
