// SPDX-License-Identifier: Apache 2

pragma solidity >=0.7.0 <0.9.0;

import {TenBridge} from "../reference_bridge/L1/contracts/TenBridge.sol";

// THIS VERSION IS FOR TESTNETS ONLY.
// It includes a function to recover native funds from the bridge contract so we can preserve sepolia eth.
contract TenBridgeTestnet is TenBridge
{
    // This function is used to recover all native funds from the bridge contract - this is only available on testnets
    function recoverTestnetFunds(address receiver) external onlyRole(ADMIN_ROLE) {
        require(receiver != address(0), "Receiver cannot be 0x0");
        require( // this check ensures this could not be used on mainnet if it were deployed there accidentally
            block.chainid == 11155111 || block.chainid == 1337,
            "Recovery only allowed on approved testnets"
        );

        uint256 balance = address(this).balance;
        if (balance > 0) {
            (bool sent, ) = receiver.call{value: balance}("");
            require(sent, "Failed to recover Ether");
            emit Withdrawal(receiver, address(0), balance);
        }
    }
}
