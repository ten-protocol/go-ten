// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

interface IBridge {
    struct TransferMessage {
        address asset;
        uint256 amount;
        address target;
    }

    enum Topics {
        TRANSFER,
        MANAGEMENT
    }

    function sendNative(address target) external payable;
    function sendAssets(address asset, uint256 amount, address receiver) external;
    function receiveAssets(address asset, uint256 amount, address receiver) external;
}