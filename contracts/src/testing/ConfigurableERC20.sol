// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.28;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

/**
 * @title ConfigurableERC20
 * @dev ERC20 token with configurable decimals and fixed initial supply
 * Useful for deploying test tokens like USDC (6 decimals) or USDT (6 decimals)
 */
contract ConfigurableERC20 is ERC20 {
    uint8 private _decimals;

    /**
     * @dev Constructor that mints initialSupply tokens to the deployer
     * @param name Token name (e.g., "USD Coin")
     * @param symbol Token symbol (e.g., "USDC")
     * @param decimals_ Number of decimals (e.g., 6 for USDC, 18 for most tokens)
     * @param initialSupply Initial token supply (already adjusted for decimals)
     */
    constructor(
        string memory name,
        string memory symbol,
        uint8 decimals_,
        uint256 initialSupply
    ) ERC20(name, symbol) {
        require(decimals_ <= 18, "Decimals must be <= 18");
        _decimals = decimals_;
        _mint(msg.sender, initialSupply);
    }

    /**
     * @dev Returns the number of decimals used by the token
     * @return The number of decimals
     */
    function decimals() public view virtual override returns (uint8) {
        return _decimals;
    }
}
