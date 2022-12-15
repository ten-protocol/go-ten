// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "./ICrossChainMessenger.sol";

contract CrossChainMessenger is ICrossChainMessenger {
    
    IMessageBus messageBusContract;
    address public crossChainSender = address(0x0);
    mapping (bytes32 => bool) messageConsumed;

    constructor(address messageBusAddr) {
        messageBusContract = IMessageBus(messageBusAddr);
    }

    function messageBus() external view returns (address) {
        return address(messageBusContract);
    }

    function consumeMessage(Structs.CrossChainMessage calldata message) private {
        require(messageBusContract.verifyMessageFinalized(message), "Message not found or finalized.");
        bytes32 msgHash = keccak256(abi.encode(message));
        require(messageConsumed[msgHash] == false, "Message already consumed.");

        messageConsumed[msgHash] = true;
    }

    function _getRevertMsg(bytes memory _returnData) internal pure returns (string memory) {
        // If the _res length is less than 68, then the transaction failed silently (without a revert message)
        if (_returnData.length < 68) return 'Transaction reverted silently';

        assembly {
            // Slice the sighash.
            _returnData := add(_returnData, 0x04)
        }
        return abi.decode(_returnData, (string)); // All that remains is the revert string
    }

    function encodeCall(address target, bytes calldata payload) public pure returns (bytes memory) {
        return abi.encode(CrossChainCall(target, payload, 0));
    }

    function relayMessage(Structs.CrossChainMessage calldata message) public {
        consumeMessage(message);

        crossChainSender = message.sender;

        CrossChainCall memory callData = abi.decode(message.payload, (CrossChainCall));
        (bool success,  bytes memory returnData ) = callData.target.call(callData.data);        
        if (!success) {
            //This doesn't seem to really work ...
            revert(_getRevertMsg(returnData));
        }

        crossChainSender = address(0x0);
    }
}