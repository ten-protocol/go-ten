// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import {CrossChainEnabledTEN} from "../../../cross_chain_messaging/lib/CrossChainEnabledTEN.sol";
import "../../L2/interfaces/ITokenFactory.sol";
import "../../common/IBridge.sol";
import "../interfaces/ITenBridge.sol";
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuardTransient.sol";

// This is the Ethereum side of the Obscuro Bridge.
// End-users can interact with it to transfer ERC20 tokens and native eth to the Layer 2 Obscuro.
contract TenBridge is
    CrossChainEnabledTEN,
    IBridge,
    ITenBridge,
    AccessControlUpgradeable,
    ReentrancyGuardTransient
{
    event Withdrawal(address indexed receiver, address indexed asset, uint256 amount);
    // This is the role that is given to the address that represents a native currency
    bytes32 public constant NATIVE_TOKEN_ROLE = keccak256("NATIVE_TOKEN");

    // This is the role that is given to addresses which are ERC20 contract.
    // If we have assigned a role to a contract it is considered whitelisted.
    bytes32 public constant ERC20_TOKEN_ROLE = keccak256("ERC20_TOKEN");
    bytes32 public constant SUSPENDED_ERC20_ROLE = keccak256("SUSPENDED_ERC20_TOKEN");

    // This is the role of the address that can perform administrative changes
    // like adding or removing tokens.
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");

    address public remoteBridgeAddress;

    function initialize(address messenger, address owner) public initializer {
        require(messenger != address(0), "Messenger cannot be 0x0");
        require(owner != address(0), "Owner cannot be 0x0");

        CrossChainEnabledTEN.configure(messenger);
        __AccessControl_init();
        _grantRole(DEFAULT_ADMIN_ROLE, owner);
        _grantRole(DEFAULT_ADMIN_ROLE, address(this));
        _setRoleAdmin(ADMIN_ROLE, DEFAULT_ADMIN_ROLE);
        _grantRole(ADMIN_ROLE, owner);
        _grantRole(NATIVE_TOKEN_ROLE, address(0x0));
    }

    function promoteToAdmin(address newAdmin) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newAdmin != address(0), "New admin cannot be 0x0");
        grantRole(ADMIN_ROLE, newAdmin);
    }

    function whitelistToken(
        address asset,
        string calldata name,
        string calldata symbol
    ) external onlyRole(ADMIN_ROLE) {
        require(asset != address(0), "Asset cannot be 0x0");
        require(!hasRole(ERC20_TOKEN_ROLE, asset), "Token already whitelisted");
        _grantRole(ERC20_TOKEN_ROLE, asset);

        bytes memory data = abi.encodeWithSelector(
            ITokenFactory.onCreateTokenCommand.selector,
            asset,
            name,
            symbol
        );
        queueMessage(
            remoteBridgeAddress,
            data,
            uint32(Topics.MANAGEMENT),
            0,
            0,
            0
        );
    }

    function pauseToken(address asset) external onlyRole(ADMIN_ROLE) {
        require(hasRole(ERC20_TOKEN_ROLE, asset), "Token is not whitelisted");
        _grantRole(SUSPENDED_ERC20_ROLE, asset);
    }

    function unpauseToken(address asset) external onlyRole(ADMIN_ROLE) {
        _revokeRole(SUSPENDED_ERC20_ROLE, asset);
    }

    function setRemoteBridge(address bridge) external onlyRole(ADMIN_ROLE) {
        require(bridge != address(0), "Bridge cannot be 0x0");
        if (remoteBridgeAddress != address(0)) {
            revert("Remote bridge address already set.");
        }
        remoteBridgeAddress = bridge;
    }

    // This cross chain message is specialized and will result in automatic increase
    // of balance on the other side.
    function sendNative(address receiver) external payable override {
        require(msg.value > 0, "Empty transfer.");
        bytes memory data = abi.encode(ValueTransfer(msg.value, receiver));
        publishRawMessage(data, uint32(Topics.VALUE), 0, 0); // No fee l1 to l2.
    }

    function sendERC20(
        address asset,
        uint256 amount,
        address receiver
    ) external payable override {
        require(!hasRole(SUSPENDED_ERC20_ROLE, asset), "Token is paused.");
        require(amount > 0, "Attempting empty transfer.");
        require(
            hasRole(ERC20_TOKEN_ROLE, asset),
            "This address has not been given a type and is thus considered not whitelisted."
        );

        // ensures the token is correctly transferred to the contract - tx reverts on failure
        SafeERC20.safeTransferFrom(
            IERC20(asset),
            msg.sender,
            address(this),
            amount
        );

        bytes memory data = abi.encodeWithSelector(
            IBridge.receiveAssets.selector,
            asset,
            amount,
            receiver
        );
        queueMessage(remoteBridgeAddress, data, uint32(Topics.TRANSFER), 0, 0, 0);
    }

    function receiveAssets(
        address asset,
        uint256 amount,
        address receiver
    ) external override onlyCrossChainSender(remoteBridgeAddress) nonReentrant {
        if (hasRole(ERC20_TOKEN_ROLE, asset)) {
            _receiveTokens(asset, amount, receiver);
        } else if (hasRole(NATIVE_TOKEN_ROLE, asset)) {
            _receiveNative(receiver, amount);
        } else {
            revert("Attempting to withdraw unknown asset.");
        }
    }

    function _receiveTokens(
        address asset,
        uint256 amount,
        address receiver
    ) private {
        SafeERC20.safeTransfer(IERC20(asset), receiver, amount);
        emit Withdrawal(receiver, asset, amount);
    }

    function _receiveNative(address receiver, uint256 amount) private {
        (bool sent, ) = receiver.call{value: amount}("");
        require(sent, "Failed to send Ether");
        emit Withdrawal(receiver, address(0), amount);
    }
}
