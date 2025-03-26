// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.20;

/**
 * @title IContractTransparencyConfig
 * @dev Implement this interface if you want to configure the visibility rules of your smart contract. the TEN platform will interpret this information
 */
interface IContractTransparencyConfig {

    /**
     * @dev Enum for the visibility of a field
     */
    enum Field{
        TOPIC1, TOPIC2, TOPIC3, // if any of these fields is in the relevantTo array, then the address in that topic will be able to query for that event
        SENDER, // the tx.origin will be able to query for the event
        EVERYONE // the event is public - visible to everyone
    }

    /**
     * @dev Enum for the visibility of a contract
     */
    enum ContractCfg{
        TRANSPARENT, //the internal state via getStorageAt will be accessible to everyone. All events will be public. This is the strongest setting.
        PRIVATE // internal state is hidden, and events can be configured.
    }

    /**
     * @dev Struct for the visibility of an event
     */
    struct EventLogConfig {
        bytes32 eventSignature;
        Field[] visibleTo;
    }

    /**
     * @dev Struct for the visibility of a contract
     */
    struct VisibilityConfig {
        ContractCfg contractCfg;
        EventLogConfig[] eventLogConfigs;  // mapping from event signature to visibility configs per event
    }

    /**
     * @dev Returns the visibility rules for the contract. Keep the logic independent of the environment.
     * Maximum gas: 1 Million
     * @return VisibilityConfig The visibility rules for the contract.
     */
    function visibilityRules() external pure returns (VisibilityConfig memory);
}
