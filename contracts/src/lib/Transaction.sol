// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title Transaction
 * @dev Library for representing a TEN transaction in the system contracts
 */
library Structs {
    struct Transaction {
        uint8 txType;
        uint256 nonce;
        uint256 gasPrice;
        uint256 gasLimit;
        address to;
        uint256 value;
        bytes data;
        address from;
        bool successful;
        uint64 gasUsed;
    }
}