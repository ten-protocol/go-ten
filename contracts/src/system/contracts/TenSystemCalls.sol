// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/proxy/utils/Initializable.sol";

contract TenSystemCalls is Initializable {

    function initialize() external initializer {
    }

    function getRandomNumber() external view returns (uint256) {
        // We inject randomness in prevrandao as the first 28 bytes followed by the last 4 bytes are the timestamp delta
        // In practice this means the whole can be used for randomness if we go through a hash function.
        return uint256(keccak256(abi.encodePacked(block.prevrandao))); 
    }

    function getTransactionTimestamp() external view returns (uint256) {
        // Extract last 4 bytes from difficulty which contains timestamp delta
        int32 timeDelta = int32(uint32(uint256(block.prevrandao)));
        // Real timestamp is current block time minus the delta
        return uint256(int256(block.timestamp*1000) - timeDelta);
    }
}