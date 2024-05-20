// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import "./Structs.sol";

interface IMerkleTreeMessageBus {
    // This function is called to add a cross chain state root to the message bus.
    function addStateRoot(bytes32 stateRoot, uint256 activationTime) external;
    // This function disables a cross chain state root from the message bus. On challenge
    function disableStateRoot(bytes32 stateRoot) external;
    // This function verifies that a cross chain message is included in the state root.
    // message - the message to verify
    // proof - merkle tree proof for the said message against root
    // root - the state root to verify against. The contract checks that such a root has been added
    function verifyMessageInclusion(Structs.CrossChainMessage calldata message, bytes32[] memory proof, bytes32 root) external view;
    // This function verifies that a value transfer message is included in the state root.
    // arguments are same as the ones for message inclusion
    function verifyValueTransferInclusion(Structs.ValueTransferMessage calldata message, bytes32[] calldata proof, bytes32 root) external view;
}