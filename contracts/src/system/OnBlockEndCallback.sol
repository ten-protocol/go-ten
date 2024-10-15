// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./Transaction.sol";


// OnBlockEndCallback is the interface that a contract needs to implement in order to support
// being called from the system transaction analyzer contract. Note that contracts are added as a callback
// with a manual authorization flow that whitelists them.  
interface OnBlockEndCallback {
    function onBlockEnd(Structs.Transaction[] calldata transactions) external;
}