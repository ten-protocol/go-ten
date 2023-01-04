// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

// The interface for the layer 1 bridge that drives the whitelist
// and has the functionality to modify it.
interface IObscuroBridgeAdmin {
    // This will whitelist a token and generate a cross chain message to the ITokenFactory
    // to create wrapped tokens in case of succcess.
    function whitelistToken(
        address asset,
        string calldata name,
        string calldata symbol
    ) external;

    // This will delist the token and queue a message for it to be delisted on L2. Notice that the token itself
    // can still be transfered between users, just not across chains.
    function removeToken(address asset) external;

    function setRemoteBridge(address bridge) external;
}
