// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import "../../common/Structs.sol";
import "./IMessageBus.sol";

/// @title IL2MessageBus
/// @dev Extension of IMessageBus with Layer 2 specific functionality
interface IL2MessageBus is IMessageBus {
    /// @dev Verifies that a cross chain message has been submitted from the other network
    ///      and that its challenge period has passed.
    function verifyMessageFinalized(
        Structs.CrossChainMessage calldata crossChainMessage
    ) external view returns (bool);

    /// @dev Returns the timestamp when a message becomes final. Reverts if the message
    ///      was never submitted.
    function getMessageTimeOfFinality(
        Structs.CrossChainMessage calldata crossChainMessage
    ) external view returns (uint256);

    /// @dev Stores a message coming from the other linked layer.
    /// @param crossChainMessage The message to store
    /// @param finalAfterTimestamp Seconds after which the message is considered final
    function storeCrossChainMessage(
        Structs.CrossChainMessage calldata crossChainMessage,
        uint256 finalAfterTimestamp
    ) external;
}
