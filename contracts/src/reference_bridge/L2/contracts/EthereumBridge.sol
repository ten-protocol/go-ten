// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "../../../common/UnrenouncableOwnable2Step.sol";
import "../../../cross_chain_messaging/lib/CrossChainEnabledTEN.sol";

import "../../common/IBridge.sol";
import "../interfaces/ITokenFactory.sol";
import "./WrappedERC20.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

contract EthereumBridge is
    IBridge,
    ITokenFactory,
    CrossChainEnabledTEN,
    Initializable,
    UnrenouncableOwnable2Step,
    PausableUpgradeable
{
    event CreatedWrappedToken(
        address remoteAddress,
        address localAddress,
        string name,
        string symbol
    );

    mapping(address localToken => WrappedERC20 wrappedToken) public wrappedTokens;
    mapping(address localToken => address remoteToken) public localToRemoteToken;
    mapping(address remoteToken => address localToken) public remoteToLocalToken;

    address public remoteBridgeAddress;

    function initialize(
        address messenger,
        address remoteBridge
    ) public initializer {
        require(messenger != address(0), "Messenger cannot be 0x0");
        require(remoteBridge != address(0), "Remote bridge cannot be 0x0");
        __Pausable_init();
        CrossChainEnabledTEN.configure(messenger);
        remoteBridgeAddress = remoteBridge;
    }

    function onCreateTokenCommand(
        address crossChainAddress,
        string calldata name,
        string calldata symbol
    ) external onlyCrossChainSender(remoteBridgeAddress) {
        WrappedERC20 newToken = new WrappedERC20(name, symbol);
        address localAddress = address(newToken);

        wrappedTokens[localAddress] = newToken;
        localToRemoteToken[localAddress] = crossChainAddress;
        remoteToLocalToken[crossChainAddress] = localAddress;

        emit CreatedWrappedToken(crossChainAddress, localAddress, name, symbol);
    }

    function hasTokenMapping(address wrappedToken) public view returns (bool) {
        return address(wrappedTokens[wrappedToken]) != address(0x0);
    }

    function erc20Fee() public view returns (uint256) {
        return _messageBus().getPublishFee();
    }

    function valueTransferFee() public view returns (uint256) {
        return _messageBus().getPublishFee();
    }

    function sendNative(address receiver) external payable whenNotPaused {
        require(msg.value > 0, "Nothing sent.");
        require(msg.value >= _messageBus().getPublishFee(), "Insufficient funds to publish value transfer");
        uint256 fee = _messageBus().getPublishFee();
        uint256 amount = msg.value - fee;
        bytes memory data = abi.encodeWithSelector(
            IBridge.receiveAssets.selector,
            address(0), 
            amount, 
            receiver
        );
        queueMessage(remoteBridgeAddress, data, uint32(Topics.VALUE), 0, 0, fee);
    }

    function sendERC20(
        address asset,
        uint256 amount,
        address receiver
    ) external payable whenNotPaused {
        require(hasTokenMapping(asset), "No mapping for token.");
        require(amount > 0, "Amount must be greater than 0.");

        WrappedERC20 token = wrappedTokens[asset];
        token.burnFor(msg.sender, amount);

        bytes memory data = abi.encodeWithSelector(
            IBridge.receiveAssets.selector,
            localToRemoteToken[asset],
            amount,
            receiver
        );

        require(msg.value >= _messageBus().getPublishFee(), "Insufficient funds to publish message");
        queueMessage(remoteBridgeAddress, data, uint32(Topics.TRANSFER), 0, 0, msg.value);
    }

    function receiveAssets(
        address asset,
        uint256 amount,
        address receiver
    ) external onlyCrossChainSender(remoteBridgeAddress) whenNotPaused {
        address localAddress = remoteToLocalToken[asset];

        WrappedERC20 token = wrappedTokens[localAddress];
        require(
            address(token) != address(0x0),
            "Receiving assets for unknown wrapped token!"
        );

        token.issueFor(receiver, amount);
    }

    fallback() external payable {
        revert("fallback() method unsupported");
    }

    receive() external payable {
        revert("Contract does not support receive()");
    }

    function pause() external onlyOwner {
        _pause();
    }

    function unpause() external onlyOwner {
        _unpause();
    }
}
