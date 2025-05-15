// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.28;

import "../../lib/Transaction.sol";
import "../interfaces/IOnBlockEndCallback.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";

/**
 * @title TransactionPostProcessor
 * @dev Contract that processes transactions after they are converted
 * 
 * TODO stefan to add docs
 */
contract TransactionPostProcessor is Initializable, AccessControl{
    using Structs for Structs.Transaction;

    bytes32 public constant EOA_ADMIN_ROLE = keccak256("EOA_ADMIN_ROLE");

    event TransactionsConverted(uint256 transactionsLength);

    struct Receipt {
        uint8 _type;
        bytes postState;
        uint64 Status;        
    }

    modifier onlySelf() {
        address maskedSelf = address(uint160(address(this)) - 1);
        require(msg.sender == maskedSelf, "Not self");
        _;
    }

    IOnBlockEndCallback[] public onBlockEndListeners;

    function initialize(address eoaAdmin) public initializer {
        _grantRole(DEFAULT_ADMIN_ROLE, eoaAdmin);
        _grantRole(EOA_ADMIN_ROLE, eoaAdmin);
    }

    function addOnBlockEndCallback(address callbackAddress) external onlyRole(EOA_ADMIN_ROLE) {
        require(callbackAddress != address(0), "Invalid callback address");
        require(callbackAddress.code.length > 0, "Callback address must be a contract");
        onBlockEndListeners.push(IOnBlockEndCallback(callbackAddress));
    }

    function onBlock(Structs.Transaction[] calldata transactions) external onlySelf {
        if (transactions.length == 0) {
            revert("No transactions to convert");
        }
        
//        emit TransactionsConverted(transactions.length);
        
        for (uint256 i = 0; i < onBlockEndListeners.length; ++i) {
            IOnBlockEndCallback callback = onBlockEndListeners[i];
            callback.onBlockEnd(transactions);
        }
    }
}