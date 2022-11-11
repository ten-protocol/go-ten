// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import * as MessageBus from "../messagebuscontract/Structs.sol";

interface Structs {
     // MetaRollup is a rollup meta data
    struct MetaRollup{
        bytes32 ParentHash;
        bytes32 Hash;
        address AggregatorID;
        bytes32 L1Block;
        uint256 Number;
    }

    // TreeElement is an element of the Tree structure
    struct TreeElement{
        uint256 ElementID;
        uint256 ParentID;
        MetaRollup rollup;
    }

    // NonExisting - 0 (Constant)
    // Tail - 1 (Constant)
    // Head - X (Variable)
    // Does not use rollup hashes as a storing ID as they can be compromised
    struct Tree {
        // rollups stores the Elements using incremental IDs
        mapping(uint256 => TreeElement) rollups;
        // map a rollup hash to a storage ID
        mapping(bytes32 => uint256) rollupsHashes;
        // map the children of a node
        mapping(uint256 => uint256[]) rollupChildren;

        uint256 _TAIL; // tail is always 1
        uint256 _HEAD;
        uint256 _nextID; // TODO use openzeppelin counters
        bool initialized;
    }

    struct HeaderCrossChainData {
        uint256 blockNumber;
        bytes32 blockHash;
        MessageBus.Structs.CrossChainMessage[] messages;
    }
}