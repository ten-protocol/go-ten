// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import "../IBridge.sol";
import "../../messaging/IMessageBus.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

contract ObscuroBridge is IBridge {

    IMessageBus messageBus;
    uint32 nonce;

    constructor(address busAddress) {
        messageBus = IMessageBus(busAddress);
        nonce = 0;
    }

    modifier onlyMessageBus() {
        require(msg.sender == address(messageBus), "Function can only be called by the message bus!");
        _;
    }

    function sendNative(address target) override external payable {
        require(msg.value > 0, "Empty transfer.");
        queueTransfer(TransferMessage(address(0x0), msg.value, target));
    }

    function sendAssets(address asset, uint256 amount, address receiver) override external {
        require(amount > 0, "Attempting empty transfer.");
        // Check the whitelist to include the ERC20 token.

        SafeERC20.safeTransferFrom(IERC20(asset), msg.sender, address(this), amount);

        queueTransfer(TransferMessage(asset, amount, receiver));
    }

    function receiveAssets(address asset, uint256 amount, address receiver) override external onlyMessageBus {
        SafeERC20.safeTransfer(IERC20(asset), receiver, amount);
        // emit event
    }

    function _receiveTokens(address asset, uint256 amount, address receiver) private {

    }

    function _receiveNative(uint256 amount, address receiver) private {
        
    }

    function queueTransfer(TransferMessage memory message) private {
        messageBus.publishMessage(nonce++, uint32(Topics.TRANSFER), abi.encode(message), 0);
    }
}