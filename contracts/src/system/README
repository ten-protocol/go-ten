# Public Callbacks Contract

The PublicCallbacks contract provides a mechanism for registering and executing callbacks in a gas-efficient manner. It is designed to be used by contracts that need to hide their gas cost whenever executing something and make the state change immutable if proxied.

## Key Features

- Register callbacks with associated value for gas payment
- Automatic gas refunds for unused gas
- Ability to reattempt failed callbacks
- Gas-limited execution to prevent abuse

## Usage

Any contract can call `registerCallback` to register a callback. The callback is built from the msg.sender and encoded data to be passed to the msg.sender when doing a callback. The value paid for is converted to gas and refunded if unused. We use the base fee at the time of registration to calculate the gas refund, so there is no issue of the price changing between registration and execution. Execution also uses no baseFee unlimited call.

## Internal

The contract uses a queue made out of a mapping and two uints. One points to where callbacks are added and the other lags behind pointing to the oldest callback.
The synthetic call DOES NOT fail if the underlying callback fails. Instead for now it gifts the stored value to coinbase and does not delete the callback, allowing for reattempting externally with whatever gas chosen. This might be a bit of a security risk, but its a failsafe as contracts normally do not have custom recovery logic if a callback fails.
