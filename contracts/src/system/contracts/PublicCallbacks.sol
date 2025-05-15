// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

/**
 * @title PublicCallbacks
 * @dev Contract that allows to register callbacks that can be executed by the system
 * 
 * TODO stefan to add docs
 */
contract PublicCallbacks is Initializable {

    modifier onlySelf() {
        address maskedSelf = address(uint160(address(this)) - 1);
        require(msg.sender == maskedSelf, "Not self");
        _;
    }


    constructor() {
        _disableInitializers();
    }

    struct Callback {
        address target;
        bytes data;
        uint256 value;
        uint256 baseFee;
    }

    mapping(uint256 callbackId => Callback callback) public callbacks;
    uint256 private nextCallbackId;
    uint256 private lastUnusedCallbackId;

    mapping(uint256 callbackId => uint256 blockNumber) public callbackBlockNumber;

    // This modifier prevents using the callback in the same block it was registered (before the automation has a chance to do it)
    // this ensures that one can't commit and uncommit in the same transaction based on the outcome of reattempting.
    modifier canReattemptCallback(uint256 callbackId) {
        require(callbackBlockNumber[callbackId] < block.number, "Callback cannot be reattempted yet");
        _;
    }

    function initialize() external initializer {
    }

    function addCallback(address callback, bytes calldata data, uint256 value) internal returns (uint256 callbackId) {
        callbackId = nextCallbackId;
        callbacks[nextCallbackId++] = Callback({target: callback, data: data, value: value, baseFee: block.basefee});
        callbackBlockNumber[callbackId] = block.number;
    }

    function getCurrentCallbackToExecute() internal view returns (Callback memory, uint256) {
        return (callbacks[lastUnusedCallbackId], lastUnusedCallbackId);
    }

    function popCurrentCallback() internal {
        delete callbacks[lastUnusedCallbackId];
        delete callbackBlockNumber[lastUnusedCallbackId];
    }

    function moveToNextCallback() internal {
        lastUnusedCallbackId++;
    }

    function calculateGas(uint256 value) internal view returns (uint256) {
        return value / block.basefee;
    }

    // This function is callable from external dApps to register a callback.
    // The bytes passed in the param are the calldata for the call to be made
    // to msg.sender. 
    // todo: Consider making the callback function named in order to avoid
    // weird potential attacks if any? 
    function register(bytes calldata callback) external payable returns (uint256) { 
        require(msg.value > 0, "No value sent");
        require(calculateGas(msg.value) > 21000, "Gas too low compared to cost of call");
        // todo - add maximum value to limit
        return addCallback(msg.sender, callback, msg.value);
    }

    // reattempt a callback that failed to execute.
    // This is callable from external users and fully passes over the gas given to this call.
    function reattemptCallback(uint256 callbackId) external canReattemptCallback(callbackId) {
        Callback memory callback = callbacks[callbackId];
        (bool success, ) = callback.target.call(callback.data);
        require(success, "Callback execution failed");
        delete callbacks[callbackId];
        delete callbackBlockNumber[callbackId];
        // nothing to refund; the callback was already paid for during its failure
    }

    event CallbackExecuted(uint256 callbackId, uint256 gasBefore, uint256 gasAfter);

    // System level call. As it is called during a synthetic transaction that does not have gas limit, 
    // the contract enforces a custom limit based on the value stored for the callback.
    // It attempts to somewhat accurately refund.
    function executeNextCallbacks() external onlySelf {
        while (nextCallbackId != lastUnusedCallbackId) {
            executeNextCallback();
        }
    }

    function executeNextCallback() internal {
        if (nextCallbackId == lastUnusedCallbackId) {
            return; // todo: change to revert if possible
        }

        (Callback memory callback, uint256 callbackId) = getCurrentCallbackToExecute();
        uint256 baseFee = callback.baseFee;
        uint256 prepaidGas = callback.value / baseFee;
        uint256 gasBefore = gasleft();
        (bool success, ) = callback.target.call{gas: prepaidGas}(callback.data);
        uint256 gasAfter = gasleft();
    

        uint256 gasUsed = (gasBefore - gasAfter);
        uint256 gasRefundValue = 0;
        if (prepaidGas > gasUsed) {
            gasRefundValue = (prepaidGas - gasUsed) * baseFee;
        }
        uint256 paymentToCoinbase = callback.value - gasRefundValue;
        address target = callback.target;

        if (success) {  
            popCurrentCallback();
        }
        moveToNextCallback();

        internalRefund(gasRefundValue, target, callbackId);
        payForCallback(paymentToCoinbase);
    }

    function internalRefund(uint256 gasRefund, address to, uint256 callbackId) internal {
        // 22k is the max refund gas limit; 21k for a call and a bit for any accounting the contract might have.
        // ordinarily such accounting should be prepared for beforehand in the callback they pay for, but we give them a
        // slight buffer. 
        bytes memory data = abi.encodeWithSignature("handleRefund(uint256)", callbackId);
        (bool success, ) = to.call{value: gasRefund, gas: 45000}(data); 
        if (!success) {
            // if they dont accept the refund, we gift it to coinbase.
            payForCallback(gasRefund);
        }
    }

    function payForCallback(uint256 gasPayment) internal {
        if (gasPayment == 0) {
            return;
        }
        // We don't care about success, should always happen.
        // If not, contract is upgradable and we can recover.
        // solc-ignore-next-line unused-call-retval
        block.coinbase.call{value: gasPayment}("");
    }
}