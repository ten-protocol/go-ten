// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "../../common/Structs.sol";

/**
 * @title IDataAvailabilityRegistry
 * @dev Interface for managing L2 rollup submissions
 */
interface IDataAvailabilityRegistry {

    // MetaRollup: Represents metadata for a rollup process, including sequence numbers, hashes, and signatures.
    struct MetaRollup {
        bytes32 Hash;
        uint256 FirstSequenceNumber;
        uint256 LastSequenceNumber;

        bytes32 BlockBindingHash;
        uint256 BlockBindingNumber;
        bytes32 crossChainRoot;
        bytes32 LastBatchHash;
        bytes Signature;
    }

    /**
     * @dev Emitted when a new rollup is added
     * @param rollupHash The hash of the rollup
     * @param signature The signature of the rollup
     */
    event RollupAdded(bytes32 rollupHash, bytes signature);
    
    /**
     * @dev Adds a new rollup batch to the chain and processes its cross-chain messages
     * @param rollup The MetaRollup containing:
     *        - LastSequenceNumber: Latest sequence number for ordering
     *        - LastBatchHash: Hash of the previous batch
     *        - BlockBindingHash: Hash of the block this rollup is bound to
     *        - BlockBindingNumber: Number of the block this rollup is bound to
     *        - crossChainRoot: Merkle root of cross-chain messages (if any)
     *        - Signature: Sequencer enclave signature of the rollup data
     * @notice Requires the sender to be an attested sequencer enclave
     * @notice If crossChainRoot is present, it will be added to MessageBus with a delay of challengePeriod
     */
    function addRollup(MetaRollup calldata rollup) external;

    /**
     * @dev Retrieves a previously submitted rollup batch by its hash
     * @param rollupHash The hash of the rollup to retrieve
     * @return bool True if the rollup exists, false otherwise
     * @return Structs.MetaRollup The rollup data if it exists, empty struct if not
     */
    function getRollupByHash(bytes32 rollupHash) external view returns (bool, MetaRollup memory);

    /**
     * @dev Returns the current rollup challenge period duration
     * @return uint256 Duration of the challenge period
     */
    function getChallengePeriod() external view returns (uint256);

    /**
     * @dev Sets the duration of the rollup challenge period
     * @param delay New duration for the challenge period
     */
    function setChallengePeriod(uint256 delay) external;

    /**
     * @dev Pauses the contract in case of emergency
     * @notice Only callable by the owner
     */
    function pause() external;

    /**
     * @dev Unpauses the contract
     * @notice Only callable by the owner
     */
    function unpause() external;
}