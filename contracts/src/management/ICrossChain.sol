// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "../messaging/Structs.sol" as MessageStructs;

interface ICrossChain {
    event WithdrawalProcessed(bytes32 indexed messageHash, address indexed receiver, uint256 amount);
    event WithdrawalFailed(bytes32 indexed messageHash, string reason);
    event WithdrawalsPaused(bool paused);

    function extractNativeValue(
        MessageStructs.Structs.ValueTransferMessage calldata msg,
        bytes32[] calldata proof,
        bytes32 root
    ) external;
    function pauseWithdrawals(bool pause) external;
    function isWithdrawalSpent(bytes32 messageHash) external view returns (bool);
    function isBundleAvailable(bytes[] memory crossChainHashes) external view returns (bool);
    function getChallengePeriod() external view returns (uint256);
    function setChallengePeriod(uint256 delay) external;
}