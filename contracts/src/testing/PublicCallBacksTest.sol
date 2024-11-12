// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IPublicCallbacks} from "../system/PublicCallbacks.sol";

contract PublicCallbacksTest {
    IPublicCallbacks public callbacks;

    constructor(address _callbacks) payable {
        callbacks = IPublicCallbacks(_callbacks);
        lastCallSuccess = true;
        testRegisterCallback();
    }

    bool lastCallSuccess = true;

    // This function will be called back by the system
    function handleCallback(uint256 expectedGas) external {
        uint256 gasGiven = gasleft();
        if (gasGiven > expectedGas - 22000) { //call + 1000 for calldata (which overshoots greatly)
            lastCallSuccess = false;
        }
        // Handle the callback here
        // For testing we'll just allow it to succeed
        
    }

    // Test function that registers a callback
    function testRegisterCallback() internal {
        require(lastCallSuccess, "Last call failed");
        // Encode the callback data - calling handleCallback()
        // Calculate expected gas based on value sent
        uint256 expectedGas = msg.value / block.basefee;
        bytes memory callbackData = abi.encodeWithSelector(this.handleCallback.selector, expectedGas);
        
        // Register the callback, forwarding any value sent to this call
        callbacks.register{value: msg.value}(callbackData);
    }

    function isLastCallSuccess() external view returns (bool) {
        return lastCallSuccess;
    }
}
