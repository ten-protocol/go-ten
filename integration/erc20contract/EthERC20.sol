// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "libs/openzeppelin/contracts/token/ERC20/ERC20.sol";
//import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

interface Structs {
    struct CrossChainMessage {
        address sender;
        uint64  sequence;
        uint32  nonce;
        bytes   topic;
        bytes   payload;
        uint8   consistencyLevel;
    }
}

interface IMessageBus {
    function publishMessage(
        uint32 nonce,
        uint32 topic,
        bytes calldata payload, 
        uint8 consistencyLevel
    ) external payable returns (uint64 sequence);
}

interface ICrossChain {
    function messageBus() external view returns (IMessageBus);
}

contract EthERC20 is ERC20 {

    event Something(string hm);

    IMessageBus bus;

    address target;

    enum Topics{ MINT, TRANSFER }
    struct AssetTransferMessage {
        address sender;
        address receiver;
        uint256 amount;
    }

    constructor(
        string memory name,
        string memory symbol,
        uint256 initialSupply,
        address crossChainContract
    )  ERC20(name, symbol) {

        if (crossChainContract != address(0x0)) {
            bus = ICrossChain(crossChainContract).messageBus();
        }

        target = crossChainContract;
        _mint(msg.sender, initialSupply);
    }

    function _beforeTokenTransfer(address from, address to, uint256 amount)
    internal virtual override {
        //Only if message bus is configured.
        if (address(bus) == address(0x0)) {
            return;
        }

        //Only deposit messages.
        if (to == target) { 
            AssetTransferMessage memory message = AssetTransferMessage(from, to, amount);
            bus.publishMessage(uint32(block.number), uint32(Topics.TRANSFER), abi.encode(message), 0);
        }
    }
}