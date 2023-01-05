// SPDX-License-Identifier: Apache 2
pragma solidity >=0.7.0 <0.9.0;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "./ObsERC20.sol";

contract WrappedERC20 is ObsERC20, AccessControl {
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");

    constructor(string memory name, string memory symbol) ObsERC20(name, symbol) {
        _grantRole(ADMIN_ROLE, msg.sender);
    }

    function issueFor(
        address receiver,
        uint256 amount
    ) external onlyRole(ADMIN_ROLE) {
        _mint(receiver, amount);
    }

    function burnFor(
        address giver,
        uint256 amount
    ) external onlyRole(ADMIN_ROLE) {
        require(balanceOf(giver) >= amount, "Insufficient balance.");
        _burn(giver, amount);
    }
}
