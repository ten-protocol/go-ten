# Obscuro JSON-RPC API

Obscuro supports a subset of Ethereum's [JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/). This 
page details which JSON-RPC API methods are supported.

## Supported JSON-RPC API methods

Obscuro nodes support the following JSON-RPC API methods over both HTTP and websockets:

* `eth_chainId`
* `eth_blockNumber`
* `eth_getBalance`
* `eth_getBlockByNumber`
* `eth_getBlockByHash`
* `eth_gasPrice`
* `eth_call`
* `eth_getTransactionReceipt`
* `eth_estimateGas`
* `eth_sendRawTransaction`
* `eth_getCode`
* `eth_getTransactionCount`
* `eth_getTransactionByHash`

## Supported subscription methods

When connecting via websockets, the following API methods are also exposed:

* `eth_subscribe`
* `eth_unsubscribe`

The only supported subscription type is `logs`.