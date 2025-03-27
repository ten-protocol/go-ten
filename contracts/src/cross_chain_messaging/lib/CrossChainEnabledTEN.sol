// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "../../cross_chain_messaging/common/ICrossChainMessenger.sol";
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";


/**
 * @title CrossChainEnabledTEN
 * @dev Abstract contract that provides cross-chain messaging capabilities for TEN protocol.
 * This contract serves as a base for implementing cross-chain functionality, providing
 * secure message passing and verification between different chains.
 * 
 * TODO: We need to upgrade the open zeppelin version to the 4.x that adds cross chain enabled
 */
abstract contract CrossChainEnabledTEN is Initializable {
    ICrossChainMessenger messenger;
    IMessageBus messageBus;
    uint32 nonce = 0;

    /**
     * @dev Configures the contract with the messenger address that will be trusted for cross-chain communication
     * @param messengerAddress The address of the messenger contract that will act as the authority
     * for cross-chain message verification and sender authentication
     */
    function configure(address messengerAddress) public onlyInitializing {
        messenger = ICrossChainMessenger(messengerAddress);
        messageBus = IMessageBus(messenger.messageBus());
    }

    /**
     * @dev Checks if the current call is a cross-chain message
     * @return bool True if the message is from the trusted messenger contract
     */
    function _isCrossChain() internal view returns (bool) {
        return msg.sender == address(messenger);
    }

    /**
     * @dev Returns the message bus instance used for cross-chain communication
     * @return IMessageBus The message bus interface
     */
    function _messageBus() internal view returns (IMessageBus) {
        return messageBus;
    }

    /**
     * @dev Returns the address of the sender from the originating chain
     * @return address The cross-chain sender's address (0x0 if no sender/null)
     */
    function _crossChainSender() internal view returns (address) {
        return messenger.crossChainSender();
    }

    /**
     * @dev Modifier that ensures the message is from a specific cross-chain sender
     * This provides a double security check:
     * 1. Verifies the message is truly cross-chain (via _isCrossChain)
     * 2. Validates the specific sender address (via _crossChainSender)
     * 
     * This combination prevents relay attacks where a valid message from a correct sender
     * might be maliciously redirected to a different contract.
     * 
     * @param sender The expected sender address from the originating chain
     */
    modifier onlyCrossChainSender(address sender) {
        require(
            _isCrossChain(),
            "Contract caller is not the registered messenger!"
        );
        require(
            _crossChainSender() == sender,
            "Cross chain message coming from incorrect sender!"
        );
        _;
    }

    /**
     * @dev Queues a cross-chain message for execution on the target chain
     * @param target The contract address that will execute the message on the target chain
     * @param message The encoded function call (using abi.encodeWithSelector)
     * @param topic Message topic identifier
     * @param gas Gas limit for the target contract call (currently unused)
     * @param consistencyLevel Number of block confirmations required before message finalization. While TEN supports 0 confirmations securely, other protocols may require more
     * @param value Native token value to be sent with the message
     */
    function queueMessage(
        address target,
        bytes memory message,
        uint32 topic,
        uint256 gas,
        uint8 consistencyLevel, 
        uint256 value
    ) internal {
        bytes memory payload = abi.encode(
            ICrossChainMessenger.CrossChainCall(target, message, gas)
        );
        messageBus.publishMessage{value: value}(nonce++, topic, payload, consistencyLevel);
    }
}
