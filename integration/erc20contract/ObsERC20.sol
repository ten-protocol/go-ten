// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "libs/openzeppelin/contracts/token/ERC20/ERC20.sol";

contract ObsERC20 is ERC20 {

    constructor(
        string memory name,
        string memory symbol,
        uint256 initialSupply
    )  ERC20(name, symbol) {
        _mint(msg.sender, initialSupply);
    }

    function balanceOf(address account) public view virtual override returns (uint256) {
// Todo - enable this to test ACL
//        require(tx.origin == account);
        return super.balanceOf(account);
    }

}