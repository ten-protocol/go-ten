// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IPublicCallbacks} from "../system/PublicCallbacks.sol";

contract PublicCallbacksTest {
    IPublicCallbacks public callbacks;

    constructor(address _callbacks) payable {
        callbacks = IPublicCallbacks(_callbacks);
        lastCallSuccess = false;
        allCallbacksRan = false;
        testRegisterCallback();
    }

    bool lastCallSuccess = false;
    bool allCallbacksRan = false;

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

    function handleAllCallbacksRan() external {
        allCallbacksRan = true;
    }

    mapping(uint256 => address) public callbackRefundees;
    mapping(address=>uint256) public pendingRefunds;

    // This is called by the system to refund the callbac
    function handleRefund(uint256 callbackId) external payable {
        // do nothing
        pendingRefunds[callbackRefundees[callbackId]] += msg.value;
        refundCalled++;
    }

    // Test function that registers a callback
    function testRegisterCallback() internal {
        // Encode the callback data - calling handleCallback()
        // Calculate expected gas based on value sent
        uint256 expectedGas = (msg.value/3) / block.basefee;
        bytes memory callbackData = abi.encodeWithSelector(this.handleCallback.selector, expectedGas);
        
        bytes memory callbackDataFail = abi.encodeWithSelector(this.handleCallbackFail.selector);
        bytes memory callbackDataAllCallbacksRan = abi.encodeWithSelector(this.handleAllCallbacksRan.selector);

        // Register the callback, forwarding any value sent to this call
        uint256 id = callbacks.register{value: msg.value/3}(callbackData);
        callbackRefundees[id] = msg.sender;
        id = callbacks.register{value: msg.value/3}(callbackDataFail);
        callbackRefundees[id] = msg.sender;
        id = callbacks.register{value: msg.value/3}(callbackDataAllCallbacksRan);
        callbackRefundees[id] = msg.sender;
    }

    uint256 refundCalled = 0;
    function isLastCallSuccess() external view returns (bool) {
        return allCallbacksRan && refundCalled == 3;
    }
}
