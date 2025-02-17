# Platform-level traffic incentives 

We want a flexible mechanism to incentivise traffic.

For example:
- every transaction is rewarded with a ZEN token
- if you hold over 1000 TEN tokens, you receive two ZEN tokens
- if you use the newly launched app X during the first week, you double your reward
- randomly, tx senders could receive hidden prizes
- every 1000 transactions, you receive 100 ZEN
- etc

## Solution

In the Genesis block, we  deploy an upgradeable contract. The address is saved.
This contract will have a method `onTransactions(tx[])` that delegates to a dynamic implementation.

*Note1: The method receives all transactions from the batch.*
*Note2: The `tx` - is a reduced object containing everything except the calldata*


At the end of every batch, our EVM implementation will create a synthetic transaction that calls `onTransactions`.


## Discusson

- the logic should be simple because it will be applied to every transaction, and it will impact performance.
- The implementation can be changed at any time
- The contract should not fail.
- 
