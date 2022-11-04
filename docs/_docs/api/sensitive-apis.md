---
---
# Sensitive APIs

Obscuro supports a subset of Ethereum's [JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/).

Some of these methods deal with sensitive information. For example, the response to an `eth_getBalance` request will 
contain the balance of an account. An attacker could intercept this response to discover a user's balance. To avoid 
this, the requests and responses for methods deemed sensitive are encrypted and decrypted by the 
[wallet extension](../wallet-extension/wallet-extension.md). To ensure a good user experience, this process is 
invisible to the end user.

This page details which JSON-RPC API methods are deemed sensitive, and the rules governing who is able to decrypt the 
response to a given method call.

## Sensitive JSON-RPC API Methods

Of the methods above, the following are deemed sensitive, and their requests and responses are encrypted in transit:
 
* `eth_call`: Response can only be decrypted by the owner of the account in the request's `from` field
* `eth_estimateGas`: Response can only be decrypted by the owner of the account in the request's `from` field
* `eth_getBalance`: Response can only be decrypted by:
    * For account addresses: The owner of the account
    * For contract addresses: The owner of the account that deployed the contract
* `eth_getLogs`: Response can only be decrypted by the owner of the account, and only includes logs relevant to that 
   account
* `eth_getTransactionByHash`: Response can only be decrypted by the signer of the transaction
* `eth_getTransactionCount`: Response can only be decrypted by the owner of the address
* `eth_getTransactionReceipt`: Response can only be decrypted by the signer of the transaction
* `eth_sendRawTransaction`: Response can only be decrypted by the signer of the transaction
