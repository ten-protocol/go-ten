// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "../interfaces/IFees.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "../../common/UnrenouncableOwnable2Step.sol";

/**
 * @title Fees
 * @dev Contract that will contain fees for contracts that need to apply them
 * 
 * TODO stefan to add explanation
 */
contract Fees is IFees, Initializable, UnrenouncableOwnable2Step {

    uint256 private _messageFee;

    event FeeChanged(uint256 oldFee, uint256 newFee);
    event FeeWithdrawn(uint256 amount);

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @dev Initializes the contract with fees and owner
     * @param flatFee Initial message fee
     * @param eoaOwner Address that will own the contract
     */
    function initialize(uint256 flatFee, address eoaOwner) public initializer {
        __UnrenouncableOwnable2Step_init(eoaOwner);
        _messageFee = flatFee;
        emit FeeChanged(0, flatFee);
    }

    // Helper function to calculate the fee for a message
    function messageFee() external view returns (uint256) {
        return _messageFee;
    }

    // The EOA owner can set the message fee to ensure sequencer is not publishing
    // at a loss
    function setMessageFee(uint256 newFeeForMessage) external onlyOwner{
        uint256 oldFee = _messageFee;
        _messageFee = newFeeForMessage;
        emit FeeChanged(oldFee, newFeeForMessage);
    }

    // The EOA owner can collect the fees
    function withdrawalCollectedFees() external onlyOwner {
        uint256 balance = address(this).balance;
        (bool success, ) = payable(owner()).call{value: balance}("");
        require(success, "Failed to send Ether");
        emit FeeWithdrawn(balance);
    }

    // For now just the whole balance as we only collect fees
    // from the message bus
    function collectedFees() external view returns (uint256) {
        return address(this).balance;
    }

    // For now channel all here as we only collect fees
    // from the message bus
    receive() external payable {
    }
}
