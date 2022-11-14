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
    }
}

interface IMessageBus {
    function publishMessage(
        uint32 nonce,
        uint32 topic,
        bytes calldata payload, 
        uint8 consistencyLevel
    ) external payable returns (uint64 sequence);

    function verifyMessageFinalized(Structs.CrossChainMessage calldata crossChainMessage) external view returns (bool);
    
    function getMessageTimeOfFinality(Structs.CrossChainMessage calldata crossChainMessage) external view returns (uint256);

    function submitOutOfNetworkMessage(Structs.CrossChainMessage calldata crossChainMessage, uint256 finalAfterTimestamp) external;

   /* function queryMessages(
        address      sender,
        bytes memory topic,
        uint256      fromIndex,
        uint256      toIndex
    ) external returns (bytes [] memory); */
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
        address l1MessageBus,
        address managementContract
    )  ERC20(name, symbol) {
        bus = IMessageBus(l1MessageBus);
        bus.publishMessage(uint32(block.number), uint32(Topics.MINT), abi.encodePacked(initialSupply), 0);
        _mint(msg.sender, initialSupply);
        target = l1MessageBus;
    }

    function _beforeTokenTransfer(address from, address to, uint256 amount)
    internal virtual override {

        AssetTransferMessage memory message = AssetTransferMessage(from, to, amount);
        uint64 sequence = bus.publishMessage(uint32(block.number), uint32(Topics.TRANSFER), abi.encode(message), 0);
        require(sequence == 1, "Sanity check fail");
    }
}