// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";

import "../IBridge.sol";
import "../ITokenFactory.sol";
import "../../messaging/messenger/CrossChainEnabledObscuro.sol";
import "../../common/WrappedERC20.sol";

contract EthereumBridge is
    IBridge,
    ITokenFactory,
    CrossChainEnabledObscuro
{
    event CreatedWrappedToken(
        address remoteAddress,
        address localAddress,
        string name,
        string symbol
    );

    mapping(address => WrappedERC20) public wrappedTokens;
    mapping(address => address) public localToRemoteToken;
    mapping(address => address) public remoteToLocalToken;

    address remoteBridgeAddress;

    function initialize(
        address messenger,
        address remoteBridge
    ) public initializer {
        CrossChainEnabledObscuro.configure(messenger);
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
        // receiveAssets selector (4 bytes) + address (20 bytes) + uint256 (32 bytes) + address (20 bytes) = 76 bytes
        uint256 dataLength = 76;
        return _messageBus().getMessageFee(dataLength);
    }

    function valueTransferFee() public view returns (uint256) {
        return _messageBus().getValueTransferFee();
    }

    function sendNative(address receiver) external payable {
        require(msg.value > 0, "Nothing sent.");
        require(msg.value >= _messageBus().getValueTransferFee(), "Insufficient funds to publish value transfer");
        _messageBus().sendValueToL2{value: msg.value}(receiver, msg.value);
    }

    function sendERC20(
        address asset,
        uint256 amount,
        address receiver
    ) external payable {
        require(hasTokenMapping(asset), "No mapping for token.");

        WrappedERC20 token = wrappedTokens[asset];
        token.burnFor(msg.sender, amount);

        bytes memory data = abi.encodeWithSelector(
            IBridge.receiveAssets.selector,
            localToRemoteToken[asset],
            amount,
            receiver
        );

        require(msg.value >= _messageBus().getMessageFee(data.length), "Insufficient funds to publish message");
        queueMessage(remoteBridgeAddress, data, uint32(Topics.TRANSFER), 0, 0, msg.value);
    }

    function receiveAssets(
        address asset,
        uint256 amount,
        address receiver
    ) external onlyCrossChainSender(remoteBridgeAddress) {
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
}
