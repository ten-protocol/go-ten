# Events in TEN - Design document

## Scope

This is a design proposal for how to handle events in TEN.

It covers two aspects:

- the visibility rules for events.
- technical implementation details.

## Requirements

* Each subscription is tied to a specific account, and will only return events relevant to that account (the concept of
  relevancy is defined below)
* Subscription requests and responses are encrypted in transit, and can only be decrypted by the enclave on one end, 
  and the viewing key for the account the subscription is tied to on the other
* Subscription requests are authenticated; a node cannot create a subscription on behalf of another account (e.g. to 
  snoop on the pattern of traffic for that account)
* This encryption reuses the wallet extension and its existing viewing keys, and is transparent to the user
* The logs in transaction receipts are also filtered to only include events relevant to the account of the transaction 
  submitter

## Event visibility design

In Ten, we aim to maintain the same events APIs that are found in Ethereum, and will try to implement privacy 
concerns with as few changes as possible. For background on how Ethereum handles events, see the section 
"Background - Ethereum Events Design", below.

### Event types

To decide whether a user Alice is entitled to view an event, we must consider three different types of events:

1. The event was emitted as a result of a transaction sent by Alice.
2. The event was emitted as a result of a transaction sent by Bob, and is relevant to Alice (we define relevancy 
   below).
3. The event was emitted as a result of a transaction sent by Bob, but is not relevant to Alice

### Event visibility rules

Each event contains an array of 32-byte hex strings, called topics. In Ten, as in Ethereum, end-users are 
identified by one or more account addresses. Since we are reusing Ethereum's API without modifications, our goal is to 
use an event's topics to ascertain who the event is relevant to, with users able to see an event if and only if the 
event is considered relevant to them.

We propose the following rules:

1. An event is considered relevant to all account owners whose addresses are used as topics in the event.
2. In case there are no account addresses in an event's topics, then the event is considered relevant to everyone 
   (known as a "lifecycle event").

The purpose for these rules is to be simple, clear, intuitive, and to work as well as possible with the existing 
contracts. See the section "Events and address fields", below, for an analysis of several events from ERC20 and 
Uniswap, grouped by whether they contain address fields. However, there are several edge-cases to be handled:

* The event contains addresses, but they are contract addresses
* The event contains a topic that is not an address, but looks like one (e.g. `0x00000000000000000000000000000001`)

To handle this, for each potential account address in the topics, we check:

* Whether the address has associated contract code, using the `GetCode` function
* Whether the address has a zero nonce, using the `GetNonce` function (only account and smart contract addresses can  
  have non-zero nonces)

If one or both of the above are true for every address in the event's topics, the event will fall under rule (2), above.

### Adjusting the event visibility rules

In case the rules above are not providing the desired functionality, the developer will have a couple of options to adjust visibility. 

For example, if one of the lifecycle events should only be visible to the administrators, the developer can add that address as a topic.

In case an event contains a account address topic that should not contribute to relevancy, the developer can remove the 
`"indexed"` and thus the event will become invisible to that user. This might not be ideal in case this address has to 
be used for subscribing. 

As the ultimate flexible mechanism we propose a programmatic way to determine whether a requester is allowed to view an event. 

If the implicit rules are not satisfactory, the smart contract developer can define a view function called 
``eventVisibility``, which will be called by the TEN VM behind the scenes.

```solidity
   // If declared in a smart contract, this function will be called by the TEN VM to determine whether the requester
   // address is allowed to view the event. 
   function eventVisibility(address requester, bytes32[] memory topics,  bytes32[] memory data) external view returns (bool){
       // based on the data from the event, passed in as an array of topic and data (which is the internal VM representation) 
       // calculate whether the requester address is entitled to view the event
       // - the first element in the topics array is the hash of the event signature
       // Note: the developer must maintain a list of the event hashes they want to handle programmatically 
        return ; 
    }
```    

To determine the visibility of an event, the TEN VM will do the following:

1. call the `eventVisibility` with the event being requested and the requester.
2. If the function exists and returns 'true', then return the event. If it returns `false`, then the event is invisible.
3. If the function does not exist, apply the implicit rules.

### Alternative event visibility rules considered

#### Make events public by default, and remove privacy-leaking events

Removing privacy-leaking events is not feasible. For example, the correct functioning of the ERC20 contract relies on a 
public, privacy-leaking `Transfer` event.

#### Make the events relevant to the signer of the transaction that generated the event as well

This rule adds another dimension to the reasoning process, because there is an implicit user to whom the event is relevant.
It also reduces flexibility in sending lifecycle events to administrators.

## TEN events implementation

Our goal is to implement the visibility rules described above without modifying the Ethereum events API. We will look 
first at the changes to the enclave, then those to the host, and finally those to the RPC client.

### TEN enclave

#### Creating and deleting logs subscriptions

The enclave exposes RPC methods to add or remove logs subscriptions. The request to create a new logs subscription must 
contain the following information:

* `Account`: The account address the events relate to.
* `Signature`: A signature over the account address using a private viewing key.
* `Filter`: A subscriber-defined filter to apply to the stream of logs.

This object is encrypted with the enclave's private key, to prevent attackers from eavesdropping on the creation of 
subscriptions.

The signature field proves that the client is authorised to create a subscription for the given account. This prevents
attackers from creating subscriptions for accounts other than their own to analyse the pattern of logs.

Each new logs subscription is assigned a unique subscription ID that is returned to the host. The method to delete a 
subscription takes the subscription ID of the subscription to be deleted. This method is not authenticated, so an 
attacker who discovered the subscription ID could request its deletion.

#### Block ingestion

Each time the host sends the enclave a new block to ingest, the enclave stores the logs for that block in its database. 
It then produces a block submission response which it sends back to the host.

This response includes a mapping from subscription IDs to a set of logs, where each set of logs is encrypted with 
the viewing key corresponding to the subscription's `Account`. The set of logs is constructed by taking any logs from 
the current block that meet the two following criteria:

* They match the `LogSubscription.Filter`
* They pass the relevancy test described earlier in this document for the account address in `LogSubscription.Account`

#### Transaction receipts

Transaction receipts also contain the logs for that transaction. Whenever the host requests a transaction receipt, the 
enclave filters out the receipt's logs to only include those that pass the relevancy test, using the transaction's 
sender as the user address.

#### Log snapshots

The enclave also exposes an RPC method to get a snapshot of logs. It expects an encrypted set of params that include 
the log filter and the address the logs are for. It then crawls the chain, extracts all the logs that match the filter 
and are relevant based on the address provided. It returns these logs encrypted with the viewing key corresponding to 
the address provided.

### TEN host

#### Logs subscriptions

For each incoming logs subscription request via RPC, the host has to do two things:

1. Route the new subscription request to the enclave and create a new Geth `rpc.Subscription` to return to the client
2. Extract the logs upon each block ingestion and route them to the corresponding client subscription

Upon creating the new subscription, the host immediately sends the subscription ID as the first message. This is 
required by the client, so that it can be returned to the user and used to unsubscribe if needed at a later time.

Since the log subscription request is encrypted by the client and can only be decrypted on the enclave, the host 
forwards it blindly, and cannot learn anything about the contents of the subscription. However, it does generate a 
fresh ID for the subscription, which it also forwards to the enclave. The enclave reuses this ID when sending back 
logs for that subscription in the block submission response.

For each block submission response, the host extracts the encrypted logs for that subscription ID, and forwards them 
on to the corresponding Geth `rpc.Subscription`.

#### Logs snapshots

For log snapshot requests, it is again forwarded blindly to the enclave, with the host unable to learn anything about 
the request or response.

### TEN encrypted RPC client

Due to their sensitive nature, logs requests and log subscriptions must pass through the encrypted RPC client.

#### Logs subscriptions

The encrypted RPC client only handles logs subscriptions via the `eth_subscribe` and `eth_unsubscribe` APIs (see 
[here](https://ethereum.org/en/developers/tutorials/using-websockets/#eth-subscribe)). A consequence of this is that 
events are only available in TEN over a websocket connection.

In response to the incoming `eth_subscribe` request, the client creates a `logs` subscription to the host by making a
`rpc.Client.Subscribe` call via the embedded Geth client. It passes as a parameter a `LogSubscription` object encrypted 
with the enclave's public key to protect the request from eavesdroppers, setting the `Account` to the client's 
account and generating the required signature.

For each received event, the encrypted RPC client decrypts the encrypted log bytes with the corresponding private key 
before returning the log events to the user. It then pushes the decrypted event onto a separate channel listened to by 
the client.

## Security and usability of the proposed design

App developers will be able to use the existing libraries unchanged, as long as they connect through a wallet extension with registered viewing keys.

Depending on the subscription, the results might be different from those returned on a normal Ethereum network, because the user might not have the right to see certain private data. 

Smart contract developers need to spend a few minutes to think about whether an event can be seen by the entity who it is relevant to, or whether it leaks data.
There is no new syntax to learn, just to be aware about the visibility rules. In case the default intuitive rules do not satisfy the requirements,
the developer can write a function to precisely control visibility.

The data access protections of smart contracts will prevent another smart contract interacting with it from extracting information and leaking it as an event.  

The fact that the wallet extension adds signed accounts to each subscription request, makes it impossible for a user to request the events of another user.

An ERC20 transfer from Alice to Bob will show up on Bob's UI if he is subscribed to it, but will not show on Charlie's UI.

## Open implementation questions

1. What's the lifespan for a subscription in the enclave, unsubscribe option + an expiry?

2. Still pretty worried about perf and DoS potential, but we can test for it and try to optimise/prioritise traffic etc. Maybe we need read-only nodes to services these requests in production (normal nodes but the host doesn't give it any transactions/other work).

3. Very minor information leak for the host about their users, host owner can see how many relevant events their subscribers are getting back. Tbf that's no different from it seeing how many transactions that user is submitting etc., just a measure of their activity I guess and not tied to their acc address at all

4. How do events interact with the revelation period?

## Appendices

### Background - Ethereum Events Design

To help dApp developers design applications with a good UX, the ethereum developers invented the concept of "events" or "logs", which
are pieces of information emitted from smart contracts, which can be streamed in real time to external applications that
subscribe to them.

*Note that there is no constraint on data access, since all data is public.*

To better understand the anatomy of events, read this [blog](https://medium.com/mycrypto/understanding-event-logs-on-the-ethereum-blockchain-f4ae7ba50378)

#### Defining and emitting events

This is how an event is declared in a smart contract:

```solidity
event Transfer(address indexed from, address indexed to, uint256 value);
```

And this is how it is emitted:

```solidity
emit Transfer(from, to, amount);
```

#### Consuming events

A web app can subscribe to events by doing something like:

```javascript
var subscription = web3.eth.subscribe('logs', {
    address: '0x123456..', 
    topics: ['0x12345...']
}, function(error, result){
    if (!error)
        console.log(result);
});
- address - String|Array: An address or a list of addresses to only get logs from particular account(s).
- topics - Array: An array of values which must each appear in the log entries. The order is important, if you want to leave topics out use null, e.g. [null, '0x00...']. You can also pass another array for each topic with options for that topic e.g. [null, ['option1', 'option2']]
```

A web app might request:

- stream all "transfer" events where the 'to' field is my address
- stream all "transfer" events from the USDC contract

Apps can also request historic events starting from any block.

### Events and address topics

Let's analyse a couple of events from ERC20 and Uniswap, grouped by whether they contain address topics.

*A topic is a field which is marked as `indexed`.*

#### Events with account address topics

All the events in this section contain at least one account address topic.

```solidity
    event Transfer(address indexed from, address indexed to, uint256 value);
```

```solidity   
    /// @notice Emitted when the owner of the factory is changed
    /// @param oldOwner The owner before the owner was changed
    /// @param newOwner The owner after the owner was changed
    event OwnerChanged(address indexed oldOwner, address indexed newOwner);
```

```solidity
    event Swap(
    address indexed sender,
    uint amount0In,
    uint amount1In,
    uint amount0Out,
    uint amount1Out,
    address indexed to
    );
```

```solidity
    /// @notice Emitted when fees are collected by the owner of a position
    /// @dev Collect events may be emitted with zero amount0 and amount1 when the caller chooses not to collect fees
    /// @param owner The owner of the position for which fees are collected
    /// @param tickLower The lower tick of the position
    /// @param tickUpper The upper tick of the position
    /// @param amount0 The amount of token0 fees collected
    /// @param amount1 The amount of token1 fees collected
    event Collect(
        address indexed owner,
        address recipient,
        int24 indexed tickLower,
        int24 indexed tickUpper,
        uint128 amount0,
        uint128 amount1
    );
```

```solidity
    /// @notice Emitted when the collected protocol fees are withdrawn by the factory owner
    /// @param sender The address that collects the protocol fees
    /// @param recipient The address that receives the collected protocol fees
    /// @param amount0 The amount of token0 protocol fees that is withdrawn
    /// @param amount0 The amount of token1 protocol fees that is withdrawn
    event CollectProtocol(address indexed sender, address indexed recipient, uint128 amount0, uint128 amount1);
```

What all these events have in common is that the address topics like: `sender`, `recipient`, `owner`, `to`, etc, represent the
accounts which are affected by this transaction, and which are thus directly interested in it.

#### Events without account address topics

All the events in this section do not contain any account address topics.

```solidity
    /// @notice Emitted when a pool is created
    /// @param token0 The first token of the pool by address sort order
    /// @param token1 The second token of the pool by address sort order
    /// @param fee The fee collected upon every swap in the pool, denominated in hundredths of a bip
    /// @param tickSpacing The minimum number of ticks between initialized ticks
    /// @param pool The address of the created pool
    event PoolCreated(
        address indexed token0,
        address indexed token1,
        uint24 indexed fee,
        int24 tickSpacing,
        address pool
    );
```

```solidity
    /// @notice Emitted when a new fee amount is enabled for pool creation via the factory
    /// @param fee The enabled fee, denominated in hundredths of a bip
    /// @param tickSpacing The minimum number of ticks between initialized ticks for pools created with the given fee
    event FeeAmountEnabled(uint24 indexed fee, int24 indexed tickSpacing);
```

```solidity
    event Sync(uint112 reserve0, uint112 reserve1);
```

```solidity
    /// @notice Emitted when the protocol fee is changed by the pool
    /// @param feeProtocol0Old The previous value of the token0 protocol fee
    /// @param feeProtocol1Old The previous value of the token1 protocol fee
    /// @param feeProtocol0New The updated value of the token0 protocol fee
    /// @param feeProtocol1New The updated value of the token1 protocol fee
    event SetFeeProtocol(uint8 feeProtocol0Old, uint8 feeProtocol1Old, uint8 feeProtocol0New, uint8 feeProtocol1New);
```

What these events have in common is that they are not user-specific. They represent a general update of the smart 
contract. While they might contain address fields, but these are addresses of smart contracts.

### Events APIs in Ethereum, Geth, and common Web3 libraries

#### Official Ethereum events API

The Ethereum [JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/) spec defines seven methods for 
handling events:

* `eth_newFilter`
* `eth_newBlockFilter`
* `eth_newPendingTransactionFilter`
* `eth_uninstallFilter`
* `eth_getFilterChanges`
* `eth_getFilterLogs`
* `eth_getLogs`

In addition, there are two methods defined as part of the 
[subscription API spec](https://ethereum.org/en/developers/tutorials/using-websockets/#eth-subscribe):

* `eth_subscribe`
* `eth_unsubscribe`

There are three types of subscription that can be created:

* `newPendingTransactions`
* `newHeads`
* `logs`

#### web3.js events API

web3.js receives events in two ways:

* `eth_getLogs`
* `eth_subscribe`/`eth_unsubscribe`

It uses all three types of subscriptions (`newPendingTransactions`, `newHeads`, and `logs`).

#### ethers.js events API

ethers.js uses the majority of the APIs above (both from the JSON-RPC API, and the subscription API). The only ones it 
doesn't appear to use are:

* `eth_newFilter`
* `eth_newBlockFilter`
* `eth_getFilterLogs`

#### Geth events API

Geth has a `FilterAPI` that is registered under the `eth_` namespace and defines the following 10 API methods:

* `NewPendingTransactionFilter`
* `NewBlockFilter`
* `NewFilter`
* `GetLogs`
* `UninstallFilter`
* `GetFilterLogs`
* `GetFilterChanges`
* `NewPendingTransactions`
* `NewHeads`
* `Logs`

The first seven methods match up with those on the Ethereum JSON-RPC spec, while the remaining three 
(`NewPendingTransactions`, `NewHeads` and `Logs`) are used to power the three types of subscriptions listed above.
