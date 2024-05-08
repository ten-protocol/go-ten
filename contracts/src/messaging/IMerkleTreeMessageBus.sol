// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import "./Structs.sol";

interface IMerkleTreeMessageBus {
    function addStateRoot(bytes32 stateRoot, uint256 activationTime) external;
    function blockStateRoot(bytes32 stateRoot) external;
    function verifyMessageInclusion(Structs.CrossChainMessage calldata message, bytes32[] memory proof, bytes32 root) external view;
}