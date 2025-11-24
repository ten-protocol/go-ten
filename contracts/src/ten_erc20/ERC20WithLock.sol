// SPDX-License-Identifier: BUSL
// Compatible with OpenZeppelin Contracts ^5.5.0
pragma solidity ^0.8.27;

import {ERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {AccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";

/**
 * @dev Abstract ERC20 extension that adds simple account-level locking and an allowlist role
 *      for destinations that may receive transfers from fully locked accounts.
 *
 *      - `VESTING_ROLE` recipients are allowed to receive from locked senders.
 *      - `locked[address]` marks an account as fully locked.
 *      - `initializeAllocations` mints initial balances and can pre-lock accounts (one-time).
 *      - `createLockedAccount` lets any user opt-in to being lock-eligible to receive locked transfers.
 *
 *      Access:
 *      - Owner may call `grantVesting`/`revokeVesting`.
 *      - The contract expects the inheriting token to call initializers for Ownable/AccessControl.
 */
abstract contract Erc20WithLock is ERC20Upgradeable, OwnableUpgradeable, AccessControlUpgradeable {
    // Role for addresses allowed to receive transfers from fully locked accounts
    bytes32 public constant VESTING_ROLE = keccak256("VESTING_ROLE");

    // Set of addresses whose tokens are completely locked
    mapping(address => bool) public locked;

    // Guard to ensure initial mint/lock configuration can only be executed once
    bool private _initialLocksConfigured;

    /**
     * @dev Initial mint plan entry.
	 * - `account`: recipient address
	 * - `amount`: amount to mint
	 * - `locked`: whether to mark the recipient as locked
	 */
    struct InitialAllocation {
        address account;
        uint256 amount;
        bool    locked;
    }

    /**
     * @dev Combined one-time initializer to mint initial supply and optionally lock accounts.
	 * Each allocation can be locked or unlocked based on the struct flag.
	 * Callable only once by the owner.
	 */
    function initializeAllocations(InitialAllocation[] memory allocations) external onlyOwner {
        require(!_initialLocksConfigured, "Initial allocations already configured");
        for (uint256 i = 0; i < allocations.length; i++) {
            address account = allocations[i].account;
            uint256 amount  = allocations[i].amount;
            require(account != address(0), "Invalid account");
            require(amount > 0, "Amount must be positive");
            _mint(account, amount);
            if (allocations[i].locked) {
                locked[account] = true;
            }
        }
        _initialLocksConfigured = true;
    }

    /**
     * @dev Owner-managed helpers to grant/revoke vesting role.
	 */
    function grantVesting(address account) external onlyOwner {
        require(account != address(0), "Invalid account");
        grantRole(VESTING_ROLE, account);
    }

    function revokeVesting(address account) external onlyOwner {
        revokeRole(VESTING_ROLE, account);
    }

    /**
     * @dev Override transfer to check locked status
	 */
    function transfer(address to, uint256 value) public virtual override returns (bool) {
        address from = _msgSender();
        _handleLockedTransfer(from, to, value);
        return super.transfer(to, value);
    }

    /**
     * @dev Override transferFrom to check locked status
	 */
    function transferFrom(address from, address to, uint256 value) public virtual override returns (bool) {
        _handleLockedTransfer(from, to, value);
        return super.transferFrom(from, to, value);
    }

    /**
     * @dev Shared logic for enforcing locked account transfer rules.
	 */
    function _handleLockedTransfer(address from, address to, uint256 value) internal {
        bool fromIsVesting = hasRole(VESTING_ROLE, from);
        bool toIsVesting = hasRole(VESTING_ROLE, to);
        bool toIsLocked = locked[to];
        bool fromIsLocked = locked[from];

        if (fromIsVesting){
            // can go back to locked accounts
            // or can go to buyers
            return;
        }

        if (toIsVesting){
            return;
        }

        // marketplace interactions have been handled already so no need to worry about them

        if (fromIsLocked){
            // locked senders may transfer only to locked recipients (OTC) or vesting (handled above)
            require (toIsLocked, "Locked: transfer only to vesting");
            // If transferring full balance to another locked account, clear sender lock
            if (balanceOf(from) == value) {
                delete locked[from];
            }
            return;
        }

        if (toIsLocked) {
            // can only receive from locked accounts
            require (fromIsLocked, "Can only transfer locked tokens to other locked accounts.");
            if (balanceOf(from) == value) {
                delete locked[from];
            }
            return;
        }

        return;
    }

    // Note: For OTC flows, locked-to-locked transfers are supported via standard transfer/transferFrom

    /**
     * @dev Allows any account to lock itself to authorize receiving locked transfers.
	 * Once locked, the account cannot unlock itself.
	 */
    function createLockedAccount() external returns (bool) {
        require(!locked[msg.sender], "Already locked");
        locked[msg.sender] = true;
        return true;
    }
}
