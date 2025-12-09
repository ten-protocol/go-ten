// SPDX-License-Identifier: BUSL
// Compatible with OpenZeppelin Contracts ^5.5.0
pragma solidity ^0.8.27;

import {ERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import {ERC20PermitUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC20PermitUpgradeable.sol";
import {ERC20VotesUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC20VotesUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {NoncesUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/NoncesUpgradeable.sol";
import {Erc20WithLock} from "./ERC20WithLock.sol";

/// @custom:security-contact security@ten.xyz
contract Ten is Initializable, Erc20WithLock, ERC20PermitUpgradeable, ERC20VotesUpgradeable {

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize() public initializer {
        __ERC20_init("TEN", "TEN");
        __ERC20Permit_init("TEN");
        __ERC20Votes_init();
        __Ownable_init(msg.sender);
        __AccessControl_init();
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
    }

    // The following functions are overrides required by Solidity.

    function transfer(address to, uint256 value)
    public
    override(ERC20Upgradeable, Erc20WithLock)
    returns (bool)
    {
        return super.transfer(to, value);
    }

    function transferFrom(address from, address to, uint256 value)
    public
    override(ERC20Upgradeable, Erc20WithLock)
    returns (bool)
    {
        return super.transferFrom(from, to, value);
    }

    function _update(address from, address to, uint256 value)
    internal
    override(ERC20Upgradeable, ERC20VotesUpgradeable)
    {
        super._update(from, to, value);
    }

    function nonces(address owner) public view override(ERC20PermitUpgradeable, NoncesUpgradeable) returns (uint256) {
        return super.nonces(owner);
    }
}