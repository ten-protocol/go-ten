// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

interface IPublicCallbacks {
    function register(bytes calldata callback) external payable;
    function reattemptCallback(uint256 callbackId) external;
}


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

    mapping(uint256 => Callback) public callbacks;
    uint256 private nextCallbackId;
    uint256 private lastUnusedCallbackId;

    function initialize() external initializer {
        nextCallbackId = 0;
        lastUnusedCallbackId = 0;
    }

    function addCallback(address callback, bytes calldata data, uint256 value) internal {
        callbacks[nextCallbackId++] = Callback({target: callback, data: data, value: value, baseFee: block.basefee});
    }

    function calculateGas(uint256 value) internal view returns (uint256) {
        return value / block.basefee;
    }

    function register(bytes calldata callback) external payable { 
        require(msg.value > 0, "No value sent");
        require(calculateGas(msg.value) > 21000, "Gas too low compared to cost of call");
        addCallback(msg.sender, callback, msg.value);
    }

    // reattempt a callback that failed to execute.
    // This is callable from external users and fully passes over the gas given to this call.
    function reattemptCallback(uint256 callbackId) external {
        Callback memory callback = callbacks[callbackId];
        (bool success, ) = callback.target.call(callback.data);
        require(success, "Callback execution failed");
        delete callbacks[callbackId];
        // nothing to refund; the callback was already paid for during its failure
    }

    event CallbackExecuted(uint256 callbackId, uint256 gasBefore, uint256 gasAfter);

    // System level call. As it is called during a synthetic transaction that does not have gas limit, 
    // the contract enforces a custom limit based on the value stored for the callback.
    // It attempts to somewhat accurately refund.
    function executeNextCallback() external onlySelf {
        if (nextCallbackId == lastUnusedCallbackId) {
            return; // todo: change to revert if possible
        }

        uint256 callbackId = lastUnusedCallbackId;
        lastUnusedCallbackId++;
        require(callbackId < lastUnusedCallbackId, "Paranoia- todo: delete");
        Callback storage callback = callbacks[callbackId];
        uint256 baseFee = callback.baseFee;
        uint256 gas = callback.value / baseFee;
        uint256 gasBefore = gasleft();
        (bool success, ) = callback.target.call{gas: gas}(callback.data);
        if (success) {
            delete callbacks[callbackId];
        }
        uint256 gasAfter = gasleft();
        emit CallbackExecuted(callbackId, gasBefore, gasAfter);
       // uint256 gasRefund = (gasBefore - gasAfter);
       // callback.value = callback.value - gasRefund;

        //internalRefund(gasRefund, callback.target);
        payForCallback(callback.value);
    }

    function internalRefund(uint256 gasRefund, address to) internal {
        // 22k is the max refund gas limit; 21k for a call and a bit for any accounting the contract might have.
        // ordinarily such accounting should be prepared for beforehand in the callback they pay for, but we give them a
        // slight buffer. 
        (bool success, ) = to.call{value: gasRefund, gas: 22000}(""); 
        if (!success) {
            block.coinbase.transfer(gasRefund); // if they dont accept the refund, we gift it to coinbase.
        }
    }

    function payForCallback(uint256 gasPayment) internal {
        block.coinbase.transfer(gasPayment);
    }
}