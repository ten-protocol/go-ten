// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "./IMessageBus.sol";
import "../../common/Structs.sol";

/**
 * @title ICrossChainMessenger
 * @dev Interface for managing cross-chain message delivery and verification.
 * Provides functionality to relay messages between chains and verify their authenticity
 * through a message bus system.
 */
interface ICrossChainMessenger {
    
    /**
     * @dev Structure representing a cross-chain function call
     * @param target The address of the contract to call on the destination chain
     * @param data The calldata to execute on the target contract
     * @param gas The amount of gas to allocate for the cross-chain execution
     */
    struct CrossChainCall {
        address target;
        bytes data;
        uint256 gas;
    }

    /**
     * @dev Returns the address of the message bus used for message verification
     * @return address The message bus contract address
     */
    function messageBus() external view returns (address);

    /**
     * @dev Returns the original sender of the current cross-chain message being relayed
     * @return address The address of the message sender from the source chain
     */
    function crossChainSender() external view returns (address);

    /**
     * @dev Returns the target address of the current cross-chain message being relayed
     * @return address The address of the target contract
     */
    function crossChainTarget() external view returns (address);

    /**
     * @dev Relays a verified cross-chain message to its target contract
     * Only processes messages that exist in the message bus and haven't been consumed
     * @param message The cross-chain message containing target and execution details
     */
    function relayMessage(Structs.CrossChainMessage calldata message) external;
}
