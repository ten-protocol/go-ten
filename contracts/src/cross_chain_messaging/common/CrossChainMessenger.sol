// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "./ICrossChainMessenger.sol";
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import "../L1/IMerkleTreeMessageBus.sol";


/**
 * @title CrossChainMessenger
 * @dev contract that provides the context for contracts
 * that inherit CrossChainEnabledTEN. It allows to deliver messages using relayMessage.
 * This contract abstracts the message verification logic away from the CrossChainEnabled contracts.
 * It works by querying the message bus that has been passed in the constructor and calling functions
 * that have been encoded with abi.encodeWithSelector(bytes4, arg).
 * It's also responsible for marking messages as consumed whenever a successful call happens. This means
 * that CrossChainEnabled contracts need not bother with anything related to verification, apart from confirming
 * from whom the messages are coming from.
 * Notice that this Messenger has no restrictions on who can relay messages, nor does it have any understanding of fees.
 * You can opt in to deploy a customer messenger for your cross chain dApp with more specialized logic.
 */
contract CrossChainMessenger is ICrossChainMessenger, Initializable {
    error CallFailed(bytes error);

    IMerkleTreeMessageBus messageBusContract;
    address public crossChainSender;
    mapping(bytes32 => bool) messageConsumed;

    /**
     * @dev Initializes the contract with a message bus address
     * @param messageBusAddr The address of the message bus contract
     * TODO initialize only once
     */
    function initialize(address messageBusAddr) external initializer {
        messageBusContract = IMerkleTreeMessageBus(messageBusAddr);
        crossChainSender =  address(0x0);
    }

    /**
     * @dev Returns the address of the message bus contract
     * @return address The message bus contract address
     */
    function messageBus() external view returns (address) {
        return address(messageBusContract);
    }

    /**
     * @dev Verifies the message exists and has not already been consumed
     * @param message The cross-chain message to consume
     */
    function consumeMessage(
        Structs.CrossChainMessage calldata message
    ) private {
        require(
            IMessageBus(address(messageBusContract)).verifyMessageFinalized(message),
            "Message not found or finalized."
        );
        bytes32 msgHash = keccak256(abi.encode(message));
        require(messageConsumed[msgHash] == false, "Message already consumed.");

        messageConsumed[msgHash] = true;
    }

    /**
     * @dev Verifies and consumes a cross-chain message using a Merkle proof
     * @param message The cross-chain message to consume
     * @param proof Merkle proof verifying the message inclusion
     * @param root The Merkle root against which to verify the proof
     */
    function consumeMessageWithProof(
        Structs.CrossChainMessage calldata message,
        bytes32[] calldata proof, 
        bytes32 root
    ) private {
        messageBusContract.verifyMessageInclusion(message, proof, root);
        bytes32 msgHash = keccak256(abi.encode(message));
        require(messageConsumed[msgHash] == false, "Message already consumed.");

        messageConsumed[msgHash] = true;
    }

    /**
     * @dev Encodes a cross-chain call for testing purposes
     * @param target The address of the contract to call
     * @param payload The calldata to send
     * @return bytes The encoded cross-chain call
     * TODO: Remove this. It does not serve any real purpose on chain, but is currently required for hardhat tests
     */
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
    /**
     * @dev This function can be called by anyone and if the message @param actually exists in the message bus, then the
     * function will push it to the targeted smart contract. Notice that anyone can queue a call to be relayed, but the
     * cross chain sender is set to be the address of the message sender on the other layer, as it is when reaching the
     * message bus.
     * @param message The cross-chain message to relay
     */
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
        (bool success, bytes memory returnData) = callData.target.call{gas: gasleft()}(
            callData.data
        );
        if (!success) {
            revert CallFailed(returnData);
        }

        crossChainSender = address(0x0);
    }

    /**
     * @dev Relays a cross-chain message verified by a Merkle proof
     * @param message The cross-chain message to relay
     * @param proof Merkle proof verifying the message inclusion
     * @param root The Merkle root against which to verify the proof
     */
    function relayMessageWithProof(Structs.CrossChainMessage calldata message, bytes32[] calldata proof, bytes32 root) public {
        consumeMessageWithProof(message, proof, root);

        crossChainSender = message.sender;

        //TODO: Do not relay to self. Do not relay to known contracts. Consider what else not to talk to.
        //Add reentracy guards and paranoid security checks as messenger contracts will have above average rights
        //when communicating with other contracts.

        CrossChainCall memory callData = abi.decode(
            message.payload,
            (CrossChainCall)
        );
        (bool success, bytes memory returnData) = callData.target.call{gas: gasleft()}(
            callData.data
        );
        if (!success) {
            revert CallFailed(returnData);
        }

        crossChainSender = address(0x0);
    }
}
