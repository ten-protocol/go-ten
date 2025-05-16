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

    event CallbackAdded(address callbackAddress);

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
        emit CallbackAdded(callbackAddress);
    }

    function removeOnBlockEndCallback(address callbackAddress) external onlyRole(EOA_ADMIN_ROLE) {
        require(callbackAddress != address(0), "Invalid callback address");
        
        uint256 length = onBlockEndListeners.length;
        for (uint256 i = 0; i < length; ++i) {
            if (address(onBlockEndListeners[i]) == callbackAddress) {
                // Move the last element to the position being deleted
                onBlockEndListeners[i] = onBlockEndListeners[length - 1];
                // Remove the last element
                onBlockEndListeners.pop();
                return;
            }
        }
        revert("Callback not found");
    }

    function onBlock(Structs.Transaction[] calldata transactions) external onlySelf {
        if (transactions.length == 0) {
            revert("No transactions to convert");
        }
        
//        emit TransactionsConverted(transactions.length);
        
        for (uint256 i = 0; i < onBlockEndListeners.length; ++i) {
            IOnBlockEndCallback callback = onBlockEndListeners[i];
            // All on block end callbacks are admin contracts, thus success is required.
            // All failures from such contracts are treated as general node failure and should revert this transaction.
            // This in turn will block the batch production. In practice no contract registered here should be able to revert
            // unless it detects some serious issue.
            callback.onBlockEnd(transactions); 
        }
    }
}