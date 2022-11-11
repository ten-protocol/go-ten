// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "./Structs.sol";

interface IMessageBus {
    function publishMessage(
        uint32 nonce,
        bytes memory topic,
        bytes memory payload, 
        uint8 consistencyLevel
    ) external payable returns (uint64 sequence);

    function verifyMessageFinalized(Structs.CrossChainMessage calldata crossChainMessage) external view returns (bool);
    
    function getMessageTimeOfFinality(Structs.CrossChainMessage calldata crossChainMessage) external view returns (uint256);

    function submitOutOfNetworkMessage(Structs.CrossChainMessage calldata crossChainMessage, uint256 finalAfterTimestamp) external;

   /* function queryMessages(
        address      sender,
        bytes memory topic,
        uint256      fromIndex,
        uint256      toIndex
    ) external returns (bytes [] memory); */
}
