// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

interface Structs {
    struct CrossChainMessage {
        address sender;
        uint64  sequence;
        uint32  nonce;
        uint32  topic;
        bytes   payload;
        uint8   consistencyLevel;
    }
}