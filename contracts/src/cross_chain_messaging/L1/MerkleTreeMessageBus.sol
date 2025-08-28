// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import "../../common/Structs.sol";
import "./IMerkleTreeMessageBus.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "../../common/UnrenouncableOwnable2Step.sol";
import "../../common/PausableWithRoles.sol";
import "../../system/interfaces/IFees.sol";
import {MessageBus} from "../common/MessageBus.sol";
/**
 * @title MerkleTreeMessageBus
 * @dev Implementation of IMerkleTreeMessageBus that uses Merkle trees for cross-chain message verification
 * This contract manages state roots and verifies message inclusion through Merkle proofs.
 * It implements a role-based access control system for state root and withdrawal management.
 */
contract MerkleTreeMessageBus is IMerkleTreeMessageBus, MessageBus {

    /**
     * @dev Role identifier for accounts that can manage state roots
     */
    bytes32 public constant STATE_ROOT_MANAGER_ROLE = keccak256("STATE_ROOT_MANAGER_ROLE");

    /**
     * @dev Role identifier for accounts that can manage withdrawals
     */
    bytes32 public constant WITHDRAWAL_MANAGER_ROLE = keccak256("WITHDRAWAL_MANAGER_ROLE");

    /**
     * @dev Mapping of state roots to their activation timestamps
     * A value of 0 indicates either the root doesn't exist or has been disabled
     */
    mapping(bytes32 stateRoot => uint256 activationTime) public rootValidAfter;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() MessageBus() {
        // Constructor intentionally left empty
    }

    /**
     * @dev Initializes the contract with provided owner
     * @param initialOwner Address that will be granted the DEFAULT_ADMIN_ROLE and STATE_ROOT_MANAGER_ROLE
     */
    function initialize(address initialOwner, address withdrawalManager, address _fees) public override(IMerkleTreeMessageBus, MessageBus) initializer {
        require(initialOwner != address(0), "Initial owner cannot be 0x0");
        require(withdrawalManager != address(0), "Withdrawal manager cannot be 0x0");
        require(_fees != address(0), "Fees address cannot be 0x0");

        // Initialize parent contracts
        __UnrenouncableOwnable2Step_init(initialOwner);
        __PausableWithRoles_init(initialOwner);

        // Set up additional roles
        _grantRole(DEFAULT_ADMIN_ROLE, initialOwner);
        _grantRole(STATE_ROOT_MANAGER_ROLE, initialOwner);
        _grantRole(WITHDRAWAL_MANAGER_ROLE, withdrawalManager);
        fees = IFees(_fees);
    }

    /**
     * @dev Grants STATE_ROOT_MANAGER_ROLE to a new address
     * @param manager Address to be granted the role
     */
    function addStateRootManager(address manager) external onlyRole(DEFAULT_ADMIN_ROLE) whenNotPaused {
        require(manager != address(0), "Manager cannot be 0x0");
        grantRole(STATE_ROOT_MANAGER_ROLE, manager);
    }

    /**
     * @dev Revokes STATE_ROOT_MANAGER_ROLE from an address
     * @param manager Address to have the role revoked
     */
    function removeStateRootManager(address manager) external onlyRole(DEFAULT_ADMIN_ROLE) whenNotPaused {
        require(manager != address(0), "Manager cannot be 0x0");
        revokeRole(STATE_ROOT_MANAGER_ROLE, manager);
    }

    /**
     * @dev Grants WITHDRAWAL_MANAGER_ROLE to a new address
     * @param withdrawalManager Address to be granted the role
     */
    function addWithdrawalManager(address withdrawalManager) external onlyRole(DEFAULT_ADMIN_ROLE) whenNotPaused {
        _grantRole(WITHDRAWAL_MANAGER_ROLE, withdrawalManager);
    }

    /**
     * @dev Revokes WITHDRAWAL_MANAGER_ROLE from an address
     * @param withdrawalManager Address to have the role revoked
     */
    function removeWithdrawalManager(address withdrawalManager) external onlyRole(DEFAULT_ADMIN_ROLE) whenNotPaused {
        revokeRole(WITHDRAWAL_MANAGER_ROLE, withdrawalManager);
    }
    /**
     * @dev Adds a new state root with its activation time
     * @param stateRoot The state root to be added
     * @param activationTime Timestamp after which the root becomes valid
     * @notice Root must not already exist in the message bus
     */
    function addStateRoot(bytes32 stateRoot, uint256 activationTime) external onlyRole(STATE_ROOT_MANAGER_ROLE) whenNotPaused{
        require(rootValidAfter[stateRoot] == 0, "Root already added to the message bus");
        require(activationTime >= block.timestamp, "Activation time must be in the future");
        rootValidAfter[stateRoot] = activationTime;
    }

    /**
     * @dev Disables an existing state root
     * @param stateRoot The state root to be disabled
     * @notice Root must exist in the message bus
     */
    function disableStateRoot(bytes32 stateRoot) external onlyRole(STATE_ROOT_MANAGER_ROLE) {
        require(rootValidAfter[stateRoot] != 0, "State root does not exist.");
        rootValidAfter[stateRoot] = 0;
    }

    /**
     * @dev Verifies inclusion of a cross-chain message in a state root using a Merkle proof
     * @param message The cross-chain message to verify
     * @param proof Merkle proof demonstrating inclusion
     * @param root State root to verify against
     * @notice Root must be published and its activation time must have passed
     */
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

    /**
     * @dev Verifies inclusion of a value transfer message in a state root using a Merkle proof
     * @param message The value transfer message to verify
     * @param proof Merkle proof demonstrating inclusion
     * @param root State root to verify against
     * @notice Root must be published and its activation time must have passed
     */
    function verifyValueTransferInclusion(Structs.ValueTransferMessage calldata message, bytes32[] calldata proof, bytes32 root) external view {
        require(rootValidAfter[root] != 0, "Root is not published on this message bus.");
        require(block.timestamp >= rootValidAfter[root], "Root is not considered final yet.");

        bytes32 leaf = keccak256(abi.encode("v", keccak256(abi.encode(message))));

        require(MerkleProof.verifyCalldata(proof, root, keccak256(abi.encodePacked(leaf))), "Invalid inclusion proof for value transfer message.");
    }
}