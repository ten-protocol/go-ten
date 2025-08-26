// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "../../common/Structs.sol";

/**
 * @title ICrossChain
 * @dev Interface for managing cross-chain value transfers and withdrawals
 */
interface ICrossChain {
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