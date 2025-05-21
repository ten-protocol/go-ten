// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import "../../common/Structs.sol";
import "../../system/interfaces/IFees.sol";
import "./IMessageBus.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "../../common/UnrenouncableOwnable2Step.sol";

/// @title BaseMessageBus
/// @dev Implements the chain-agnostic publish API shared by L1 and L2 buses.
abstract contract BaseMessageBus is IMessageBus, Initializable, UnrenouncableOwnable2Step {
    IFees internal fees;
    mapping(address sender => uint64 sequence) internal addressSequences;

    constructor() {
        _transferOwnership(msg.sender);
    }

    /// @dev Initializes the contract with an owner and optional fee contract.
    function initialize(address caller, address feesAddress) public virtual initializer {
        __Ownable_init(caller);
        fees = IFees(feesAddress);
    }

    /// @dev Increments and returns the sequence number for a sender.
    function incrementSequence(address sender) internal returns (uint64 sequence) {
        sequence = addressSequences[sender];
        addressSequences[sender] += 1;
    }

    /// @inheritdoc IMessageBus
    function getPublishFee() public view virtual override returns (uint256) {
        return fees.messageFee();
    }

    /// @inheritdoc IMessageBus
    function publishMessage(
        uint32 nonce,
        uint32 topic,
        bytes calldata payload,
        uint8 consistencyLevel
    ) external payable virtual override returns (uint64 sequence) {
        if (address(fees) != address(0)) {
            uint256 fee = getPublishFee();
            require(msg.value >= fee, "Insufficient funds to publish message");
            (bool ok, ) = address(fees).call{value: fee}("");
            require(ok, "Failed to send fees to fees contract");
        }

        sequence = incrementSequence(msg.sender);
        emit LogMessagePublished(
            msg.sender,
            sequence,
            nonce,
            topic,
            payload,
            consistencyLevel
        );
        return sequence;
    }

    /// @inheritdoc IMessageBus
    function retrieveAllFunds(address receiver) external virtual override onlyOwner {
        (bool ok, ) = receiver.call{value: address(this).balance}("");
        require(ok, "failed sending value");
    }
}
