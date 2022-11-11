// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

interface Structs {
    struct CrossChainMessage {
        address sender;
        uint64  sequence;
        uint32  nonce;
        bytes   topic;
        bytes   payload;
    }
}