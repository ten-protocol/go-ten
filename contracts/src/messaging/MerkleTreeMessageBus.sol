// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import "../common/Structs.sol";
import "./IMerkleTreeMessageBus.sol";
import "./MessageBus.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";

contract MerkleTreeMessageBus is IMerkleTreeMessageBus, MessageBus {

    address public admin;

    // The list of addresses that are allowed to call the addStateRoot functioned. The owner of this contract
    // manually adds the rollup contract to this mapping once the contracts have been deployed.
    mapping(address => bool) public stateRootManagers;

    constructor(address _admin) MessageBus() {
        admin = _admin;
        stateRootManagers[_admin] = true;
    }

    modifier onlyAdmin() {
        require(msg.sender == admin, "Only admin can call this function");
        _;
    }

    modifier onlyStateRootManager() {
        require(stateRootManagers[msg.sender], "Only state root managers can call this function");
        _;
    }

    function addStateRootManager(address manager) external onlyAdmin {
        stateRootManagers[manager] = true;
    }

    function removeStateRootManager(address manager) external onlyAdmin {
        stateRootManagers[manager] = false;
    }

    function transferAdmin(address newAdmin) external onlyAdmin {
        admin = newAdmin;
    }

    mapping(bytes32 => uint256) rootValidAfter; //When a xchain messages root becomes valid represented as a timestamp in seconds to be compared against block timestamp

    function addStateRoot(bytes32 stateRoot, uint256 activationTime) external onlyStateRootManager {
        require(rootValidAfter[stateRoot] == 0, "Root already added to the message bus");
        rootValidAfter[stateRoot] = activationTime;
    }

    function disableStateRoot(bytes32 stateRoot) external onlyStateRootManager {
        require(rootValidAfter[stateRoot] != 0, "State root does not exist.");
        rootValidAfter[stateRoot] = 0;
    }

    function verifyMessageInclusion(Structs.CrossChainMessage calldata message, bytes32[] calldata proof, bytes32 root) external view {
        require(rootValidAfter[root] != 0, "Root is not published on this message bus.");
        require(block.timestamp >= rootValidAfter[root], "Root is not considered final yet.");

        bytes32 messageHash = keccak256(abi.encode(
            message.sender,
            message.sequence,
            message.nonce,
            message.topic,
            message.payload,
            message.consistencyLevel
        ));
        bytes32 leaf = keccak256(abi.encode("m", messageHash));

        require(MerkleProof.verifyCalldata(proof, root, keccak256(abi.encodePacked(leaf))), "Invalid inclusion proof for cross chain message.");
    }

    function verifyValueTransferInclusion(Structs.ValueTransferMessage calldata message, bytes32[] calldata proof, bytes32 root) external view {
        require(rootValidAfter[root] != 0, "Root is not published on this message bus.");
        require(block.timestamp >= rootValidAfter[root], "Root is not considered final yet.");

        bytes32 leaf = keccak256(abi.encode("v", keccak256(abi.encode(message))));

        require(MerkleProof.verifyCalldata(proof, root, keccak256(abi.encodePacked(leaf))), "Invalid inclusion proof for value transfer message.");
    }
}