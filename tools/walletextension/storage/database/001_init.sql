/*
 This file is used to create the database and set the necessary permissions for the user that will be used by the gateway.
 */

-- Create the database
CREATE DATABASE IF NOT EXISTS ogdb;

-- Grant the necessary permissions
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, DROP ON ogdb.* TO 'obscurouser';

-- Reload the privileges from the grant tables in the mysql database
FLUSH PRIVILEGES;