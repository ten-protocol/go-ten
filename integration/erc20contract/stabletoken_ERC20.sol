// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract StableToken is ERC20 {
    constructor() ERC20("StableToken", "STK") {
        _mint(msg.sender, 1000000000 * 10 ** decimals());
    }
}