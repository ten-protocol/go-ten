// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

/**
 * @title UnrenouncableOwnable2Step
 * @dev Contract that extends Ownable2StepUpgradeable but prevents renouncing ownership
 */
contract UnrenouncableOwnable2Step is Ownable2StepUpgradeable {
    /**
     * @dev Initializes the contract setting the deployer as the initial owner.
     * @param initialOwner The address that will be the initial owner
     */
    function __UnrenouncableOwnable2Step_init(address initialOwner) internal onlyInitializing {
        __Ownable2Step_init();
        _transferOwnership(initialOwner);
    }

    /**
     * @dev Prevents renouncing ownership
     */
    function renounceOwnership() public view override onlyOwner {
        revert("UnrenouncableOwnable2Step: cannot renounce ownership");
    }
}