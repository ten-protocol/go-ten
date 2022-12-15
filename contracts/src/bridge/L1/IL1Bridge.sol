// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;


interface IL1Bridge {
    function whitelistToken(address asset, string calldata name, string calldata symbol) external;

    function removeToken(address asset) external;

    function setRemoteBridge(address bridge) external;
}