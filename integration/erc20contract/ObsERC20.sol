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
        // Respond to balance requests made by the account owner only.
        // In case the requester spoofs the "from" of the call, they will not be able to read
        // the result since it will be returned encrypted with the viewing key of the declared "from".
        require(msg.sender == account);
        return super.balanceOf(account);
    }

}