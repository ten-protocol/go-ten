# Events in Obscuro - Design document

## Scope

This is a design proposal for how to handle events in Obscuro.

It covers two aspects:

- the visibility rules for events.
- technical implementation details.

## Background - Ethereum Events Design

To help dApp developers design applications with a good UX, the ethereum developers invented the concept of "events" or "logs", which
are pieces of information emitted from smart contracts, which can be streamed in real time to external applications that
subscribe to them.

To better understand the anatomy of events, read this [blog](https://medium.com/mycrypto/understanding-event-logs-on-the-ethereum-blockchain-f4ae7ba50378)

### Smart contracts

This is how an event is declared in a smart contract.

```solidity
event Transfer(address indexed from, address indexed to, uint256 value);
```

And this is how it is emitted.
```solidity
emit Transfer(from, to, amount);
```

### Consuming events

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

### How events work

`Note: this section might not be 100% accurate.`

The query made on the UI side is transformed in a server-side query on the node, and there is some logic after a tx is
executed and events are emmitted, to match them against the filters requested by users and distribute them to the
requester.

*Note that there is no constraint on data access, since all data is public.*

## Obscuro Events Design

In Obscuro, we aim to maintain the same building blocks that are found in Ethereum: events and subscriptions, and will try
to implement the privacy concerns with as little disruption as possible.

### Event types

There are a couple of cases that must be considered in order to decide whether Alice is entitled to view an event:

1. The event was emitted by a smart contract as a result of executing a Tx sent by Alice.
2. The event that is relevant to Alice was emitted as a result of a Tx sent by Bob. (See below for a definition of
   relevancy.)
3. The event that is not relevant to Alice was emitted as a result of a Tx sent by Bob.

### Event relevancy

In Obscuro (inherited from Ethereum), end users can have multiple accounts. The account address is how accounts are
referenced.

*Note: Smart contracts also have accounts referenced by an address. Given an account address, we can query the "code" property 
to determine whether it is an end user or a contract.*

Events are structured objects containing multiple entries (topics or data fields).

If we were designing events from scratch, with privacy in mind, we could add metadata to declare which address should be
able to view an event. Since we're trying to maintain the API of Ethereum unchanged, we'll try to infer this information
from the existing information available in the event, and to also allow the developers to achieve the desired outcome.

Let's analyse a couple of events from ERC20 and Uniswap, grouped by whether they contain address fields.

#### With end-user address topics

All the events in this section contain at least one end-user address topic.

*Note: a topic is a field which is marked as `indexed`*

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

#### Without end-user address fields

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

What these events have in common is that they are not user-specific. They represent a general update of the smart contract.

*Note that they might contain address fields, but these are addresses of smart contracts.*

### Event visibility rules

Users should be able to request and read all events that are relevant to them. By relevant, we mean that the user was 
somehow involved in the transaction that emitted that event, and this event might be of interest to them.

The implicit rules we propose are:

1. An event is considered relevant to all account owners whose addresses are used as topics in the event.
2. In case there is no user address in a topic, then the event is considered a lifecycle event, and thus relevant to everyone.

The purpose for these rules is to be simple, clear, intuitive, and to work as good as possible with the existing contracts.

There are several edge-cases:

1. The event contains addresses, but they are contract address, as in the ``PoolCreated`` event from above
2. The event contains a topics that looks like an address, but is not

To handle this, for each potential account address in the topics, the VM must check:

* If there is associated contract code, using the `GetCode` function
* If it has a zero nonce, using the `GetNonce` function (only account and smart contract addresses have non-zero nonces)

If either of these is true, the event will fall under rule 1.

### Adjusting the event visibility rules

In case the rules above are not providing the desired functionality, the developer will have a couple of options to adjust visibility. 

For example, if one of the lifecycle events should only be visible to the administrators, the developer can add that address as a topic.

In case an event contains a user address topic that should not contribute to relevancy, the developer can remove the 
`"indexed"` and thus the event will become invisible to that user. This might not be ideal in case this address has to 
be used for subscribing. 

As the ultimate flexible mechanism we propose a programmatic way to determine whether a requester is allowed to view an event. 

If the implicit rules are not satisfactory, the smart contract developer can define a view function called 
``eventVisibility``, which will be called by the Obscuro VM behind the scenes.

```solidity
   // If declared in a smart contract, this function will be called by the Obscuro VM to determine whether the requester
   // address is allowed to view the event. 
   function eventVisibility(address requester, bytes32[] memory topics,  bytes32[] memory data) external view returns (bool){
       // based on the data from the event, passed in as an array of topic and data (which is the internal VM representation) 
       // calculate whether the requester address is entitled to view the event
       // - the first element in the topics array is the hash of the event signature
       // Note: the developer must maintain a list of the event hashes they want to handle programmatically 
        return ; 
    }
```    

To determine the visibility of an event, the Obscuro VM will do the following:

1. call the `eventVisibility` with the event being requested and the requester.
2. If the function exists and returns 'true', then return the event. If it returns `false`, then the event is invisible.
3. If the function does not exist, apply the implicit rules.

## Obscuro events technical implementation

The task is to implement the visibility rules described above without changing the query and subscription API from a user's point of view.

### Constraints and Considerations 

We already have a tool called the "Wallet Extension", which acts as a proxy between the wallet and the obscuro node, and manages viewing keys.

- Applications will connect to the "wallet extension", which will translate the plain web3 "Subscribe"
  call into an encrypted Obscuro compatible one. The stream of received events will be decrypted automatically with the appropriate viewing keys.

- Events should not leave the enclave space unencrypted or encrypted with a non-relevant account key. Transactions are 
  executed inside a secure enclave, and events emitted during that, need to be collected, filtered, and encrypted before being returned from the enclave.
  Optimisations need to be created as the load on the enclave could be significant. 

- An account should be able to monitor only the events relevant to itself, and not subscribe to anything else. 
  Basically, subscriptions need to be authenticated. Otherwise, someone could setup a subscription to monitor well-known addresses, 
  and receive useful information, even if they cannot decrypt the actual event. 
  They could determine  for example when a high profile individual has transferred some ERC20, even if they wouldn't know 
  how much or to whom.

- Events included in transaction receipts should be filtered to only include events which are visible to the transaction submitter.

### Implementation

#### Obscuro enclave

##### Creating and deleting subscriptions

The enclave exposes two methods, one to add a new logs subscription, and one to delete a logs subscription.

The method to create a subscription takes a `LogSubscription` object, defined as below. This object is encrypted with 
the enclave's private key, to prevent attackers from eavesdropping on the creation of subscriptions.

```
LogSubscription {
    Account   // The account address the events relate to.
    Signature // A signature over the account address using a private viewing key.
    Filter    // A subscriber-defined filter to apply to the stream of logs.
}
```

The signature ensures that the client is authorised to create a subscription for the given account. This prevents
attackers from creating subscriptions to analyse the pattern of logs. Each new logs subscription is assigned a unique 
subscription ID that is returned to the host.

The method to delete a subscription takes the subscription ID of the subscription to be deleted. This method is not 
authenticated, so an attacker who discovered the subscription ID could request its deletion.

##### Block ingestion

Each time the host sends the enclave a new block to ingest, the enclave responds with a block submission response. 
This response will be extended to include a mapping from subscription IDs to the associated logs, where the logs are 
encrypted with the viewing key corresponding to the `LogSubscription.Account`. By associated logs, we mean any logs 
from any transactions in the current or historical blocks that meet the two following criteria:

* They match the `LogSubscription.Filter`
* They pass the relevancy test described above for the user address in `LogSubscription.Account`

This means that the set of logs for a given subscription can become very large. For example, a subscription with no 
filter would return all logs passing the relevancy test for that user since the start of time.

##### Transaction receipts

Logs are also included in transaction receipts. Whenever the host requests a transaction receipt, the enclave filters 
out the receipt's logs to only include those that pass the relevancy test (using the transaction's sender as the user 
address).

#### Obscuro host

For each incoming logs subscription request via RPC, the host has to do two things:

1. Request the creation of the new subscription on the enclave
2. Create a Geth `rpc.Subscription` object to route any new logs for the subscription back to the client after each 
   block ingestion

##### Enclave subscription

For (1), the host takes a log subscription object as a parameter, and forwarding it to the enclave via the API 
described above. As above, the log subscription object is encrypted with the enclave's private key, to prevent 
attackers from eavesdropping on the creation of subscriptions.

##### Client subscription

For (2), we reuse Geth's `PublicFilterAPI`, passing in a custom `Backend` object. The `Backend` object has a 
`SubscribeLogsEvent` method, which takes as a parameter a channel, and is responsible for sending any new logs onto 
that channel. We implement this method by pushing a new log onto the channel whenever a new log is returned from the 
enclave via a block submission response.

The host then creates individual logs subscriptions using the `PublicFilterAPI.Logs` method, which takes a filter 
criteria and automatically applies it to the contents of the master list returned by `SubscribeLogsEvent` to determine 
the set of logs to return on that subscription.

There are two key limitations that had to be overcome here:

1. The log events contained in the block submission response are encrypted by the enclave, and are thus of type 
   `[]byte`, whereas the channel provided to `SubscribeLogsEvent` expects Geth's `types.Log` objects. To overcome this, 
   the host performs an additional step where a fake "wrapper" log object is created, and the encrypted log bytes are 
   placed into the `data` field of that wrapper log
2. How do we prevent `PublicFilterAPI.Logs` from automatically returning all events over the subscription? To achieve 
   this, when creating the wrapper log, we include the subscription ID as a topic. Then when we set up the logs 
   subscription using `PublicFilterAPI.Logs`, we pass in a fake filter that filters based on that subscription ID as 
   the topic. This ensures that only the logs for the specific subscription ID are returned. (Note that the "real" 
   filter has already been passed encrypted to the enclave, and applied before returning the results as part of the 
   block submission response)

#### Obscuro encrypted RPC client

Due to their sensitive nature, logs subscription requests and responses must pass through the encrypted RPC client.

The encrypted RPC client only handles logs subscriptions via the `eth_subscribe` and `eth_unsubscribe` APIs (see 
[here](https://ethereum.org/en/developers/tutorials/using-websockets/#eth-subscribe)). A consequence of this is that 
events are only available in Obscuro over a websocket connection.

In response to the incoming `eth_subscribe` request, the client creates a `logs` subscription to the host by making a
`rpc.Client.Subscribe` call via the embedded Geth client. It passes as a parameter a `LogSubscription` object encrypted 
with the enclave's public key to protect the request from eavesdroppers, setting the `Account` to the client's 
account and generating the required signature.

For each received event, the encrypted RPC client must retrieve the encrypted log bytes from the `data` field and 
decrypt them with the corresponding private key before returning the log events to the user. It then pushes the 
decrypted event onto a separate channel listened to by the client.

### Security and usability of the proposed design

App developers will be able to use the existing libraries unchanged, as long as they connect through a wallet extension with registered viewing keys.

Depending on the subscription, the results might be different from those returned on a normal Ethereum network, because the user might not have the right to see certain private data. 

Smart contract developers need to spend a few minutes to think about whether an event can be seen by the entity who it is relevant to, or whether it leaks data.
There is no new syntax to learn, just to be aware about the visibility rules. In case the default intuitive rules do not satisfy the requirements,
the developer can write a function to precisely control visibility.

The data access protections of smart contracts will prevent another smart contract interacting with it from extracting information and leaking it as an event.  

The fact that the wallet extension adds signed accounts to each subscription request, makes it impossible for a user to request the events of another user.

An ERC20 transfer from Alice to Bob will show up on Bob's UI if he is subscribed to it, but will not show on Charlie's UI.

### Open implementation questions

1. What's the lifespan for a subscription in the enclave, unsubscribe option + an expiry?

2. Still pretty worried about perf and DoS potential, but we can test for it and try to optimise/prioritise traffic etc. Maybe we need read-only nodes to services these requests in production (normal nodes but the host doesn't give it any transactions/other work).

3. Very minor information leak for the host about their users, host owner can see how many relevant events their subscribers are getting back. Tbf that's no different from it seeing how many transactions that user is submitting etc., just a measure of their activity I guess and not tied to their acc address at all

4. Todo: Can you clarify in the doc whether the events are discarded after being distributed, or are stored for future subscribers?

5. How do events interact with the revelation period?

### Alternatives considered

#### All events are public by default 

This is not possible as it breaks the most fundamental contact, the `ERC20`, which contains the `Transfer` event.
If all events were public by default then, we either break the ERC20 api by removing the event, or we lose privacy

#### Add a third visibility rule giving the signer of the transaction total visibility

Rule 3: A signer of a transaction can view all events emitted during the execution of that transaction.

This rule adds another dimension to the reasoning process, because there is an implicit user to whom the event is relevant.
It also reduces flexibility in sending lifecycle events to administrators.

Note: Adding this rule simplifies transaction receipts.

## Appendices

### Events in Ethereum, Geth, and common Web3 libraries

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

## Geth events implementation

Obscuro gets access to Geth's pub/sub functionality for free because it reuses Geth's RPC layer. This PR shows a 
working example of a subscription API and client in Obscuro: https://github.com/obscuronet/go-obscuro/pull/604.

The events API and supporting code is implemented in the `eth/filters` package. It's possible we could reuse this code 
wholesale. We could use the `func NewFilterAPI(backend Backend, lightMode bool, timeout time.Duration) *FilterAPI` 
constructor, passing in our own implementation of the `filters.Backend` interface:

```
type Backend interface {
	ChainDb() ethdb.Database
	HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error)
	HeaderByHash(ctx context.Context, blockHash common.Hash) (*types.Header, error)
	GetReceipts(ctx context.Context, blockHash common.Hash) (types.Receipts, error)
	GetLogs(ctx context.Context, blockHash common.Hash) ([][]*types.Log, error)
	PendingBlockAndReceipts() (*types.Block, types.Receipts)

	SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
	SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription
	SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription
	SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription
	SubscribePendingLogsEvent(ch chan<- []*types.Log) event.Subscription

	BloomStatus() (uint64, uint64)
	ServiceFilter(ctx context.Context, session *bloombits.MatcherSession)
}
```
