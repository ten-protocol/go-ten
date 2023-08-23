Implements a very primitive database migration.

Each supported db folder must contain a list of sql files of the format "$order_$desc.sql".
Where "$order" is a number, and these numbers have to be consecutive.

Upon starting, there will be logic fetching the number stored in the database, and then execute 
all sql files which are more recent, in a single transaction. 