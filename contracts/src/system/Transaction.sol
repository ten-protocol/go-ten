// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";



library Structs {
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
}