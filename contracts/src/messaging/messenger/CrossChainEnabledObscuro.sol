// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "./ICrossChainMessenger.sol";

// TODO: We need to upgrade the open zeppelin version to the 4.x that adds cross chain enabled
abstract contract CrossChainEnabledObscuro {
    
    ICrossChainMessenger messenger;
    IMessageBus messageBus;
    uint32 nonce = 0;

    constructor(address messengerAddress) {
        messenger = ICrossChainMessenger(messengerAddress);
        messageBus = IMessageBus(messenger.messageBus());
    }

    function _isCrossChain() internal view returns (bool) {
        return msg.sender == address(messenger);
    }

    function _crossChainSender() internal view returns (address) {
        return messenger.crossChainSender();
    }

    modifier onlyCrossChainSender(address sender) {
        require(_isCrossChain(), "Contract caller is not the registered messenger!");
        require(_crossChainSender() == sender, "Cross chain message coming from incorrect sender!");
        _;
    }

    function queueMessage(address target, bytes memory message, uint32 topic, uint256 gas, uint8 consistencyLevel) internal {
        bytes memory payload = abi.encode(ICrossChainMessenger.CrossChainCall(target, message, gas));
        messageBus.publishMessage(nonce++, topic, payload, consistencyLevel);
    }
}