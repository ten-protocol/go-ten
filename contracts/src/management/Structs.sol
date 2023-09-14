// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import * as MessageBusStructs from "../messaging/Structs.sol";

interface Structs {
    struct MetaRollup{
        bytes32 Hash;
        address AggregatorID;
        uint256 LastSequenceNumber;
    }

    struct RollupStorage {
        mapping(bytes32=>MetaRollup) byHash;
    }

    struct HeaderCrossChainData {
        MessageBusStructs.Structs.CrossChainMessage[] messages;
    }
}