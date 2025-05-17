// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

/**
 * @title Logger
 * @dev Library for logging messages
 */
library Logger {
    /**
     * @dev Emitted when a message is logged
     * @param message The message to log
     */
    event LogMessage(string message);

    /**
     * @dev Emits a message
     * @param message The message to emit
     */
    function emitLog(string calldata message) external {
        emit LogMessage(message);
    }
}