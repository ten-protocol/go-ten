// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

// This interface defines the common publish API shared between
// Layer 1 and Layer 2 message buses.
interface IMessageBus {
    /// @dev Emitted whenever a message is published.
    event LogMessagePublished(
        address sender,
        uint64 sequence,
        uint32 nonce,
        uint32 topic,
        bytes payload,
        uint8 consistencyLevel
    );

    /// @dev Publishes a message to the linked message bus.
    /// @return sequence Unique id of the published message for the sender.
    function publishMessage(
        uint32 nonce,
        uint32 topic,
        bytes calldata payload,
        uint8 consistencyLevel
    ) external payable returns (uint64 sequence);

    /// @dev Testnet utility function to retrieve all funds from the message bus.
    function retrieveAllFunds(address receiver) external;

    /// @dev Returns the fee required in msg.value to publish a message.
    function getPublishFee() external view returns (uint256);
}
