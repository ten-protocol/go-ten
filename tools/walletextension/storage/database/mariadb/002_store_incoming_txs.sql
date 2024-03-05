/*
    This is a migration file for MariaDB that creates transactions table for storing incoming transactions
*/

CREATE TABLE IF NOT EXISTS ogdb.transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id varbinary(20),
    tx_hash TEXT,
    tx TEXT,
    tx_time DATETIME DEFAULT CURRENT_TIMESTAMP
);