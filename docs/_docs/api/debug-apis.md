---
---
# Debug JSON-RPC API

Obscuro supports a subset of Geth's [DEBUG JSON-RPC API](https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-debug). This 
page details which Debug JSON-RPC API methods are supported.

## Supported JSON-RPC API methods

Obscuro nodes support the following JSON-RPC API methods over both HTTP and websockets:

* `debug_traceTransaction`
* `debug_eventLogRelevancy`: returns a list of log entries, each containing the relevancy and lifecycleEvent fields.

## debug_LogVisibility

Request Payload:
```go
{
    "jsonrpc": "2.0",
    "method": "debug_eventLogRelevancy",
    "params": [
        "0xb29737963fd6768587ede453ab90ff7668115db16915a7833705ef134e793814"
    ],
    "id": 1
}
```

Request result:
```go
{
    "jsonrpc": "2.0",
    "id": 1,
    "result": [
        {
            "relAddress1": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "relAddress2": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "relAddress3": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "relAddress4": "0x0000000000000000000000000000000000000000000000000000000000000000",
            "lifecycleEvent": true,
            "address": "0x9802f661d17c65527d7abb59daad5439cb125a67",
            "topics": [
                "0xebfcf7c0a1b09f6499e519a8d8bb85ce33cd539ec6cbd964e116cd74943ead1a"
            ],
            "data": "0x000000000000000000000000987e0a0692475bcc5f13d97e700bb43c1913effe0000000000000000000000000000000000000000000000000000000000000001",
            "blockNumber": "0x98",
            "transactionHash": "0xb29737963fd6768587ede453ab90ff7668115db16915a7833705ef134e793814",
            "transactionIndex": "0x0",
            "blockHash": "0x92c15c0a6784c4f7e3133d9bd31d3202b90de8aa66636b2f597fea9c1e76fc1b",
            "logIndex": "0x0",
            "removed": false
        }
    ]
}
```
