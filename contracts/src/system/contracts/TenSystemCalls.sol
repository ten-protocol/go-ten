// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/proxy/utils/Initializable.sol";

contract TenSystemCalls is Initializable {

    function initialize() external initializer {
    }

    function getRandomNumber() external view returns (uint256) {
        return block.prevrandao; // We inject randomness in prevrandao
    }

    function getTransactionTimestamp() external view returns (uint256) {
        return block.difficulty; // We override block.difficulty in the sequencer to inject exact time of arrival
    }
}