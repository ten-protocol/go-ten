// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "../../cross_chain_messaging/common/ICrossChainMessenger.sol";
import "@openzeppelin/contracts/utils/SlotDerivation.sol";
import "@openzeppelin/contracts/utils/StorageSlot.sol";


/**
 * @title CrossChainEnabledTEN
 * @dev Abstract contract that provides cross-chain messaging capabilities for TEN protocol.
 * This contract serves as a base for implementing cross-chain functionality, providing
 * secure message passing and verification between different chains.
 * 
 * TODO: We need to upgrade the open zeppelin version to the 4.x that adds cross chain enabled
 */
abstract contract CrossChainEnabledTEN  {
    using SlotDerivation for string; // for slot derivation
    using StorageSlot for bytes32;

    string private constant MESSENGER_SLOT = "CrossChainEnabledTEN.messenger";
    string private constant MESSAGE_BUS_SLOT = "CrossChainEnabledTEN.messageBus";
    string private constant NONCE_SLOT = "CrossChainEnabledTEN.nonce";

    function messenger() public view returns (ICrossChainMessenger) {
        return ICrossChainMessenger(MESSENGER_SLOT.erc7201Slot().getAddressSlot().value);
    }

    function messageBus() public view returns (IMessageBus) {
        return IMessageBus(MESSAGE_BUS_SLOT.erc7201Slot().getAddressSlot().value);
    }

    // Open zeppelin library does not support uint64 thus we must oversize.
    function nonce() public view returns (uint256) {
        return NONCE_SLOT.erc7201Slot().getUint256Slot().value;
    }

    // The messenger contract passed will be the authority that we trust to tell us
    // who has sent the cross chain message and that the message is indeed cross chain.
    function configure(address messengerAddress) internal {
        MESSENGER_SLOT.erc7201Slot().getAddressSlot().value = messengerAddress;
        MESSAGE_BUS_SLOT.erc7201Slot().getAddressSlot().value = messenger().messageBus();
        NONCE_SLOT.erc7201Slot().getUint256Slot().value = 0;
    }

    // Returns if the message is considered to be a cross chain one.
    function _isCrossChain() internal view returns (bool) {
        return msg.sender == address(messenger());
    }

    function _messageBus() internal view returns (IMessageBus) {
        return messageBus();
    }

    // Returns the address of the sender of the current cross chain message.
    // address 0x0 is considered null/no sender.
    function _crossChainSender() internal view returns (address) {
        return messenger().crossChainSender();
    }

    // Returns the address of the target of the current cross chain message.
    function _crossChainTarget() internal view returns (address) {
        return messenger().crossChainTarget();
    }

    // Ensures that the message is coming from another chain and for this contract.
    // Combined usage of _isCrossChain and _crossChainSender prevents attacks where
    // a message is relayed from the correct sender to a different contract that
    // maliciously calls into this contract.
    modifier onlyCrossChainSender(address sender) {
        require(
            _isCrossChain(),
            "Contract caller is not the registered messenger!"
        );
        require(
            _crossChainSender() == sender,
            "Cross chain message coming from incorrect sender!"
        );
        require(
            _crossChainTarget() == address(this),
            "Cross chain message coming from incorrect target!"
        );
        _;
    }

    // This function should be called to queue messages that are encodings of function calls
    // using the abi.encodeWithSelector(bytes4, arg) pattern.
    // target - the contract on which the call will be executed.
    // message - the result of abi.encodeWithSelector(bytes4, arg)
    // topic - TODO: determine if should be removed.
    // gas - when doing target.call{} this gas will be the limit set for the call. *Currently Unused*
    // consinstencyLevel - Block confirmations before finalizing message. Obscuro allows that even 0 is secure,
    // but it might not be suitable for all protocols.
    function queueMessage(
        address target,
        bytes memory message,
        uint32 topic,
        uint256 gas,
        uint8 consistencyLevel, 
        uint256 value
    ) internal {
        require(target != address(0), "Target cannot be 0x0");
        bytes memory payload = abi.encode(
            ICrossChainMessenger.CrossChainCall(target, message, gas)
        );
        messageBus().publishMessage{value: value}(uint64(NONCE_SLOT.erc7201Slot().getUint256Slot().value++), topic, payload, consistencyLevel);
    }

    // This function is used to publish a raw message instead of a function call to be relayed.
    // This is useful for sending arbitrary data across chains.
    function publishRawMessage(
        bytes memory message,
        uint32 topic,
        uint256 fee,
        uint8 consistencyLevel
    ) internal {
        messageBus().publishMessage{value: fee}(uint64(NONCE_SLOT.erc7201Slot().getUint256Slot().value++), topic, message, consistencyLevel);
    }
}
