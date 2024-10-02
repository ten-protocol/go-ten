// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import "./OnBlockEndCallback.sol";
import "./Transaction.sol";

//TODO: @PR Review - Pick appropriate name
contract TransactionsAnalyzer is Initializable, AccessControl{
    using Structs for Structs.Transaction;

    bytes32 public constant EOA_ADMIN_ROLE = keccak256("EOA_ADMIN_ROLE");
    bytes32 public constant HOOK_CALLER_ROLE = keccak256("HOOK_CALLER_ROLE");

    event TransactionsConverted(uint256 transactionsLength);

    struct Receipt {
        uint8 _type;
        bytes postState;
        uint64 Status;        
    }

    OnBlockEndCallback[] onBlockEndListeners;

    function initialize(address eoaAdmin, address authorizedCaller) public initializer {
        _grantRole(DEFAULT_ADMIN_ROLE, eoaAdmin);
        _grantRole(EOA_ADMIN_ROLE, eoaAdmin);
        _grantRole(HOOK_CALLER_ROLE, authorizedCaller);
    }

    function addOnBlockEndCallback(address callbackAddress) public {
        onBlockEndListeners.push(OnBlockEndCallback(callbackAddress));
    }

    function onBlock(Structs.Transaction[] calldata transactions) public onlyRole(HOOK_CALLER_ROLE) {
        if (transactions.length == 0) {
            revert("No transactions to convert");
        }
        
        emit TransactionsConverted(transactions.length);
        
        for (uint256 i = 0; i < onBlockEndListeners.length; ++i) {
            OnBlockEndCallback callback = onBlockEndListeners[i];
            callback.onBlockEnd(transactions);
        }
    }
}