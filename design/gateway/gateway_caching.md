# TEN Gateway Caching Design

## 1. Why cache requests on Gateway?

Currently, all `eth_` requests that hit the gateway are forwarded to the TEN Nodes and are executed by the TEN enclave. 
This is not ideal since there is only one Sequencer and it can be overloaded with requests
and there can be timeouts or errors in the sequencer.

To solve this problem, we can cache the responses of the `eth_` requests on the Gateway.
Not all requests can be cached, that is why I analyzed the Ethereum JSON-RPC Specification and formed two groups of requests that can be cached and those that cannot be cached.
For some methods, I used a rule of thumb since results can be cached for a certain period of time, but quickly become outdated.


Cacheable request methods are:

- `eth_accounts`
- `eth_chainID`
- `eth_coinbase`
- `eth_getBlockByHash`
- `eth_getBlockByNumber`
- `eth_getBlockRecipients`
- `eth_getBlockTransactionCountByHash`
- `eth_getBlockTransactionCountByNumber`
- `eth_getCode`
- `eth_getStorageAt`
- `eth_getTransactionByBlockHashAndIndex`
- `eth_getTransactionByBlockNumberAndIndex`
- `eth_getTransactionByHash`
- `eth_getTransactionReceipt`
- `eth_getUncleCountByBlockHash`
- `eth_getUncleCountByBlockNumber`
- `eth_maxPriorityFeePerGas`
- `eth_sign`
- `eth_signTransaction`

Methods cachable for a short period of time (TTL until next batch):

- `eth_getBalance`
- `eth_blockNumber`
- `eth_call`
- `eth_createAccessList`
- `eth_estimateGas`
- `eth_feeHistory`
- `eth_getProof`
- `eth_gasPrice`
- `eth_getFilterChanges`
- `eth_getFilterLogs`
- `eth_getTransactionCount`
- `eth_newBlockFilter`
- `eth_newFilter`
- `eth_newPendingTransactionFilter`
- `eth_sendRawTransaction`
- `eth_sendTransaction`
- `eth_syncing`
- `eth_uninstallFilter`

### Expiration time
Some cacheable methods will always produce the same result and can be cached indefinitely, while others will have a short expiration time.
Even for those that can be cached indefinitely I don't think it's a good idea to cache them for a long time,
since a large percentage of the requests will probably be requested only once
and caching will only consume memory and don't provide any benefit. 
As far as I can see it can help
to cache the results for shorter amount of time to help reduce the load on the sequencer / help in case of DDoS attacks
(but we need to do also other things to prevent them).



## 3. Implementation

For the implementation we can use some of the Go libraries for caching, such as: https://github.com/patrickmn/go-cache.
All the requests come from a single point in the code, so we can easily add caching to the gateway.



## Pros and cos of caching
I want to present also some negative effects of caching, so we can make a better decision.

### Pros
- it can reduce the load on the sequencer
- cah help if there is a spike of requests (e.g. DDoS attack / high load from specific user)

### Cons
- caching can consume additional memory
- caching can lead to outdated results (only in some methods)
- it is not a good way to prevent DDoS attacks, since users can easily request non-cacheable methods
