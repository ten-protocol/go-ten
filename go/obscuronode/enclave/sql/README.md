This package contains a sql implementation of ethdb.Database.

Note: it seems quite odd to be creating a key value store with a sql database, but this allows us to plug in edgeless DB (which is a mysql-based database) as an enclave-secure persistent storage solution for Obscuro nodes