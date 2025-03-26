// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

interface IPublicCallbacks {
    function register(bytes calldata callback) external payable returns (uint256);
    function reattemptCallback(uint256 callbackId) external;
}