# Handling Sensitive Information

Obscuro supports a subset of Ethereum's [JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/).

Some of these methods deal with sensitive information. For example, the response to an `eth_getBalance` request will 
contain the balance of an account. An attacker could intercept this response to discover a user's balance. To avoid 
this, the requests and responses for methods deemed sensitive are encrypted and decrypted by the 
[wallet extension](wallet-extension.md). To provide a good user experience, this process is invisible to the end user.

This page details which JSON-RPC API methods are supported, which ones are deemed sensitive, and the rules governing 
who is able to decrypt the response to a given method call.