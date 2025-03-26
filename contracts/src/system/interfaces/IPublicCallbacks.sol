// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

/**
 * @title IPublicCallbacks
 * @dev Interface for the   PublicCallbacks contract
 */
interface IPublicCallbacks {
    /**
     * @dev Register a callback
     * @param callback The callback to register
     * @return The callback ID
     */
    function register(bytes calldata callback) external payable returns (uint256);
    /**
     * @dev Reattempt a callback
     * @param callbackId The callback ID
     */
    function reattemptCallback(uint256 callbackId) external;
}