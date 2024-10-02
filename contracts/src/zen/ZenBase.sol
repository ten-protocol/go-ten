// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.0;

/*
// Import OpenZeppelin Contracts
import "../system/TransactionDecoder.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

*/

import "../system/OnBlockEndCallback.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";




interface ITransactionDecoder {
    function recoverSender(Structs.Transaction calldata txData) external view returns (address sender);
}

/**
 * @title ZenBase
 * @dev ERC20 Token with minting functionality.
 */
contract ZenBase is OnBlockEndCallback, ERC20, Ownable {
    using Structs for Structs.Transaction;

    event TransactionProcessed(address sender, uint256 amount);
    /**
     * @dev Constructor that gives msg.sender all of existing tokens.
     * You can customize the initial supply as needed.
     */
    constructor(address transactionAnalyzer, address transactionDecoder) ERC20("Zen", "ZEN") Ownable(msg.sender) {
        require(transactionAnalyzer != address(0), "Invalid transaction analyzer address");
        require(transactionDecoder != address(0), "Invalid transaction decoder address");
        _caller = transactionAnalyzer;
        _transactionDecoder = transactionDecoder;
    }
    

    address private _caller;
    address private _transactionDecoder;

    modifier onlyCaller() {
        require(msg.sender == _caller, "Caller: caller is not the designated address");
        _;
    }

    function onBlockEnd(Structs.Transaction[] calldata transactions) external onlyCaller {
        if (transactions.length == 0) {
            revert("No transactions to convert");
        }
        // Implement custom logic here
        for (uint256 i=0; i<transactions.length; i++) {
            // Process transactions
            _mint(transactions[i].from, 1);
            emit TransactionProcessed(transactions[i].from, 1);
        }
    }
}
