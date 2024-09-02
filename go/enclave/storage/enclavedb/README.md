This package contains the logic to access the relational database.

Currently, it is at the first version which is a little rough.

### Historical reasoning for the current design:

1. Geth is implemented on top of levelDB, and so makes heavy use of the k/v mechanism.
2. TEN uses EdglessDb in production, and sqlite in testing.
3. TEN started out by implementing the k/v primitives on top of sql.
4. Then we decided to use the SQL features to simplify some business logic, which required exposing the sql primitives.
5. Eventually, we ended up using SQL for everything, except the evm functionality.


The current abstractions are reflecting the first 4 steps.

Note: This package also mixes the `EnclaveDB` and `TransactionDB` implementations 

### SQL implementation

The implementation to read/write data to the db is also quite crude.

It mixes the SQL logic with the implementations for the EnclaveDB and DBTransaction. The reason it has to do that
is because they depend on the SQL k/v logic.

