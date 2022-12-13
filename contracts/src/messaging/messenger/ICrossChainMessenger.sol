// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "../IMessageBus.sol";
import "../Structs.sol";

interface ICrossChainMessenger {
    struct CrossChainCall {
        address target;
        bytes data;
        uint256 gas;
    }

    function messageBus() external view returns (address);
    function crossChainSender() external view returns (address);

    function relayMessage(Structs.CrossChainMessage calldata message) external;
}