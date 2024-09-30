// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.20;

// implement this interface if you want to configure the visibility rules of your smart contract
// the TEN platform will interpret this information
interface ContractTransparencyConfig {
    // configuration per event log type
    enum Field{
        TOPIC1, TOPIC2, TOPIC3,
        SENDER, // tx.origin - msg.sender
        EVERYONE // the event is public - visible to everyone
    }

    enum ContractCfg{
        TRANSPARENT, //the internal state via getStorageAt will be accessible to everyone. All events will be public. This is the strongest setting.
        PRIVATE // internal state is hidden, and events can be configured.
    }

    // configuration per event log type
    struct EventLogConfig {
        uint256 eventSignature;
        Field[] visibleTo;
    }

    struct VisibilityConfig {
        ContractCfg contractCfg;
        EventLogConfig[] eventLogConfigs;  // mapping from event signature to visibility configs per event
    }

    // keep the logic independent of the environment
    // max gas: 1 Million
    function visibilityRules() external pure returns (VisibilityConfig memory);
}
