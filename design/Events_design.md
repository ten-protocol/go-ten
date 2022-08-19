# Events in Obscuro - Design document

## Scope

This is a design proposal for how to handle events in Obscuro.

## What are events in Ethereum

### Smart contracts can emit events which are streamed from a node and generally used by UIs.

For example:

```emit Transfer(from, to, amount);```


### A web app can subscribe to events by doing something like:

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


### How events work 

The query made on the UI side is transformed in a server-side query on the node, and there is some logic after a tx is executed and events are emmitted, to match them against the filters requested by users and distribute them to the requester.(? - need to confirm this is 100% how it works) .

*Note that there is no constraint on data access, since all data is public.*


## Requirements and Problems in Obscuro

- The high level requirement is to determine which account owner can view which events, and make sure, through encryption that nobody else can see those events.
- Ideally we want to achieve this without any change to the API.


#### Problem1: How do we authenticate event subscriptions? 
We have a wallet extension with viewing key support, which we have to use ideally, to avoid creating another mechanism.

#### Problem2: Assuming the enclave knows who is subscribed, how will the filtering work?

E.g: Alice is requesting to listen to all transfer events. (this is a perfectly reasonable request on a public chain)

There are 2 cases:
1. The event was emitted by a smart contract as a result of executing a tx sent by Alice.
2. The event was emitted as result of a Tx sent by Bob but it is somehow relevant to Alice.


#### Problem 3: How to encrypt events before they leave the enclave

#### Problem 4: The mechanism should not leak general information on what activity is happening.


## Proposed solution

### Prerequisites
An event `E1` is emitted after executing transaction `TX1` with a from address `A1`.
`E1` is formed of multiple topics, 2 of which are addresses (`A2` and `A3`), and there are some other random values.

An event `E2` is emitted by a transaction `TX2` from an address `A4`.
`E2` is formed of multiple topics, none of which are addresses.


### Solution

1. We modify the subscription call to accept a list of owning accounts. 
Each account must be signed with the latest viewing key, to prevent someone from asking random events, just to leak info. 
The call will fail if there are no viewing keys for all those accounts.
This is possible because the subscription call is implemented on Obscuro, and made by the wallet_extension, so it doesn't have to be compatible with Ethereum.
For our RPC client it would an authenticated subscription.

2. We need to pass into the enclave the "Subscription" objects, and create an rpc streaming for each. 

3. The enclave runs all the usual filtering logic on this event, and determines if there is any subscription that requested it.

4. Then, there is an extra step (inside the enclave as well) to determine whether that party is allowed to view that event. 

This can be done by running an intersection between:
  - all accounts found in the event topics plus the `from` of the transaction: `E1->(A1,A2,A3)` 
  - and the accounts from the subscription. 

If there is at least one element, the event is encrypted with the first of the result.

6. For event `E2`, that has no addresses in the topics, there are a couple of options:

    1. Because it is not address specific, we consider it a public event, and send it to everyone who requests it. Basically, we send it encrypted with the first viewing key from any subscription that matches it.
This has the advantage that the UIs will keep working as expected for existing contracts, but the disadvantage that it might leak unexpected data.   
    2. We only send it to the `from`. This is more secure, but will impact usability.

My personal preference is for Option 1. Which is in line with the other decisions around not changing the behaviour of public variables, and documenting it for developers.

8. The encrypted event is bubbled out of the enclave on the right subscription channel, and then sent to the wallet extension.


### Scenarios

#### 1. Curious user deploys a contract that intends to leak Erc20 balances.
This is prevented by the protections in place in the ERC20 contract, who will only allow the initiator of a tx to read the balance.


#### 2. Curious user listens for all possible events.
Under Option1. They will get back all events that don't have addresses.
