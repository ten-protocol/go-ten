// SPDX-License-Identifier: MIT
pragma solidity >=0.7.0 <0.9.0;

/**
 * @title Storage
 * @dev Library for managing storage slots in upgradeable contracts
 * This pattern prevents storage collisions between different versions of a contract
 * by using explicit storage slots for each variable
 */
library Storage {
    /**
     * @dev Reads an address from a specific storage slot
     * @param _slot The storage slot to read from, typically calculated as keccak256(variable_name) - 1
     * @return addr_ The address stored in the slot
     * @notice Uses assembly for direct storage access which is more gas efficient
     */
    function getAddress(bytes32 _slot) internal view returns (address addr_) {
        assembly {
            addr_ := sload(_slot)
        }
    }

    /**
     * @dev Writes an address to a specific storage slot
     * @param _slot The storage slot to write to
     * @param _address The address to store
     * @notice Storage slots should be unique and deterministic to avoid collisions
     * @notice Formula for slot calculation: bytes32(uint256(keccak256("variable_name")) - 1)
     */
    function setAddress(bytes32 _slot, address _address) internal {
        assembly {
            sstore(_slot, _address)
        }
    }
}