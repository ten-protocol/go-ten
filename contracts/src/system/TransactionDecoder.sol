// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./Transaction.sol";

library TransactionDecoder {
    using Structs for Structs.Transaction;

    function decode(bytes calldata rawTransaction) internal pure returns (Structs.Transaction memory) {
        require(rawTransaction.length > 0, "Empty transaction data");
        uint8 _txType = uint8(rawTransaction[0]);
        
        if (_txType == 0) {
            return decodeLegacyTransaction(rawTransaction);
        } else if (_txType == 1) {
            return decodeEIP2930Transaction(rawTransaction);
        } else if (_txType == 2) {
            return decodeEIP1559Transaction(rawTransaction);
        } else {
            revert("Unsupported transaction type");
        }
    }

    function decodeLegacyTransaction(bytes calldata rawTransaction) private pure returns (Structs.Transaction memory) {
        Structs.Transaction memory _tx;
        _tx.txType = 0;

        (_tx.nonce, _tx.gasPrice, _tx.gasLimit, _tx.to, _tx.value, _tx.data, _tx.v, _tx.r, _tx.s) = 
            abi.decode(rawTransaction, (uint256, uint256, uint256, address, uint256, bytes, uint8, bytes32, bytes32));

        return _tx;
    }

    function decodeEIP2930Transaction(bytes calldata rawTransaction) private pure returns (Structs.Transaction memory) {
        Structs.Transaction memory _tx;
        _tx.txType = 1;

        (_tx.chainId, _tx.nonce, _tx.gasPrice, _tx.gasLimit, _tx.to, _tx.value, _tx.data, _tx.accessList, _tx.v, _tx.r, _tx.s) = 
            abi.decode(rawTransaction[1:], (uint256, uint256, uint256, uint256, address, uint256, bytes, address[], uint8, bytes32, bytes32));

        return _tx;
    }

    function decodeEIP1559Transaction(bytes calldata rawTransaction) private pure returns (Structs.Transaction memory) {
        Structs.Transaction memory _tx;
        _tx.txType = 2;

        (_tx.chainId, _tx.nonce, _tx.maxPriorityFeePerGas, _tx.maxFeePerGas, _tx.gasLimit, _tx.to, _tx.value, _tx.data, _tx.accessList, _tx.v, _tx.r, _tx.s) = 
            abi.decode(rawTransaction[1:], (uint256, uint256, uint256, uint256, uint256, address, uint256, bytes, address[], uint8, bytes32, bytes32));

        return _tx;
    }
}