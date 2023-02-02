---
---
# The Standard Obscuro Bridge

The standard Obscuro bridge is a decentralized asset bridge that uses a wrapped token mint and burn pattern. Under the hood it is based on the cross chain messaging protocol and exists entirely as a smart contract without the need of separate runnables or nodes.

## Contract Addresses

## Interacting with the contract

Users interact with the bridge contract using either the `sendERC20` or `sendNative` functions. In order to send ERC20 tokens they need to be whitelisted. 

The pattern of interacting on both layers with `sendERC20` is as follows: 
 1. Approve the bridge contract to act on your behalf for X amount of tokens using the standard ERC20 `approve()` method.
 2. Call `sendERC20` with the amount approved and the address of the recipient on the other side.
 3. Wait a tiny bit for the cross chain protocol to synchronize the state of the layers.
 4. Use your assets on the other layer.

Interacting with `sendNative` does not require approving an allowance, but instead just add value to the transaction (or evm call).
The value received by the bridge contract during the execution of `sendNative` is what will be logged as transfer to the other layer.

## Layer 1 To Layer 2 Specifics

The `ObscuroBridge.sol` contract is responsible for managing the layer 1 side of the bridge. It's the "bridge to Obscuro".
The `EthereumBridge.sol` contract is responsible for managing the layer 2 side of the bridge. It's the "bridge to Ethereum".

In order to bridge tokens over they need to be whitelisted. **Only accounts with the admin role can whitelist tokens!**

When an asset is whitelisted, the bridge internally uses the `publishMessage` call on the MessageBus which can be found within the `ManagementContract`. In the message that is published the bridge "tells" the other side of it, which resides on layer 2 that a token has been whitelisted. This in turn creates a wrapped version of the token on the other side. This version of the token can only be minted and burned by the layer 2 bridge.

 * Minting allows to create fresh funds on the L2 when they get locked on L1.
 * Burning allows to destroy supply on L2 in order to release it from the bridge on L1.

The protocol to bridge assets using `mint`/`burn` is based on the `MessageBus`'s `publishMessage` too. The `ObscuroBridge` tells the layer 2 that it has locked a certain amount of tokens and because of that address Y should receive the same amount of the wrapped token counterpart. This is when `mint` comes into play. On the other hand, the `EthereumBridge` burns tokens when received and tells the L1 that it has in fact done so. This is interpreted by the L1 smart contract as evidence to release locked funds.


## Security

The bridge inherits the `CrossChainEnabled` contract. It provides a modifier that can be attached to public smart contract functions in order to limit them. Normally you can limit a function to be called by `msg.sender`. With this new modifier you can limit a function to be called from the equivalent of `msg.crosschainsender`. 

The bridge links its two sides together - the ObscuroBridge and the EthereumBridge and limits the cross chain sender to only be the opposite contract on the `receiveAssets` functions. This means that only the `CrossChainMessenger` can call them when it is provided with a valid cross chain message produced by the bridge. As the cross chain sender is determined by `msg.sender` at the point of calling `publishMessage`, no one else but the bridge has the capability of queueing a message that will pass.


The layer 2 bridge works much the same way as the L1 bridge and uses the `MessageBus` to facilitate communication. 

## Building Bridges

The standard bridge was designed to work without requiring access to any private platform capabilities. This means that anyone can build a bridge using the same cross chain messaging API. There is no permissioning required to do so.

## Interface 

```solidity
interface IBridge {
    function sendNative(address receiver) external payable;

    function sendERC20(
        address asset,
        uint256 amount,
        address receiver
    ) external;
}
```