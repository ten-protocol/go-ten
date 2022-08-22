# Events in Obscuro - Design document

## Scope

This is a design proposal for how to handle events in Obscuro.

It covers two aspects:

- the visibility rules for events.
- technical implementation details.

## Ethereum Events

To help dApp developers design applications with a good UX, the ethereum developers invented the concept of "events" or "logs", which
are pieces of information emitted from smart contracts, which can be streamed in real time to external applications that
subscribe to them.

To better understand the anatomy of events, read this [blog](https://medium.com/mycrypto/understanding-event-logs-on-the-ethereum-blockchain-f4ae7ba50378)

### Smart contracts

This is how an event is declared in a smart contract.

```
event Transfer(address indexed from, address indexed to, uint256 value);
```

And this is how it is emitted.
```
emit Transfer(from, to, amount);
```

### Consuming events

A web app can subscribe to events by doing something like:

```
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

Todo: Can you clarify in the doc whether the events are discarded after being distributed, or are stored for future
subscribers?

*Note that there is no constraint on data access, since all data is public.*

## Obscuro Design

In Obscuro, we aim maintain the same building blocks that are found in Ethereum: events and subscriptions, and will try
to implement the privacy concerns with as little disruption as possible.

### Event visibility rules

There are a couple of cases that must be considered in order to decide whether Alice is entitled to view an event:

1. The event was emitted by a smart contract as a result of executing a Tx sent by Alice.
2. The event was emitted as result of a Tx sent by Bob which is relevant to Alice. (See below for a definition of
   relevancy.)
3. The event was emitted as result of a Tx sent by Bob which is not relevant to Alice.

*Note that we assume that all events emitted during a transaction signed by Alice, will be visible to Alice.*  

The developer of a smart contract is responsible with making sure to not include any data that must be kept secret from
the transaction originator.
(The events are also included in the transaction receipt, which is available to the transaction submitter.)

#### Event relevancy

In Obscuro (inherited from Ethereum), end users can have multiple accounts. The account address is how accounts are
referenced.

Events are structured objects containing multiple entries (topics or data fields).

If we were designing events from scratch, with privacy in mind, we could add metadata to declare which address should be
able to view an event. Since we're trying to maintain the API of Ethereum unchanged, we'll try to infer this information
from the existing information available in the event, and to also allow the developers to achieve the desired outcome.

Let's analyse a couple of events from ERC20 and Uniswap, grouped by whether they contain address fields.

##### With end-user address topics

All the events in this section contain at least one end-user address topic.

*Note: a topic is a field which is marked as `indexed`*

```
    event Transfer(address indexed from, address indexed to, uint256 value);
```

```   /// @notice Emitted when the owner of the factory is changed
    /// @param oldOwner The owner before the owner was changed
    /// @param newOwner The owner after the owner was changed
    event OwnerChanged(address indexed oldOwner, address indexed newOwner);
```

```
    event Swap(
    address indexed sender,
    uint amount0In,
    uint amount1In,
    uint amount0Out,
    uint amount1Out,
    address indexed to
    );
```

```
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

```
    /// @notice Emitted when the collected protocol fees are withdrawn by the factory owner
    /// @param sender The address that collects the protocol fees
    /// @param recipient The address that receives the collected protocol fees
    /// @param amount0 The amount of token0 protocol fees that is withdrawn
    /// @param amount0 The amount of token1 protocol fees that is withdrawn
    event CollectProtocol(address indexed sender, address indexed recipient, uint128 amount0, uint128 amount1);
```

What all these events have in common is that the address topics like: `sender`, `recipient`, `owner`, `to`, etc, represent the 
accounts which are affected by this transaction, and which are thus directly interested in it.


##### Without end-user address fields

```
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

```
    /// @notice Emitted when a new fee amount is enabled for pool creation via the factory
    /// @param fee The enabled fee, denominated in hundredths of a bip
    /// @param tickSpacing The minimum number of ticks between initialized ticks for pools created with the given fee
    event FeeAmountEnabled(uint24 indexed fee, int24 indexed tickSpacing);
```

```
    event Sync(uint112 reserve0, uint112 reserve1);
```

```
    /// @notice Emitted when the protocol fee is changed by the pool
    /// @param feeProtocol0Old The previous value of the token0 protocol fee
    /// @param feeProtocol1Old The previous value of the token1 protocol fee
    /// @param feeProtocol0New The updated value of the token0 protocol fee
    /// @param feeProtocol1New The updated value of the token1 protocol fee
    event SetFeeProtocol(uint8 feeProtocol0Old, uint8 feeProtocol1Old, uint8 feeProtocol0New, uint8 feeProtocol1New);
```

What these events have in common is that they are not user-specific. They represent a general update of the smart contract.

##### Visibility Rules

Users should be able to request and read all events that are relevant to them.

The rules we propose are:

1. An event is considered relevant to all account owners whose addresses are used as topics in the event.
2. In case there is no address in a topic, then the event is considered a lifecycle event, and thus relevant to everyone.
3. An event is considered relevant to the signer of the transaction that created it. 

The purpose for these rules is to be simple, clear, intuitive, and to work as good as possible with the existing contracts.  

In case the rules above are not providing the desired functionality, the developer can adjust the topics.

For example, if one of the lifecycle events should only be visible to the administrators, the developer can add that address as a topic.

In case an event with an address field that should not contribute to relevancy, the developer can remove the `"indexed"` and thus the event will become invisible to that user.

*Note: It is easy to detect if a value is an address since it has 20 non-zero bytes. There is no risk that a collision occurs between an end-user address, and some other field because there are 2^160 possibilities.*

### Technical implementation

The task is to implement the visibility rules described above without changing the query and subscription API from a user's point of view.

#### Constraints and Considerations 

We already have a tool called the "Wallet Extension", which acts as a proxy between the wallet and the obscuro node, and manages viewing keys.

- Applications will connect to the "wallet extension", which will translate the plain web3 "Subscribe"
  call into an encrypted Obscuro compatible one. The stream of received events will be decrypted automatically with the appropriate viewing keys.


- Events should not leave the enclave space unencrypted or encrypted with a non-relevant account key. Transactions are 
  executed inside a secure enclave, and events emitted during that, need to be collected, filtered, and encrypted before being returned from the enclave.
  Optimisations need to be created as the load on the enclave could be significant. 

- An account should be able to monitor only the events relevant to itself, and not subscribe to anything else. 
  Basically, subscriptions need to be authenticated. This is not desireable,
  because someone could setup a subscription to monitor well-known addresses, and receive useful information, even if
  they cannot decrypt the actual event. 
  They could determine  for example when a high profile individual has transferred some ERC20, even if they wouldn't know 
  how much or to whom.

  
#### Prerequisites

An event `E1` is emitted after executing transaction `TX1` with a from address `A1`.
`E1` is formed of multiple topics, 2 of which are addresses (`A2` and `A3`), and there are some other random values.

An event `E2` is emitted by a transaction `TX2` from an address `A4`.
`E2` is formed of multiple topics, none of which are addresses.

User `U1`- owner of `A1` subscribes to all events.

#### Implementation

1. The Obscuro `Subscription` call, and the `Event query` call must take a list of signed owning accounts. Each account must be signed with the latest
   viewing key( to prevent someone from asking random events, just to leak info). The call will fail if there are no
   viewing keys for all those accounts. 

   Note: This is possible because the subscription call is implemented on Obscuro, and
   made by the wallet_extension, so it doesn't have to be compatible with Ethereum. For our RPC client it would an
   authenticated subscription.


2. The "Obscuro Host" is responsible in setting up the subscriptions and dispatching the events it receives from the enclave.


3. Upon ingesting a rollup included in a block, the enclave runs all the usual filtering logic on this event, and determines if there 
   is any subscription that requested it. 


4. Then, there is an extra step (inside the enclave as well) to determine whether the event is relevant to the account which is
   authenticated for that subscription.
 
   Based on the rules described above, this can be done by running an intersection between:
   - all accounts found in the event topics plus the `from` of the transaction. E.g.: `for E1 -> (A1,A2,A3)`
   - and the accounts from the subscription. `U1 -> A1`

    If there is at least one common element, the event is encrypted and returned in the `BlockSubmissionResponse` with metadata.

    *Note: This extra step is skipped for rvent `E2`, which is returned on every subscription where it matches, encrypted with the first viewing key*


5. The encrypted event is streamed from the host and then sent to the wallet extension, where it is decrypted, and streamed further to the App.


### Security and Usability of the proposed design

App developers will be able to use the existing libraries unchanged, as long as they connect through a wallet extension with registered viewing keys.

Depending on the subscription, the results might be different from those returned on a normal Ethereum network, because the user might not have the right to see certain private data. 

Smart contract developers need to spend a few minutes to think about whether an event can be seen by the entity who it is relevant to, or whether it leaks data.
There is no new syntax to learn.

The data access protections of smart contracts will prevent another smart contract interacting with it from extracting information and leaking it as an event.  

The fact that the wallet extension adds signed accounts to each subscription request, makes it impossible for a user to request the events of another user.

An ERC20 transfer from Alice to Bob will show up on Bob's UI if he is subscribed to it, but will not show on Charlie's UI. 