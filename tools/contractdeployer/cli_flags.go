package contractdeployer

// Flag names, defaults and usages.
const (
	nodeHostName  = "nodeHost"
	nodeHostUsage = "The host on which to connect the RPC client"

	nodePortName  = "nodePort"
	nodePortUsage = "The port on which to connect the RPC client"

	contractNameName  = "contractName"
	contractNameUsage = "The name of the contract bytecode to be deploy (e.g. `MGMT` or `ERC20`)"

	privateKeyName  = "privateKey"
	privateKeyUsage = "The private key for the node account"

	chainIDName  = "chainID"
	chainIDUsage = "The ID of the chain"
)
