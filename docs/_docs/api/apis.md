---
---
# Obscuro JSON-RPC API

Obscuro supports a subset of Ethereum's [JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/). This 
page details which JSON-RPC API methods are supported.

## Supported JSON-RPC API methods

Obscuro nodes support the following JSON-RPC API methods over both HTTP and websockets:

* `eth_blockNumber`
* `eth_call`
* `eth_chainId`
* `eth_estimateGas`
* `eth_gasPrice`
* `eth_getBalance`
* `eth_getBlockByHash`
* `eth_getBlockByNumber`
* `eth_getCode`
* `eth_getLogs`
* `eth_getTransactionByHash`
* `eth_getTransactionCount`
* `eth_getTransactionReceipt`
* `eth_sendRawTransaction`

## Supported subscription methods

When connecting via websockets, the following API methods are also exposed:

* `eth_subscribe`
* `eth_unsubscribe`

The only supported subscription type is `logs`.