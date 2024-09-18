// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";


library TransactionDecoder {
    using MessageHashUtils for bytes32;
    using ECDSA for bytes32;

    struct Transaction {
        uint8 txType;
        uint256 chainId;
        uint256 nonce;
        uint256 gasPrice;
        uint256 gasLimit;
        address to;
        uint256 value;
        bytes data;
        uint8 v;
        bytes32 r;
        bytes32 s;
        uint256 maxPriorityFeePerGas;
        uint256 maxFeePerGas;
        address[] accessList;
    }

    function recoverSender(Transaction calldata txData) external pure returns (address sender) {
        // Step 1: Hash the transaction data excluding the signature fields (v, r, s)
        bytes32 dataHash = keccak256(
            abi.encode(
                txData.txType,
                txData.chainId,
                txData.nonce,
                txData.gasPrice,
                txData.gasLimit,
                txData.to,
                txData.value,
                txData.data,
                txData.maxPriorityFeePerGas,
                txData.maxFeePerGas,
                txData.accessList
            )
        );

        // Step 2: Prefix the hash with the standard Ethereum message prefix
        bytes32 ethSignedHash = dataHash.toEthSignedMessageHash();

        // Step 3: Recover the address using the signature parameters
        sender = ECDSA.recover(ethSignedHash, txData.v, txData.r, txData.s);

        return sender;
    }


    function decode(bytes calldata rawTransaction) internal pure returns (Transaction memory) {
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

    function decodeLegacyTransaction(bytes calldata rawTransaction) private pure returns (Transaction memory) {
        Transaction memory _tx;
        _tx.txType = 0;

        (_tx.nonce, _tx.gasPrice, _tx.gasLimit, _tx.to, _tx.value, _tx.data, _tx.v, _tx.r, _tx.s) = 
            abi.decode(rawTransaction, (uint256, uint256, uint256, address, uint256, bytes, uint8, bytes32, bytes32));

        return _tx;
    }

    function decodeEIP2930Transaction(bytes calldata rawTransaction) private pure returns (Transaction memory) {
        Transaction memory _tx;
        _tx.txType = 1;

        (_tx.chainId, _tx.nonce, _tx.gasPrice, _tx.gasLimit, _tx.to, _tx.value, _tx.data, _tx.accessList, _tx.v, _tx.r, _tx.s) = 
            abi.decode(rawTransaction[1:], (uint256, uint256, uint256, uint256, address, uint256, bytes, address[], uint8, bytes32, bytes32));

        return _tx;
    }

    function decodeEIP1559Transaction(bytes calldata rawTransaction) private pure returns (Transaction memory) {
        Transaction memory _tx;
        _tx.txType = 2;

        (_tx.chainId, _tx.nonce, _tx.maxPriorityFeePerGas, _tx.maxFeePerGas, _tx.gasLimit, _tx.to, _tx.value, _tx.data, _tx.accessList, _tx.v, _tx.r, _tx.s) = 
            abi.decode(rawTransaction[1:], (uint256, uint256, uint256, uint256, uint256, address, uint256, bytes, address[], uint8, bytes32, bytes32));

        return _tx;
    }

    function getTo(Transaction memory _tx) internal pure returns (address) {
        return _tx.to;
    }

    function getValue(Transaction memory _tx) internal pure returns (uint256) {
        return _tx.value;
    }

    function getData(Transaction memory _tx) internal pure returns (bytes memory) {
        return _tx.data;
    }

    function getNonce(Transaction memory _tx) internal pure returns (uint256) {
        return _tx.nonce;
    }

    function getGasPrice(Transaction memory _tx) internal pure returns (uint256) {
        if (_tx.txType == 2) {
            return _tx.maxFeePerGas;
        }
        return _tx.gasPrice;
    }
}