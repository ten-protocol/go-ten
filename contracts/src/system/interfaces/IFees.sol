// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

/**
 * @title IFees
 * @dev Interface for the Fees contract
 */     
interface IFees {
    /**
     * @dev Returns the fee for a message
     * @return The fee for a message
     */
    function messageFee() external view returns (uint256);
}
