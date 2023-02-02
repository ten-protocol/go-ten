---
---
# Obscuro Cross Chain Messaging

Obscuro is equipped with a cross chain messaging protocol that enables sending and receiving messages across layers securely and in an entirely decentralized fashion.

The core contract that provides this functionality is the `MessageBus`. It exists both on L1 and L2. In the L1 you can find it under the management contract whilst on the L2 it is created and managed by the enclave.

## How It Works

Users call the `publishMessage` function on the `MessageBus`. This function emits an event that records what message is being published, who is calling the method and some other bits. This event later on gets picked on by the enclave which ensures that the message is available on the counterpart `MessageBus`.

When the message is made available, users can call the `verifyMessageFinalized` method in the `MessageBus` on the receipient layer. This method will return `true` if the message was indeed sent on the other layer as presented. Same `msg.sender`, same `payload`, same `nonce`.

This allows to have a behaviour much like querying questions - **'Has the bridge on address 0xAB0FF received 25 WETH tokens?'**. If the bridge on this address has called `publishMessage` saying 'I have received 25 WETH tokens?` the query will return true. 

When messages are published on the Obscuro layer (L2) the transport to L1 is done by the management contract upon rollup submission. Messages delivered however need to wait for the challenge period of the rollup before being considered final. This ensures that rollups along with the messages they carry can be challenged.

## Advanced capabilities

The `MessageBus` provides a way to inspect delivered messages and the time they become valid at. This is done using the `getMessageTimeOfFinality` method. This can be done to provide advanced functionality, for example withdrawal bonds - when running a validator node one can ensure that a rollup is valid before the challenge period expires and give the pending funds early whilst taking control of the ones to be released from the bridge for example.

## Security

The enclave only listens to the events on the contract address that is bound to the management contract. Messages are bound to L1 blocks and in the event of reorganization the state will be recalculated. L2 messages are bound to L1 block hashes too, resulting in a system where rollups that update the state can only be applied to the correct fork. If there is a mismatch, the state update will be rejected by the management contract.

The enclave controlls the keys for the L2 contract. Even if a hacker acquires the private key and gains the ability to sign correct transactions that would normally store messages on the L2, the enclave will reject any such transactions coming from this key that arrive on the normal route. The only acceptable way for the Obscuro protocol is for those transactions to be generated internally based on events. This means that every enclave regenerates those transactions when processing and validating state.

## Interface

```solidity
interface Structs {
    struct CrossChainMessage {
        address sender; // The contract/address which called the publishMessage on the message bus.
        uint64  sequence; // The sequential index of the message for the sending address.
        uint32  nonce; // Provided by the smart contract.
        uint32  topic; // Can be used to separate messages and provide basic versioning.
        bytes   payload; // The actual encoded message.
        uint8   consistencyLevel;
    }
}

interface IMessageBus {
    function publishMessage(
        uint32 nonce,
        uint32 topic,
        bytes calldata payload, 
        uint8 consistencyLevel
    ) external returns (uint64 sequence);

    function verifyMessageFinalized(
        Structs.CrossChainMessage calldata crossChainMessage
    ) external view returns (bool);
    
    function getMessageTimeOfFinality(
        Structs.CrossChainMessage calldata crossChainMessage
    ) external view returns (uint256);
}
```