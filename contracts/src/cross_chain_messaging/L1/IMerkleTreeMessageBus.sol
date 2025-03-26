// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import "../../common/Structs.sol";

/**
 * @title IMerkleTreeMessageBus
 * @dev Interface for a message bus that uses Merkle trees to verify cross-chain messages.
 * Provides functionality for managing state roots and verifying message inclusion
 * through Merkle proofs.
 */
interface IMerkleTreeMessageBus {
    /**
     * @dev Adds a new cross-chain state root to the message bus
     * @param stateRoot The state root to be added
     * @param activationTime The timestamp after which this state root becomes valid
     */
    function addStateRoot(bytes32 stateRoot, uint256 activationTime) external;

    /**
     * @dev Disables a previously added state root from the message bus
     * This can be used in challenge scenarios to invalidate a disputed state root
     * @param stateRoot The state root to be disabled
     */
    function disableStateRoot(bytes32 stateRoot) external;

    /**
     * @dev Verifies that a cross-chain message is included in a specific state root
     * using a Merkle proof
     * @param message The cross-chain message to verify
     * @param proof Array of hashes representing the Merkle proof
     * @param root The state root to verify against (must be previously added)
     */
    function verifyMessageInclusion(Structs.CrossChainMessage calldata message, bytes32[] memory proof, bytes32 root) external view;

    /**
     * @dev Verifies that a value transfer message is included in a specific state root
     * using a Merkle proof
     * @param message The value transfer message to verify
     * @param proof Array of hashes representing the Merkle proof
     * @param root The state root to verify against (must be previously added)
     */
    function verifyValueTransferInclusion(Structs.ValueTransferMessage calldata message, bytes32[] calldata proof, bytes32 root) external view;

    /**
     * @dev Initializes the message bus with an owner and withdrawal manager
     * @param initialOwner Address of the initial owner of the contract
     * @param withdrawalManager Address of the withdrawal manager contract
     */
    function initialize(address initialOwner, address withdrawalManager) external;
}