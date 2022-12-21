// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "./ICrossChainMessenger.sol";

// CrossChainMessenger is the contract that provides the context for contracts
// that inherit CrossChainEnabledObscuro. It allows to deliver messages using relayMessage.
// This contract abstracts the message verification logic away from the CrossChainEnabled contracts.
// It works by querying the message bus that has been passed in the constructor and calling functions
// that have been encoded with abi.encodeWithSelector(bytes4, arg).
// It's also responsible for marking messages as consumed whenever a successful call happens. This means
// that CrossChainEnabled contracts need not bother with anything related to verification, apart from confirming
// from whom the messages are coming from.
// Notice that this Messenger has no restrictions on who can relay messages, nor does it have any understanding of fees.
// You can opt in to deploy a customer messenger for your cross chain dApp with more specialized logic.
contract CrossChainMessenger is ICrossChainMessenger {
    error CallFailed(bytes error);

    IMessageBus messageBusContract;
    address public crossChainSender = address(0x0);
    mapping(bytes32 => bool) messageConsumed;

    constructor(address messageBusAddr) {
        messageBusContract = IMessageBus(messageBusAddr);
    }

    function messageBus() external view returns (address) {
        return address(messageBusContract);
    }

    // Will verify that the message exists & has not been already consumed and will
    // mark it as consumed.
    function consumeMessage(
        Structs.CrossChainMessage calldata message
    ) private {
        require(
            messageBusContract.verifyMessageFinalized(message),
            "Message not found or finalized."
        );
        bytes32 msgHash = keccak256(abi.encode(message));
        require(messageConsumed[msgHash] == false, "Message already consumed.");

        messageConsumed[msgHash] = true;
    }

    // TODO: Remove this. It does not serve any real purpose on chain, but is currently required for hardhat tests
    // as producing the same result in JS has proven difficult...
    function encodeCall(
        address target,
        bytes calldata payload
    ) public pure returns (bytes memory) {
        return abi.encode(CrossChainCall(target, payload, 0));
    }

    // This function can be called by anyone and if the message @param actually exists in the message bus,
    // then the function will push it to the targeted smart contract.
    // Notice that anyone can queue a call to be relayed, but the cross chain sender is set to be
    // the address of the message sender on the other layer, as it is when reaching the message bus.
    function relayMessage(Structs.CrossChainMessage calldata message) public {
        consumeMessage(message);

        crossChainSender = message.sender;

        //TODO: Do not relay to self. Do not relay to known contracts. Consider what else not to talk to.
        //Add reentracy guards and paranoid security checks as messenger contracts will have above average rights
        //when communicating with other contracts.

        CrossChainCall memory callData = abi.decode(
            message.payload,
            (CrossChainCall)
        );
        (bool success, bytes memory returnData) = callData.target.call(
            callData.data
        );
        if (!success) {
            revert CallFailed(returnData);
        }

        crossChainSender = address(0x0);
    }
}
