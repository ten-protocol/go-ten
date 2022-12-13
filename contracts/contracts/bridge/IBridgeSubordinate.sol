// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;


interface IBridgeSubordinate {
    function createWrappedToken(address crossChainAddress, string calldata name, string calldata symbol) external;
}