// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IPublicCallbacks} from "../system/PublicCallbacks.sol";

contract PublicCallbacksTest {
    IPublicCallbacks public callbacks;

    constructor(address _callbacks) payable {
        callbacks = IPublicCallbacks(_callbacks);
        lastCallSuccess = false;
        testRegisterCallback();
    }

    bool lastCallSuccess = false;

    // This function will be called back by the system
    function handleCallback(uint256 expectedGas) external {
        uint256 gasGiven = gasleft();
        if (gasGiven >= expectedGas - 2100) { //call + 1000 for calldata (which overshoots greatly)
            lastCallSuccess = true;
        }
        // Handle the callback here
        // For testing we'll just allow it to succeed
    }

    function handleCallbackFail() external {
        lastCallSuccess = true;
        require(false, "This is a test failure");
    }

    // Test function that registers a callback
    function testRegisterCallback() internal {
        // Encode the callback data - calling handleCallback()
        // Calculate expected gas based on value sent
        uint256 expectedGas = (msg.value/2) / block.basefee;
        bytes memory callbackData = abi.encodeWithSelector(this.handleCallback.selector, expectedGas);
        
        bytes memory callbackDataFail = abi.encodeWithSelector(this.handleCallbackFail.selector);

        // Register the callback, forwarding any value sent to this call
        callbacks.register{value: msg.value/2}(callbackData);
        callbacks.register{value: msg.value/2}(callbackDataFail);
    }

    function isLastCallSuccess() external view returns (bool) {
        return lastCallSuccess;
    }
}
