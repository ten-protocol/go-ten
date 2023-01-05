// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "../IMessageBus.sol";
import "../Structs.sol";

interface ICrossChainMessenger {
    struct CrossChainCall {
        address target;
        bytes data;
        uint256 gas;
    }

    // Returns the address of the message bus which is being used to verify messages.
    function messageBus() external view returns (address);

    // Returns the message.sender of the current message that is being relayed.
    function crossChainSender() external view returns (address);

    // Relays messages that exist in the message bus and have not been already consumed to the
    // target encdoded in the message, with the params encoded in the message.
    function relayMessage(Structs.CrossChainMessage calldata message) external;
}
