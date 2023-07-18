// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import * as MessageBusStructs from "../messaging/Structs.sol";

interface Structs {
     // MetaRollup is a rollup meta data
    struct MetaRollup{
        bytes32 Hash;
        address AggregatorID;
        bytes32 L1Block;
        uint256 LastSequenceNumber;
    }

    struct RollupStorage {
        mapping(bytes32=>MetaRollup) byHash;
    }

    struct HeaderCrossChainData {
        uint256 blockNumber;
        bytes32 blockHash;
        MessageBusStructs.Structs.CrossChainMessage[] messages;
    }
}