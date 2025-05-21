// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import "../../common/Structs.sol";
import "../../system/contracts/Fees.sol";

import "../../system/interfaces/IFees.sol";
import "./BaseMessageBus.sol";
import "./IL2MessageBus.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

/**
 * @title MessageBus
 * @dev Implementation of the IMessageBus interface for cross-layer message handling.
 * Manages message publishing, verification, and value transfers between L1 and L2.
 */
contract MessageBus is BaseMessageBus, IL2MessageBus {

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
    mapping(bytes32 messageHash => uint256 messageFinalityTimestamp) messageFinalityTimestamps;

    // The stored messages, currently unconsumed.
    mapping(address sender => mapping(uint32 topic => Structs.CrossChainMessage[] messages)) messages;

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

    fallback() external {
        revert("unsupported");
    }
}
