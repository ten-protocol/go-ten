// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import "./IMessageBus.sol";
import "./Structs.sol";

contract MessageBus is IMessageBus {

    function messageFee() virtual internal returns (uint256) { return 0; }

    event LogMessagePublished(address sender, uint64 sequence, uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel);

    mapping (bytes32 => uint256) messageFinalityTimestamps;
    mapping ( address => mapping ( uint32 => Structs.CrossChainMessage[] ) ) messages; 


    function publishMessage(
        uint32 nonce,
        uint32 topic,
        bytes calldata payload, 
        uint8 consistencyLevel
    ) override external returns (uint64 sequence) {
        //require(msg.value >= messageFee());
        
        sequence = 0;
        emit LogMessagePublished(msg.sender, sequence, nonce, topic, payload, consistencyLevel);
        sequence = 1;
        return sequence;
    }

    function verifyMessageFinalized(Structs.CrossChainMessage calldata crossChainMessage) external view override returns (bool) 
    {
        bytes32 msgHash = keccak256(abi.encode(crossChainMessage));
        return messageFinalityTimestamps[msgHash] >= block.number;
    }

    function getMessageTimeOfFinality(Structs.CrossChainMessage calldata crossChainMessage) external view override returns (uint256) {
        bytes32 msgHash = keccak256(abi.encode(crossChainMessage));
        uint256 timeOfFinality = messageFinalityTimestamps[msgHash];

        require(timeOfFinality > 0, "This message was never submitted.");
        return timeOfFinality;
    }

    function submitOutOfNetworkMessage(Structs.CrossChainMessage calldata crossChainMessage, uint256 finalAfterTimestamp) external override 
    {
        uint256 finalAtHeight = block.number + finalAfterTimestamp;
        bytes32 msgHash = keccak256(abi.encode(crossChainMessage));
        
        require(messageFinalityTimestamps[msgHash] == 0, "Message submitted more than once!");

        messageFinalityTimestamps[msgHash] = finalAtHeight;

        messages[crossChainMessage.sender][crossChainMessage.topic].push(crossChainMessage);
    }

    fallback() external payable {revert("unsupported");}

    receive() external payable {revert("the Wormhole contract does not accept assets");}
}