// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "../messaging/IMessageBus.sol";

// This is an implementation of a canonical ERC20 as used in the Obscuro network
// where access to data has to be restricted.
contract ObsERC20 is ERC20 {
    address bridge = 0xdeB34A740ECa1eC42C8b8204CBEC0bA34FDD27f3;

    constructor(
        string memory name,
        string memory symbol
    ) ERC20(name, symbol) {
    }

    function balanceOf(
        address account
    ) public view virtual override returns (uint256) {
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

    function allowance(
        address owner,
        address spender
    ) public view virtual override returns (uint256) {
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
