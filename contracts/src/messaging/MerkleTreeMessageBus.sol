// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";

import "./Structs.sol";
import "./IMerkleTreeMessageBus.sol";
import "./MessageBus.sol";

contract MerkleTreeMessageBus is IMerkleTreeMessageBus, MessageBus {
    constructor() MessageBus() {}

    mapping(bytes32 => uint256) rootValidAfter; //When a xchain messages root becomes valid represented as a timestamp in seconds to be compared against block timestamp

    function addStateRoot(bytes32 stateRoot, uint256 activationTime) external onlyOwner {
        require(rootValidAfter[stateRoot] == 0, "Root already added to the message bus");
        rootValidAfter[stateRoot] = activationTime;
    }

    function blockStateRoot(bytes32 stateRoot) external onlyOwner {
        require(rootValidAfter[stateRoot] != 0, "State root does not exist.");
        rootValidAfter[stateRoot] = 0;
    }

    function verifyMessageInclusion(Structs.CrossChainMessage calldata message, bytes32[] calldata proof, bytes32 root) external view {
        require(rootValidAfter[root] != 0, "Root is not published on this message bus.");
        require(block.timestamp >= rootValidAfter[root], "Root is not considered final yet.");

        bytes32 leaf = keccak256(abi.encode("message", keccak256(abi.encode(message))));

        require(MerkleProof.verifyCalldata(proof, root, leaf), "Invalid inclusion proof for cross chain message.");
    }

    function verifyValueTransferInclusion(Structs.ValueTransferMessage calldata message, bytes32[] calldata proof, bytes32 root) external view {
        require(rootValidAfter[root] != 0, "Root is not published on this message bus.");
        require(block.timestamp >= rootValidAfter[root], "Root is not considered final yet.");

        bytes32 leaf = keccak256(abi.encode("value", keccak256(abi.encode(message))));

        require(MerkleProof.verifyCalldata(proof, root, leaf), "Invalid inclusion proof for value transfer message.");
    }
}