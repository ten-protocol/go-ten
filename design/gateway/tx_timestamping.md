# TEN Gateway - as a timestamp oracle

## Background 
Fast paced games like "Aviator" require users to make quick decisions.
This doesn't work well on-chain unless latencies are very very low.

In case of Aviator, the user has to decide in real-time when to exit with the winnings. 

## Ideal game mechanics for on-chain Aviator

When a new game starts, the smart contract reads the current block timestamp, adds 1s, and generates a random number that is the number of ms when the game ends.
When the start transaction is processed, it emits an event that starts the UIs. 
People will click "Exit" along the way. Their transactions will contain the timestamp when they made the action. The timestamp translates into a multiplier with a formula.

The contract will process transactions, and based on the timestamp will reward or not.

This game mechanics can only work on a chain with private shared state (the target timestamp is invisible) and if the timestamps are trusted.


## TEN Gateway - as a timestamp oracle

TEN already has the private state. 

The TEN Gateway is a service running in SGX, which manages Viewing Keys and Session Keys on behalf of users.

Because it is the first point of contact for a user, it can act as a Timestamp oracle.

Note: SGX itself cannot control the clock of a server. The trust element is more around it being a neutral central point.

**The proposal is for the TEN GW to attest the timestamps from the transactions.**


## Game implementation

- The TEN GW will generate a timestamp certificate ( a key-pair generated encrypted).
- The cert will be exposed on an rpc endpoint.
- The "Aviator" (or similar game) will have an admin interface and will add this cert to the contract as a "trusted timestamp oracle" (TTO)
- Session Keys will be signed with this timestamp certificate (TCS). This is a signature by the GW over the "from" address.
- The TCS will be submitted to the game smart contract.
- The contract can check the signature against the TTO and will know to trust the timestamp of transctions from that "from".
- On receiving the TX, the TEN enclave needs to do something(?) to extract the timestamp declared inside the TX and compare it with the timestamp of the "Transaction" object. And reject if it's more then x ms older.
- On receiving a tx, the contract will read the "from" and compare it against the trusted addressed that registered to play.
- 




