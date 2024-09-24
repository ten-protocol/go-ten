// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./Transaction.sol";


interface OnBlockEndCallback {
    function onBlockEnd(Structs.Transaction[] calldata transactions) external;
}