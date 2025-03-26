// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.0;

import "../../lib/Transaction.sol";
import "../interfaces/IOnBlockEndCallback.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

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

    IOnBlockEndCallback[] onBlockEndListeners;

    function initialize(address eoaAdmin) public initializer {
        _grantRole(DEFAULT_ADMIN_ROLE, eoaAdmin);
        _grantRole(EOA_ADMIN_ROLE, eoaAdmin);
    }

    function addOnBlockEndCallback(address callbackAddress) public {
        onBlockEndListeners.push(IOnBlockEndCallback(callbackAddress));
    }

    function onBlock(Structs.Transaction[] calldata transactions) public onlySelf {
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