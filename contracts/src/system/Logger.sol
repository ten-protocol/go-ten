// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library Logger {
    // Define the event with a string parameter
    event LogMessage(string message);

    // Function to emit the LogMessage event
    function emitLog(string calldata message) external {
        emit LogMessage(message);
    }
}