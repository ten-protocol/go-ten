---
---
# The Standard Ten Bridge

The standard Ten bridge is a trustless and decentralised asset bridge that uses a wrapped token mint and burn pattern. Under the hood it is based on the cross chain messaging protocol and exists entirely as a smart contract without the need of separate runnables or nodes.

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

The `TenBridge.sol` contract is responsible for managing the layer 1 side of the bridge. It's the "bridge to Ten".
The `EthereumBridge.sol` contract is responsible for managing the layer 2 side of the bridge. It's the "bridge to Ethereum".

In order to bridge tokens over they need to be whitelisted. **Initially only accounts with the admin role can whitelist tokens!** When the protocol has matured this whitelisting functionality will change. If you want your token to be whitelisted, get in touch with us through our discord.

When an asset is whitelisted, the bridge internally uses the `publishMessage` call on the `MessageBus` contract which is deployed and exposed by the `ManagementContract` on layer 1. In the message that is published the bridge "tells" the other side of it, which resides on layer 2 that a token has been whitelisted. This in turn creates a wrapped version of the token on the other side. This version of the token can only be minted and burned by the layer 2 bridge. Notice that when the bridge "tells" its counter part, the message is not automatically delivered by Ten. To automate the process one needs to have a relayer in place that will do it automatically and a system in place that funds the gas costs of said relayer. As the network grows general purpose relayers might pop up. If you are developing or have developed such a relayer, contact us on discord to get it listed.

 * Minting allows to create fresh funds on the L2 when they get locked on L1.
 * Burning allows to destroy supply on L2 in order to release it from the bridge on L1.

The protocol to bridge assets using `mint`/`burn` is based on the `MessageBus`'s `publishMessage` too. The `TenBridge` tells the layer 2 that it has locked a certain amount of tokens and because of that address Y should receive the same amount of the wrapped token counterpart. This is when `mint` comes into play. On the other hand, the `EthereumBridge` burns tokens when received and tells the L1 that it has in fact done so. This is interpreted by the L1 smart contract as evidence to release locked funds. This message once again needs to be relayed for the process to work, but by a separate L1 relayer.


## Security

The bridge inherits the `CrossChainEnabledTen` contract. It provides a modifier `onlyCrossChainSender` that can be attached to public smart contract functions in order to limit them. Normally you can limit a function to be called by `msg.sender`. With this new modifier you can limit a function to be called from the equivalent of `msg.crosschainsender`. In the `TenBridge` we limit the cross chain sender to the remote bridge only like this `onlyCrossChainSender(remoteBridgeAddress)`.

The bridge initialization phase links the two bridge contracts (TenBridge, EthereumBridge) together. This linking is finalized when both contracts know the remote address on the opposite layer. Using those remote addresses, the bridges limit incoming `receiveAsset` messages to only be callable by the `CrossChainMessenger` contract sitting on the same layer. This `CrossChainMessenger` contract provides the necessary context to the bridge about the sender of the cross chain message. This allows to validate message is coming from the `remoteBridgeAddress`. The sender address is determined inside of `publishMessage` by taking the `msg.sender`. This `msg.sender` is put inside the metadata of each message. 

The result of this setup is that `receiveAsset` is only callable by having a valid cross chain message coming from the correct contract. Anything else will get rejected and cause the call to revert.

The layer 2 bridge works much the same way as the L1 bridge and uses the `MessageBus` to facilitate communication. 

## Building Bridges

The standard bridge was designed to work without requiring access to any private platform capabilities. This means that anyone can build a bridge using the same cross chain messaging API. There is no permissioning required to do so. We encourage you to build novel bridges and we'd love to assist you in any way possible. You are welcome at our discord if you have questions, need support or just want to discuss your bridge idea!

## Interface 

This is the interface that can be used in order to interact with the bridge. 

`sendNative` is where you send raw ethereum as a call value. Said value received by the bridge will be forwarded to the `receiver`.
`sendERC20 ` is the previously described function that allows for moving tokens.

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