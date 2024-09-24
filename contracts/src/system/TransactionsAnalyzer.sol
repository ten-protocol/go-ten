// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import "./TransactionDecoder.sol";
import "./OnBlockEndCallback.sol";

//TODO: @PR Review - Pick appropriate name
contract TransactionsAnalyzer is Initializable, AccessControl{
    bytes32 public constant EOA_ADMIN_ROLE = keccak256("EOA_ADMIN_ROLE");
    bytes32 public constant HOOK_CALLER_ROLE = keccak256("HOOK_CALLER_ROLE");

    struct Receipt {
        uint8 _type;
        bytes postState;
        uint64 Status;        
        /*
        CumulativeGasUsed uint64
        Bloom             Bloom 
        Logs              []*Log

        TxHash            common.Hash    
        ContractAddress   common.Address 
        GasUsed           uint64         
        EffectiveGasPrice *big.In
        BlobGasUsed       uint64  
        BlobGasPrice      *big.Int

        BlockHash        common.Ha
        BlockNumber      *big.Int 
        TransactionIndex uint  
        */  
    }

    struct BlockTransactions {
        bytes[] transactions;
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

    function onBlock(BlockTransactions calldata _block) public onlyRole(HOOK_CALLER_ROLE) {
        if (_block.transactions.length == 0) {
            return;
        }

        Structs.Transaction[] memory transactions = new Structs.Transaction[](_block.transactions.length);
        
        for (uint256 i = 0; i < _block.transactions.length; ++i) {
            transactions[i] = (TransactionDecoder.decode(_block.transactions[i]));            
        }


        for (uint256 i = 0; i < onBlockEndListeners.length; ++i) {
            OnBlockEndCallback callback = onBlockEndListeners[i];
            callback.onBlockEnd(transactions);
        }
    }
}