// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

/**
 * @title PausableWithRoles
 * @dev Contract that implements pausable functionality with role-based pause/unpause permissions
 * 
 * Roles:
 * - PAUSER_ROLE: Can pause the contract (typically deployer key for quick response)
 * - UNPAUSER_ROLE: Can unpause the contract (typically multisig wallet for controlled recovery)
 */
abstract contract PausableWithRoles is Initializable, AccessControlUpgradeable {
    
    /// @dev Role that can pause the contract
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
    
    /// @dev Role that can unpause the contract
    bytes32 public constant UNPAUSER_ROLE = keccak256("UNPAUSER_ROLE");

    /// @dev Paused state
    bool private _paused;

    /// @dev Emitted when the pause is triggered by `account`.
    event Paused(address account);

    /// @dev Emitted when the pause is lifted by `account`.
    event Unpaused(address account);

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @dev Initializes the contract with the deployer as the initial pauser and unpauser
     * @param deployer The address that will have both PAUSER_ROLE and UNPAUSER_ROLE initially
     */
    function __PausableWithRoles_init(address deployer) internal onlyInitializing {
        __AccessControl_init();
        
        // Grant roles to deployer initially
        _grantRole(PAUSER_ROLE, deployer);
        _grantRole(UNPAUSER_ROLE, deployer);
        
        // Set role admin - only deployer can grant/revoke roles initially
        _setRoleAdmin(PAUSER_ROLE, DEFAULT_ADMIN_ROLE);
        _setRoleAdmin(UNPAUSER_ROLE, DEFAULT_ADMIN_ROLE);
    }

    /**
     * @dev Returns true if the contract is paused, and false otherwise.
     */
    function paused() public view virtual returns (bool) {
        return _paused;
    }

    /**
     * @dev Modifier to make a function callable only when the contract is not paused.
     */
    modifier whenNotPaused() {
        require(!paused(), "Pausable: paused");
        _;
    }

    /**
     * @dev Modifier to make a function callable only when the contract is paused.
     */
    modifier whenPaused() {
        require(paused(), "Pausable: not paused");
        _;
    }

    /**
     * @dev Pauses the contract
     * @notice Only callable by accounts with PAUSER_ROLE
     */
    function pause() external onlyRole(PAUSER_ROLE) {
        _paused = true;
        emit Paused(_msgSender());
    }

    /**
     * @dev Unpauses the contract
     * @notice Only callable by accounts with UNPAUSER_ROLE
     */
    function unpause() external onlyRole(UNPAUSER_ROLE) {
        _paused = false;
        emit Unpaused(_msgSender());
    }

    /**
     * @dev Grants PAUSER_ROLE to an address
     * @param account The address to grant the role to
     */
    function grantPauserRole(address account) external onlyRole(DEFAULT_ADMIN_ROLE) {
        _grantRole(PAUSER_ROLE, account);
    }

    /**
     * @dev Revokes PAUSER_ROLE from an address
     * @param account The address to revoke the role from
     */
    function revokePauserRole(address account) external onlyRole(DEFAULT_ADMIN_ROLE) {
        _revokeRole(PAUSER_ROLE, account);
    }

    /**
     * @dev Grants UNPAUSER_ROLE to an address
     * @param account The address to grant the role to
     */
    function grantUnpauserRole(address account) external onlyRole(DEFAULT_ADMIN_ROLE) {
        _grantRole(UNPAUSER_ROLE, account);
    }

    /**
     * @dev Revokes UNPAUSER_ROLE from an address
     * @param account The address to revoke the role from
     */
    function revokeUnpauserRole(address account) external onlyRole(DEFAULT_ADMIN_ROLE) {
        _revokeRole(UNPAUSER_ROLE, account);
    }

    /**
     * @dev Transfers UNPAUSER_ROLE to multisig wallet (typically after setup)
     * @param multisig The multisig wallet address
     */
    function transferUnpauserRoleToMultisig(address multisig) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(multisig != address(0), "Invalid multisig address");
        
        _revokeRole(UNPAUSER_ROLE, msg.sender);
        _grantRole(UNPAUSER_ROLE, multisig);
    }
}