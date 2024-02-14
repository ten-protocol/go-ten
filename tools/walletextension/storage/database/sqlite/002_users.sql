-- Create users table
CREATE TABLE IF NOT EXISTS users (
     user_id binary(20) PRIMARY KEY,
     private_key binary(32)
);