// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

interface Structs {
    struct CrossChainMessage {
        address sender; // The contract/address which called the publishMessage on the message bus.
        uint64  sequence; // The sequential index of the message for the sending address.
        uint32  nonce; // Provided by the smart contract. Can be used to create message groups for multi step protocols.
        uint32  topic; // Can be used to separate messages and provide basic versioning.
        bytes   payload; // The actual encoded message.
        uint8   consistencyLevel; //
    }

    struct ValueTransferMessage {
        address sender;
        address receiver;
        uint256 amount;
        uint64  sequence;
    }
}