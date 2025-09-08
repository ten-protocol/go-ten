// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";

contract UpgradeManager is Initializable, Ownable2StepUpgradeable {
    constructor() {
        _disableInitializers();
    }

    function initialize(address owner) public initializer {
        __Ownable2Step_init(owner);
    }

    function upgradeFeature(string calldata featureName, bytes calldata featureData) external onlyOwner {
        emit Upgraded(featureName, featureData);
    }
}