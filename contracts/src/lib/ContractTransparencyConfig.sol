// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.20;

// implement this interface if you want to configure the visibility rules of your smart contract
// the TEN platform will interpret this information
interface ContractTransparencyConfig {
    // configuration per event log type
    struct EventLogConfig {
        bytes eventSignature;
        bool isPublic;  // everyone can see and query for this event
        bool topic1CanView;    // If the event is private, and this is true, it means that the address from topic1 is an EOA that can view this event
        bool topic2CanView;    // same
        bool topic3CanView;    // same
        bool visibleToSender; // if true, the tx signer will see this event. Default false
    }

    struct VisibilityConfig {
        bool isTransparent; // If true - the internal state via getStorageAt will be accessible to everyone. All events will be public. Default false
        EventLogConfig[] eventLogConfigs;  // mapping from event signature to visibility configs per event
    }

    // keep the logic independent of the environment
    // max gas: 1 Million
    function visibilityRules() external pure returns (VisibilityConfig memory);
}
