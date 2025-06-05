// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.28;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";


contract ConstantSupplyERC20 is ERC20 {
    constructor(string memory name, string memory symbol, uint256 initialSupply) 
    ERC20(name, symbol)
    {
        _mint(msg.sender, initialSupply);
    }
}
