// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";

contract UnrenouncableOwnable2Step is Ownable2StepUpgradeable {

    function renounceOwnership() public view override onlyOwner {
        revert("UnrenouncableOwnable2Step: cannot renounce ownership");
    }
}