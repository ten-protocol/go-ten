// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import "../common/Structs.sol";
import "./IMerkleTreeMessageBus.sol";
import "./MessageBus.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";

// This contract implements the IMerkleTreeMessageBus interface
contract MerkleTreeMessageBus is IMerkleTreeMessageBus, MessageBus, AccessControlUpgradeable {

    bytes32 public constant STATE_ROOT_MANAGER_ROLE = keccak256("STATE_ROOT_MANAGER_ROLE");
    bytes32 public constant WITHDRAWAL_MANAGER_ROLE = keccak256("WITHDRAWAL_MANAGER_ROLE");

    // When a xchain messages root becomes valid represented as a timestamp in seconds to be compared against block timestamp
    mapping(bytes32 => uint256) rootValidAfter;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() MessageBus() {
        // Constructor intentionally left empty
    }

    /**
     * @dev Initializes the contract with provided owner
     * @param initialOwner Address that will be granted the DEFAULT_ADMIN_ROLE and STATE_ROOT_MANAGER_ROLE
     * @param withdrawalManager Address that will be granted the WITHDRAWAL_MANAGER_ROLE
     */
    function initialize(address initialOwner, address withdrawalManager) public override(IMerkleTreeMessageBus, MessageBus) initializer {
        // Initialize parent contracts
        //super.initialize(initialOwner, address(0));
        __Ownable_init(initialOwner);
        __AccessControl_init();
        
        // Set up roles
        _grantRole(DEFAULT_ADMIN_ROLE, initialOwner);
        _grantRole(STATE_ROOT_MANAGER_ROLE, initialOwner);
        _grantRole(WITHDRAWAL_MANAGER_ROLE, withdrawalManager);
    }

    function addStateRootManager(address manager) external onlyRole(DEFAULT_ADMIN_ROLE) {
        grantRole(STATE_ROOT_MANAGER_ROLE, manager);
    }

    function removeStateRootManager(address manager) external onlyRole(DEFAULT_ADMIN_ROLE) {
        revokeRole(STATE_ROOT_MANAGER_ROLE, manager);
    }

    function addStateRoot(bytes32 stateRoot, uint256 activationTime) external onlyRole(STATE_ROOT_MANAGER_ROLE) {
        require(rootValidAfter[stateRoot] == 0, "Root already added to the message bus");
        rootValidAfter[stateRoot] = activationTime;
    }

    function disableStateRoot(bytes32 stateRoot) external onlyRole(STATE_ROOT_MANAGER_ROLE) {
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