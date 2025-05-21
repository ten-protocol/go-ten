# Timestamping transactions in TEN

## Background 
Fast-paced games like "Aviator" require users to make quick decisions.
This doesn't work well on-chain unless latencies are very low.

In case of the popular "Aviator" game, the user has to decide in real-time when to exit with the winnings. 

A standard EVM smart contract does not have access to the timestamp of a transaction. The reason is that the "clock" of the user is not trusted, and a network like "Ethereum" is decentralised and there is no trusted "clock".

### External Timestamp oracle

A user wanting to timestamp a transaction would have to use a trusted "Oracle" to sign over a payload containing the tx hash and the timestamp. And then create another Ethereum transaction with this payload, and the smart contract somehow matching the previous action with this proof.
This approach is very complex with multiple on-chain transactions.

### Protocol timestamp oracle

The goal of this document is to propose a timestamp oracle embedded in the TEN protocol. 

A smart contract running on TEN will have access to an authenticated time window. 

The trusted authority is the "Sequencer" node. The Sequencer will make sure that the moment a transaction is added to the mempool is within the indicated time window.


## Game mechanics for a simple fully on-chain Aviator

When a new game starts, the smart contract reads the current block timestamp, adds 1s, and generates a random number that is the number of ms when the game ends (and the plane crashes).

Note 1: When the "start" transaction is processed, it emits an event that starts the UIs. 

Note 2: The "crash" time is private state.

Players will click "Exit" along the way. Their transactions will contain the "time window" when they made the action. 

The protocol will authenticate the time-window, and will only invoke the smart contract if the time-window is correct.

The contract will process transactions, and can use the time-window to determine whether the player has jumped before the crash.



## Implementation

### Estimate gas

The gas estimation should ignore the time window.

### Tx is submitted to the TEN Gateway

The GW will extract the time-window and check that there is enough time left to be accepted by the network. 

This is just a UX check so that it returns a nice error.

### Tx is submitted to the network and reaches a validator

A validator replicates the check done by the gateway, and also return an error.


### Tx reaches the sequencer

The sequencer replicates the above check and only adds to the mempool if `time.now` is in the declared `time-window`.

Note that transactions with an invalid time window will fail silently at this point. It is important that the gateway and the validator that received the tx before have communicated to the dApp already.

Note that the transaction is not included in the batch and is dropped.


### Tx is in a batch and executed by a validator

Transactions included in batches signed by the sequencer are assumed to have gone through the logic 

### Smart contract

A function that uses this feature will have as first parameter a `bytes32` parameter with a specified format:

`byte[0:1]="tw"`
`byte[2:17]=timestamp of the middle of the window`
`byte[18:29]=window size in ms`
`byte[30:31]="tw"`

```solidity

    function makeMove(bytes32 timeWindow, .. other stuff) public returns (bool) {
      uint from,to;
      (from,to)=extractTimeWindow(timeWindow); // utility that converts from the above format to a time window 
      // more stuff
    }


```
