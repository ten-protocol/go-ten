---
---
# Ten Cross Chain Messaging

Ten is equipped with a cross chain messaging protocol that enables sending and receiving messages across layers securely and in an entirely decentralized fashion.

The core contract that provides this functionality is the `MessageBus`. It exists both on L1 and L2. In the L1 you can find it under the management contract whilst on the L2 it is created and managed by the protocol.

## How It Works

Users call the `publishMessage` function on the `MessageBus`. This function emits an event that records what message is being published, who is calling the method and some other bits. This event later on gets picked up by the protocol which ensures that the message is available on the counterpart `MessageBus`.

When the message is made available, users can call the `verifyMessageFinalized` method in the `MessageBus` on the recipient layer. This method will return `true` if the message was indeed sent on the other layer as presented. Same `msg.sender`, same `payload`, same `nonce`.

This allows to have a behaviour much like querying questions - **'Has the address 0xAAAA.. received 25WETH tokens on the bridge with address 0XAB0FF?'**. If the bridge on this address has called `publishMessage` saying **'I have received 25 WETH tokens with recipient 0xAAAA.`** the query will return true. 

When messages are published on the Ten layer (L2) the transport to L1 is done by the management contract upon rollup submission. Messages delivered however need to wait for the challenge period of the rollup before being considered final. This ensures that rollups along with the messages they carry can be challenged. This is the logical equivalent of challenge period for optimistic rollups.

## Advanced capabilities

The `MessageBus` contract provides a way to query non finalized delivered messages. Those are the messages that are still within the time window allowing for a rollup to be challenged. Furthermore the contract gives the time they will become final at. This is done using the `getMessageTimeOfFinality` method.

Building on top of this message "preview" functionality allows for dApps that provide faster than challenge period finality. For example, a contract can provide the withdrawn amount from a bridge immediately with a fee. Doing this will transfer the rights for withdrawing the deposit to said party when the message is finalized. However providing such an early withdrawal at a fee exposes the provider to the risk of the rollup being challenged and proven invalid. This risk should be lower than average as it depends on submitting an invalid rollup which requires breaking SGX. 

## Security

The protocol only listens to the events on the contract address that is bound to the management contract. Messages are bound to L1 blocks and in the event of reorganization the state will be recalculated. L2 messages are bound to L1 block hashes too, resulting in a system where rollups that update the state can only be applied to the correct fork. If there is a mismatch, the state update will be rejected by the management contract.

The protocol controls the keys for the L2 contract. They are hidden within SGX and there are no codepaths that expose them externally. To create a transaction that publishes messages on the L2, one would need to break SGX in order to extract the private key controlling the `MessageBus` L2 contract.
However if one were to produce such a transaction, the code path would reject it automatically as it scans every external transaction. The protocol rejects external transactions that are signed by the hidden signer. This means there is no way to get such a transaction to be processed and even if local SGX is bypassed, validators will immediately pick it up.


## Interface

Those are the current interfaces required to interact with the `MessageBus` contracts. If you encounter issues or have questions, please reach out on our discord.

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