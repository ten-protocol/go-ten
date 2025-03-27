// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "../../common/Structs.sol";

/**
 * @title IMessageBus
 * @dev Interface that allows dApps and accounts to send and verify received messages
 * between layer 1 and layer 2.
 */
interface IMessageBus {

    /**
     * @dev Emitted when a message is published
     * The enclave listens for it on the deployed message bus addresses
     * @param sender The address publishing the message
     * @param sequence The unique message ID for the sender
     * @param nonce Deduplication nonce, can group messages together
     * @param topic The topic for which the payload is published
     * @param payload The actual message content
     * @param consistencyLevel Number of block confirmations to wait
     */
    event LogMessagePublished(
        address sender, 
        uint64 sequence, 
        uint32 nonce, 
        uint32 topic, 
        bytes payload, 
        uint8 consistencyLevel
    );

    /**
     * @dev Emitted when value is transferred between layers
     * @param sender The address sending the value
     * @param receiver The address receiving the value
     * @param amount The amount being transferred
     * @param sequence The unique transfer ID
     */
    event ValueTransfer(
        address indexed sender,
        address indexed receiver,
        uint256 amount,
        uint64 sequence
    );

    /**
     * @dev Emitted when native tokens are deposited
     * @param receiver The address receiving the deposit
     * @param amount The amount being deposited
     */
    event NativeDeposit(
        address indexed receiver,
        uint256 amount
    );

    /**
     * @dev Publishes messages to the other linked message bus
     * @param nonce Deduplication nonce, can group messages together
     * @param topic The topic for which the payload is published
     * @param payload The actual message content
     * @param consistencyLevel Block confirmations to wait. Level 0 is secure but more prone to reorganizations
     * @return sequence Unique ID of the published message for the calling address
     */
    function publishMessage(
        uint32 nonce,
        uint32 topic,
        bytes calldata payload, 
        uint8 consistencyLevel
    ) external payable returns (uint64 sequence);

    /**
     * @dev Sends value to L2
     * @param receiver The address to receive the value on L2
     * @param amount The amount to send
     */
    function sendValueToL2(
        address receiver,
        uint256 amount
    ) external payable;

    /**
     * @dev Receives value from L2
     * @param receiver The address to receive the value
     * @param amount The amount being received
     */
    function receiveValueFromL2(
        address receiver,
        uint256 amount
    ) external;

    /**
     * @dev Verifies a cross chain message has been submitted and its challenge period has passed
     * @param crossChainMessage The message to verify
     * @return bool True if the message is finalized
     */
    function verifyMessageFinalized(Structs.CrossChainMessage calldata crossChainMessage) external view returns (bool);
    
    /**
     * @dev Gets the timestamp when a message becomes final (after challenge period)
     * @param crossChainMessage The message to check
     * @return uint256 The timestamp of finality
     */
    function getMessageTimeOfFinality(Structs.CrossChainMessage calldata crossChainMessage) external view returns (uint256);

    /**
     * @dev Stores messages sent from the other linked layer
     * Called by ManagementContract on L1 and the enclave on L2
     * Must be access controlled according to consistencyLevel and platform rules
     * @param crossChainMessage The message to store
     * @param finalAfterTimestamp The timestamp after which the message is considered final
     */
    function storeCrossChainMessage(Structs.CrossChainMessage calldata crossChainMessage, uint256 finalAfterTimestamp) external;

    /**
     * @dev Testnet function to retrieve all funds from the message bus
     * @param receiver The address to receive the funds
     */
    function retrieveAllFunds(address receiver) external;

    /**
     * @dev Gets the required fee to publish a value transfer
     * @return uint256 The fee amount that must be paid in msg.value
     */
    function getPublishFee() external view returns (uint256);
}
