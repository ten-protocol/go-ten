// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;



import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

import "../IBridge.sol";
import "../IBridgeSubordinate.sol";
import "../../messaging/messenger/CrossChainEnabledObscuro.sol";
import "./ObscuroERC20.sol";

contract ObscuroL2Bridge is IBridge, IBridgeSubordinate, CrossChainEnabledObscuro {

    event CreatedWrappedToken(address remoteAddress, address localAddress, string name, string symbol);

    mapping(address => ObscuroERC20) public wrappedTokens;
    mapping(address => address) public localToRemoteToken;
    mapping(address => address) public remoteToLocalToken;

    address remoteBridgeAddress;

    constructor(address messenger, address remoteBridge) 
    CrossChainEnabledObscuro(messenger) 
    {
        remoteBridgeAddress = remoteBridge;
    }

    function createWrappedToken(address crossChainAddress, string calldata name, string calldata symbol) 
    external 
    onlyCrossChainSender(remoteBridgeAddress) 
    {
        ObscuroERC20 newToken = new ObscuroERC20(name, symbol);
        address localAddress = address(newToken);

        wrappedTokens[localAddress] = newToken;
        localToRemoteToken[localAddress] = crossChainAddress;
        remoteToLocalToken[crossChainAddress]= localAddress;

        emit CreatedWrappedToken(crossChainAddress, localAddress, name, symbol);
    }

    function hasTokenMapping(address wrappedToken) public view returns (bool) {
        return address(wrappedTokens[wrappedToken]) != address(0x0);
    }

    function sendNative(address receiver) external payable {
        require(msg.value > 0, "Nothing sent.");
        require(hasTokenMapping(address(0x0)), "No mapping for token.");

        bytes memory data = abi.encodeWithSelector(IBridge.receiveAssets.selector, localToRemoteToken[address(0x0)], msg.value, receiver);
        queueMessage(remoteBridgeAddress, data, uint32(Topics.TRANSFER), 0, 0);
    }

    function sendERC20(address asset, uint256 amount, address receiver) external {
        require(hasTokenMapping(asset), "No mapping for token.");
       
        ObscuroERC20 token = wrappedTokens[asset];
        token.burnFor(msg.sender, amount);

        bytes memory data = abi.encodeWithSelector(IBridge.receiveAssets.selector, localToRemoteToken[asset], amount, receiver);
        queueMessage(remoteBridgeAddress, data, uint32(Topics.TRANSFER), 0, 0);
    }

    function receiveAssets(address asset, uint256 amount, address receiver) 
    external 
    onlyCrossChainSender(remoteBridgeAddress) {

        address localAddress = remoteToLocalToken[asset];

        ObscuroERC20 token = wrappedTokens[localAddress];
        require(address(token) != address(0x0), "Receiving assets for unknown wrapped token!");

        token.issueFor(receiver, amount);
    }

    fallback() external payable {revert("fallback() method unsupported");}

    receive() external payable {revert("Contract does not support receive()");}
}