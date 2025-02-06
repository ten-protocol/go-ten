// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import * as MessageBusStructs from "../messaging/Structs.sol";

interface Structs {
    struct MetaRollup{
        bytes32 Hash;
        uint256 LastSequenceNumber;

        bytes32 BlockBindingHash;
        uint256 BlockBindingNumber;
        bytes32 crossChainRoot;
        bytes Signature;
    }

    struct RollupStorage {
        mapping(bytes32=>MetaRollup) byHash;
        uint256 nextFreeSequenceNumber;
    }

    struct HeaderCrossChainData {
        MessageBusStructs.Structs.CrossChainMessage[] messages;
    }
}