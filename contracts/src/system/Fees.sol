// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

interface IFees {
    function messageFee(uint256 messageSize) external view returns (uint256);
}

// Contract that will contain fees for contracts that need to apply them
contract Fees is Initializable, OwnableUpgradeable {

    uint256 private _messageFeePerByte;

    constructor() {
        _disableInitializers();
    }

    function initialize(uint256 initialMessageFeePerByte, address eoaOwner) public initializer {
        __Ownable_init(eoaOwner);
        _messageFeePerByte = initialMessageFeePerByte;
    }

    function messageFee(uint256 messageSize) external view returns (uint256) {
        return _messageFeePerByte * messageSize;
    }

    function setMessageFee(uint256 newMessageFeePerByte) external onlyOwner{
        _messageFeePerByte = newMessageFeePerByte;
    }

    function withdrawalCollectedFees() external onlyOwner {
        payable(owner()).transfer(address(this).balance);
    }

    function collectedFees() external view returns (uint256) {
        return address(this).balance;
    }

    // For now channel all here as we only collect fees
    // from the message bus
    receive() external payable {
    }
}
