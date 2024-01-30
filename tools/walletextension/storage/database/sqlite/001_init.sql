-- Enable foreign keys in SQLite
PRAGMA foreign_keys = ON;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    user_id binary(20) PRIMARY KEY,
    private_key binary(32)
    );

-- Create accounts table
CREATE TABLE IF NOT EXISTS accounts (
    user_id binary(20),
    account_address binary(20),
    signature binary(65),
    FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
    );