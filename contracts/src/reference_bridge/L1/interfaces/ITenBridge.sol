// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

// The interface for the layer 1 bridge that drives the whitelist
// and has the functionality to modify it.
interface ITenBridge {
    // This will whitelist a token and generate a cross chain message to the ITokenFactory
    // to create wrapped tokens in case of success.
    function whitelistToken(
        address asset,
        string calldata name,
        string calldata symbol
    ) external;

    // This will pause deposits for this token on the L1 bridge. Withdrawals are still fine.
    function pauseToken(address asset) external;

    // This will unpause deposits for this token on the L1 bridge.
    function unpauseToken(address asset) external;

    function setRemoteBridge(address bridge) external;

    // This will retrieve all funds from the bridge and send to receiver address. Testnet only.
    function retrieveAllFunds(address receiver) external;
}
