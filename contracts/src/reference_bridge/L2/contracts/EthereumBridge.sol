// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "../../../common/UnrenouncableOwnable2Step.sol";
import "../../../cross_chain_messaging/lib/CrossChainEnabledTEN.sol";

import "../../common/IBridge.sol";
import "../interfaces/ITokenFactory.sol";
import "./WrappedERC20.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "../../../common/PausableWithRoles.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

// Minimal WETH interface used only for unwrapping
interface IWETH {
    function withdraw(uint256 wad) external;
}

contract EthereumBridge is
    IBridge,
    ITokenFactory,
    CrossChainEnabledTEN,
    Initializable,
    PausableWithRoles
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
    address public localWETH;

    function initialize(
        address messenger,
        address remoteBridge,
        address localWeth
    ) public initializer {
        require(messenger != address(0), "Messenger cannot be 0x0");
        require(remoteBridge != address(0), "Remote bridge cannot be 0x0");
        require(localWeth != address(0), "Local WETH cannot be 0x0");
        __PausableWithRoles_init(msg.sender);
        CrossChainEnabledTEN.configure(messenger);
        remoteBridgeAddress = remoteBridge;
        localWETH = localWeth;
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
        require(amount > 0, "Amount must be greater than 0.");

        if (asset == localWETH) {
            _bridgeWETHToL1(amount, receiver);
            return;
        }

        require(hasTokenMapping(asset), "No mapping for token.");
        WrappedERC20 token = wrappedTokens[asset];
        token.burnFor(msg.sender, amount);

        address remoteToken = localToRemoteToken[asset];
        require(msg.value >= _messageBus().getPublishFee(), "Insufficient funds to publish message");

        bytes memory data = abi.encodeWithSelector(
            IBridge.receiveAssets.selector,
            remoteToken,
            amount,
            receiver
        );
        queueMessage(remoteBridgeAddress, data, uint32(Topics.TRANSFER), 0, 0, msg.value);
    }

    function _bridgeWETHToL1(uint256 amount, address receiver) internal {
        // Collect WETH from the user and unwrap to native ETH
        SafeERC20.safeTransferFrom(IERC20(localWETH), msg.sender, address(this), amount);
        IWETH(localWETH).withdraw(amount);

        // Single publish fee covers both messages
        uint256 fee = _messageBus().getPublishFee();
        require(msg.value >= fee, "Insufficient funds for publish fee");

        // 1) Increase L1 bridge native balance by sending native to the remote bridge address
        this.sendNative{value: amount + fee}(remoteBridgeAddress);

        // 2) Instruct L1 to wrap/pay out to the receiver from its increased balance
        bytes memory nativeWrappedNotice = abi.encodeWithSelector(
            IBridge.receiveNativeWrapped.selector,
            receiver,
            amount
        );
        // No extra fee for the second message
        queueMessage(remoteBridgeAddress, nativeWrappedNotice, uint32(Topics.TRANSFER), 0, 0, 0);
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

    function receiveNativeWrapped(
        address receiver,
        uint256 amount
    ) external onlyCrossChainSender(remoteBridgeAddress) whenNotPaused {
        // L1 called sendNative (increasing our L2 native balance) and also queued this call
        // We wrap the received native amount into local WETH and transfer to receiver
        (bool success, ) = localWETH.call{value: amount}(abi.encodeWithSignature("deposit()"));
        require(success, "WETH deposit failed");
        IERC20(localWETH).transfer(receiver, amount);
    }

    fallback() external payable {
        revert("fallback() method unsupported");
    }

    receive() external payable {}
}
