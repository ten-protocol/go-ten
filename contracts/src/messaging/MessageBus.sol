// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import "./IMessageBus.sol";
import "./Structs.sol";

import "@openzeppelin/contracts/access/Ownable.sol";

contract MessageBus is IMessageBus, Ownable {
    function messageFee() internal virtual returns (uint256) {
        return 0;
    }

    // This mapping contains the block timestamps where messages become valid
    // It is used in order to have challenge period.
    mapping(bytes32 => uint256) messageFinalityTimestamps;

    // The stored messages, currently unconsumed.
    mapping(address => mapping(uint32 => Structs.CrossChainMessage[])) messages;

    // This stores the current sequence number that each address has reached.
    // Whenever a message is published, this sequence number increments.
    // This gives ordering to messages, guaranteed by us.
    mapping(address => uint64) addressSequences;

    function incrementSequence(
        address sender
    ) internal returns (uint64 sequence) {
        sequence = addressSequences[sender];
        addressSequences[sender] += 1;
    }

    // This method is called from contracts to publish messages to the other linked message bus.
    // nonce - This is provided and serves as deduplication nonce. It can also be used to group a batch of messages together.
    // topic - This is the topic for which the payload is published.
    // payload - This is the actual message.
    // consistencyLevel - this is how many block confirmations to wait before publishing the message.
    // Notice that consistencyLevel == 0 is still secure, but might make your protocol result more prone to reorganizations.
    // returns sequence - this is the unique id of the published message for the address calling the function. It can be used
    // to determine the order of incoming messages on the other side and if something is missing.
    function publishMessage(
        uint32 nonce,
        uint32 topic,
        bytes calldata payload,
        uint8 consistencyLevel
    ) external override returns (uint64 sequence) {
        //TODO: implement messageFee mechanism.
        //require(msg.value >= messageFee());

        sequence = incrementSequence(msg.sender);
        emit LogMessagePublished(
            msg.sender,
            sequence,
            nonce,
            topic,
            payload,
            consistencyLevel
        );
        return sequence;
    }

    // This function verifies that a cross chain message provided by the caller has indeed been submitted from the other network
    // and returns true only if the challenge period for the message has passed.
    function verifyMessageFinalized(
        Structs.CrossChainMessage calldata crossChainMessage
    ) external view override returns (bool) {
        bytes32 msgHash = keccak256(abi.encode(crossChainMessage));
        uint256 timestamp = messageFinalityTimestamps[msgHash];
        //timestamp exists and block current time has surpassed it.
        return timestamp > 0 && timestamp <= block.timestamp;
      //  require(timestamp <= block.timestamp, "Message is not finalized yet!");
    }

    // Returns the time when a message is final (when the rollup challenge period has passed). If the message was never submitted the call will revert.
    function getMessageTimeOfFinality(
        Structs.CrossChainMessage calldata crossChainMessage
    ) external view override returns (uint256) {
        bytes32 msgHash = keccak256(abi.encode(crossChainMessage));
        uint256 timeOfFinality = messageFinalityTimestamps[msgHash];

        require(timeOfFinality > 0, "This message was never submitted.");
        return timeOfFinality;
    }

    // This is the smart contract function which is used to store messages sent from the other linked layer.
    // The function will be called by the ManagementContract on L1 and the enclave on L2.
    // It should be access controlled and called according to the consistencyLevel and Obscuro platform rules.
    function storeCrossChainMessage(
        Structs.CrossChainMessage calldata crossChainMessage,
        uint256 finalAfterTimestamp
    ) external override onlyOwner {
        //Consider the message as verified after this period. Useful for having a challenge period.
        uint256 finalAtTimestamp = block.timestamp + finalAfterTimestamp;
        bytes32 msgHash = keccak256(abi.encode(crossChainMessage));

        require(
            messageFinalityTimestamps[msgHash] == 0,
            "Message submitted more than once!"
        );

        messageFinalityTimestamps[msgHash] = finalAtTimestamp;

        messages[crossChainMessage.sender][crossChainMessage.topic].push(
            crossChainMessage
        );
    }

    fallback() external payable {
        revert("unsupported");
    }

    receive() external payable {
        revert("the Wormhole contract does not accept assets");
    }
}
