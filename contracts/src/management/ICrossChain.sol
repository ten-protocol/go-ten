// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "../common/Structs.sol";

/**
 * @title ICrossChain
 * @dev Interface for managing cross-chain value transfers and withdrawals
 */
interface ICrossChain {

    event WithdrawalsPaused(bool paused);

    /**
     * @dev Extracts native tokens from L2 to L1 using a verified message
     * @param msg The value transfer message containing amount and recipient details
     * @param proof Merkle proof verifying the message inclusion
     * @param root The Merkle root against which to verify the proof
     */
    function extractNativeValue(
        Structs.ValueTransferMessage calldata msg,
        bytes32[] calldata proof,
        bytes32 root
    ) external;

    /**
     * @dev Enables or disables the withdrawal functionality
     * @param pause True to pause withdrawals, false to enable
     */
    function pauseWithdrawals(bool pause) external;

    /**
     * @dev Checks if a withdrawal has already been processed to prevent double-spending
     * @param messageHash Hash of the withdrawal message
     * @return bool True if withdrawal was already processed
     */
    function isWithdrawalSpent(bytes32 messageHash) external view returns (bool);

    /**
     * @dev Verifies if a bundle of cross-chain messages is available
     * @param crossChainHashes Array of cross-chain message hashes to verify
     * @return bool True if the bundle is available
     */
    function isBundleAvailable(bytes[] memory crossChainHashes) external view returns (bool);
}