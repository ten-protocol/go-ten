// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import "../../common/Structs.sol";
import "../../system/contracts/Fees.sol";

import "../../system/interfaces/IFees.sol";
import "./IMessageBus.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

/**
 * @title MessageBus
 * @dev Implementation of the IMessageBus interface for cross-layer message handling.
 * Manages message publishing, verification, and value transfers between L1 and L2.
 */
contract MessageBus is IMessageBus, Initializable, OwnableUpgradeable {

    constructor() {
        _transferOwnership(msg.sender);
    }

     /**
     * @dev Initializes the contract with an owner and fees contract
     * @param caller The address to set as the owner
     * @param feesAddress The address of the fees contract
     */
    function initialize(address caller, address feesAddress) public virtual initializer {
        __Ownable_init(caller);
        fees = IFees(feesAddress);
    }

    
    /**
     * @dev Modifier to restrict access to owner or self
     * Since this contract exists on L2, when messages are added from L1,
     * we can have the from address be the same as self.
     * This ensures no EOA collision can occur and no key needs to be stored
     * on L2 or shared with validators.
     */
    modifier ownerOrSelf() {
        address maskedSelf = address(uint160(address(this)) - 1);
        require(msg.sender == owner() || msg.sender == maskedSelf, "Not owner or self");
        _;
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

    IFees fees;

    /**
     * @dev Increments and returns the sequence number for a sender
     * @param sender The address to increment the sequence for
     * @return sequence The previous sequence number
     */
    function incrementSequence(
        address sender
    ) internal returns (uint64 sequence) {
        sequence = addressSequences[sender];
        addressSequences[sender] += 1;
    }

    /**
     * @dev Sends value to L2
     * @param receiver The address to receive the value on L2
     * @param amount The amount to send
     */
    function sendValueToL2(
        address receiver,
        uint256 amount
    ) external payable {
        require(msg.value > 0 && msg.value == amount, "Attempting to send value without providing Ether");
        
        uint256 amountToBridge = msg.value;
        if (address(fees) != address(0)) {
            uint256 fee = getPublishFee();
            require(msg.value >= fee, "Insufficient funds to send value");
            amountToBridge = msg.value - fee;
            (bool ok, ) = address(fees).call{value: fee}("");
            require(ok, "Failed to send fees to fees contract");
        }

        uint64 sequence = incrementSequence(msg.sender);
        emit ValueTransfer(msg.sender, receiver, amountToBridge, sequence);
    }

    /**
     * @dev Receives value from L2, restricted to owner
     * @param receiver The address to receive the value
     * @param amount The amount being received
     */
    function receiveValueFromL2(
        address receiver,
        uint256 amount
    ) external virtual onlyOwner {
        _receiveValueFromL2Internal(receiver, amount);
    }

    /**
     * @dev Internal function with the core logic for receiving value from L2
     * Used to avoid duplicating code when overriding receiveValueFromL2
     * @param receiver The address to receive the value
     * @param amount The amount to send
     */
    function _receiveValueFromL2Internal(address receiver, uint256 amount) internal {
        require(address(this).balance >= amount, "Insufficient funds to send value");
        (bool ok, ) = receiver.call{value: amount}("");
        require(ok, "failed sending value");
    }

    function getPublishFee() public view returns (uint256) {
        return fees.messageFee();
    }

    /**
     * @dev Publishes a message to the other linked message bus
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
    ) external payable override returns (uint64 sequence) {
        if (address(fees) != address(0)) { // No fee required for L1 to L2 messages.
            uint256 fee = getPublishFee();
            require(msg.value >= fee, "Insufficient funds to publish message");
            (bool ok, ) = address(fees).call{value: fee}("");
            require(ok, "Failed to send fees to fees contract");
        }

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

    /**
     * @dev Verifies  that a cross chain message provided by the caller has indeed been submitted from the other network
     *  and returns true only if the challenge period for the message has passed.
     * @param crossChainMessage The message to verify
     * @return bool True if the message's challenge period has passed
     */
    function verifyMessageFinalized(
        Structs.CrossChainMessage calldata crossChainMessage
    ) external view override returns (bool) {
        bytes32 msgHash = keccak256(abi.encode(crossChainMessage));
        
        uint256 timestamp = messageFinalityTimestamps[msgHash];
        //timestamp exists and block current time has surpassed it.
        return timestamp > 0 && timestamp <= block.timestamp;

    }

    /**
     * @dev Gets the finality timestamp for a message (after the rollup challenge period has passed). If the message was never submitted the call will revert.
     * @param crossChainMessage The message to check
     * @return uint256 The timestamp when the message becomes final
     */
    function getMessageTimeOfFinality(
        Structs.CrossChainMessage calldata crossChainMessage
    ) external view override returns (uint256) {
        bytes32 msgHash = keccak256(abi.encode(crossChainMessage));
        uint256 timeOfFinality = messageFinalityTimestamps[msgHash];

        require(timeOfFinality > 0, "This message was never submitted.");
        return timeOfFinality;
    }

    /**
     * @dev Stores messages from the other linked layer. It should be access controlled and called according to the consistencyLevel and TEN platform rules.
     * @param crossChainMessage The message to store
     * @param finalAfterTimestamp Time to wait before considering message final
     */
    function storeCrossChainMessage(
        Structs.CrossChainMessage calldata crossChainMessage,
        uint256 finalAfterTimestamp
    ) external override ownerOrSelf {
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

    /**
     * @dev Notifies of a deposit event
     * @param receiver The address receiving the deposit
     * @param amount The amount deposited
     */
    function notifyDeposit(
        address receiver,
        uint256 amount
    ) external ownerOrSelf {
        emit NativeDeposit(receiver, amount);
    }

    /**
     * @dev Retrieves all funds from the contract (Testnet only - to be removed before mainnet deployment)
     * @param receiver The address to receive the funds
     */
    function retrieveAllFunds(
        address receiver
    ) external onlyOwner {
        (bool ok, ) = receiver.call{value: address(this).balance}("");
        require(ok, "failed sending value");
    }

    fallback() external payable {
        revert("unsupported");
    }

    receive() external payable {
        this.sendValueToL2{value: msg.value}(msg.sender, msg.value);
    }
}
