// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "libs/openzeppelin/contracts/token/ERC20/ERC20.sol";
//import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

// todo - can't import from the /contracts folder
// implement this interface if you want to configure the visibility rules of your smart contract
// the TEN platform will interpret this information
interface ContractTransparencyConfig {
    // configuration per event log type
    struct EventLogConfig {
        bytes eventSignature;
        bool isPublic;  // everyone can see and query for this event
        bool topic1CanView;    // If the event is private, and this is true, it means that the address from topic1 is an EOA that can view this event
        bool topic2CanView;    // same
        bool topic3CanView;    // same
        bool visibleToSender; // if true, the tx signer will see this event. Default false
    }

    struct VisibilityConfig {
        bool isTransparent; // If true - the internal state via getStorageAt will be accessible to everyone. All events will be public. Default false
        EventLogConfig[] eventLogConfigs;  // mapping from event signature to visibility configs per event
    }

    // keep the logic independent of the environment
    // max gas: 1 Million
    function visibilityRules() external pure returns (VisibilityConfig memory);
}

interface Structs {
    struct CrossChainMessage {
        address sender;
        uint64 sequence;
        uint32 nonce;
        bytes topic;
        bytes payload;
        uint8 consistencyLevel;
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

// This is an implementation of a canonical ERC20 as used in the Obscuro network
// where access to data has to be restricted.
contract ObsERC20 is ERC20, ContractTransparencyConfig {

    address bridge = 0xdeB34A740ECa1eC42C8b8204CBEC0bA34FDD27f3;

    IMessageBus bus;

    enum Topics{
        MINT,
        TRANSFER
    }

    struct AssetTransferMessage {
        address sender;
        address receiver;
        uint256 amount;
    }

    constructor(
        string memory name,
        string memory symbol,
        uint256 initialSupply,
        address busAddress
    )  ERC20(name, symbol) {
        _mint(msg.sender, initialSupply);
        bus = IMessageBus(busAddress);
    }

    function visibilityRules() public override pure returns (VisibilityConfig memory){
        EventLogConfig[]  memory configs = new EventLogConfig[](1);
        // erc20 transfer
        configs[0] = EventLogConfig(hex"ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", false, true, true, false, false);
        return VisibilityConfig(false, configs);
    }

    function _beforeTokenTransfer(address from, address to, uint256 amount)
    internal virtual override {
        //Only deposit messages.
        if (address(bus) == address(0x0)) {
            return;
        }

        if (to == bridge) {
            AssetTransferMessage memory message = AssetTransferMessage(from, to, amount);
            uint64 sequence = bus.publishMessage(uint32(block.number), uint32(Topics.TRANSFER), abi.encode(message), 0);
        }
    }

    function balanceOf(address account) public view virtual override returns (uint256) {
        // 1. Human owner of an account asking for the balance.
        // 2. Human owner of an account interacting with a smart contract which in turn asks for the balance of the original asker.
        // Note: In case the requester spoofs the "from" of the call, they will not be able to read
        //  the result since it will be returned encrypted with the viewing key of the declared "from".
        if (tx.origin == account) {
            return super.balanceOf(account);
        }

        // 3. Contract asking for its own balance.
        if (msg.sender == account) {
            return super.balanceOf(account);
        }

        revert("Not allowed to read the balance");
    }

    function allowance(address owner, address spender) public view virtual override returns (uint256) {
        // 1. Human owner of an account asking for the allowance it has empowered someone to spend,
        // or Human owner of an account asking for how much it is allowed to spend by someone else.
        // 2. Human owner of an account interacting with a smart contract which in turn asks for the above
        // Note: In case the requester spoofs the "from" of the call, they will not be able to read
        //  the result since it will be returned encrypted with the viewing key of the declared "from".
        if (tx.origin == owner || tx.origin == spender) {
            return super.allowance(owner, spender);
        }

        // 3. Contract asking how much it is empowered to spend, or how much it has empowered someone else to spend.
        if (msg.sender == owner || msg.sender == spender) {
            return super.allowance(owner, spender);
        }

        revert("Not allowed to read the allowance");
    }
}