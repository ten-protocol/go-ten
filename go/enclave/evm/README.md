TEN executes EVM compatible transactions on top of a database that implements the Go-Ethereum interfaces.

The entry point is the evm_facade.

The approach we took was to depend on Go-Ethereum, mock out all consensus related dependencies, and just use the transaction execution functionality.