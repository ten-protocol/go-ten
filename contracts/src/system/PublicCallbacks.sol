// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

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
        address callback;
        bytes data;
        uint256 value;
    }

    mapping(uint256 => Callback) public callbacks;
    uint256 private nextCallbackId;
    uint256 private lastUnusedCallbackId;

    function initialize() external initializer {
        nextCallbackId = 0;
        lastUnusedCallbackId = 0;
    }

    function addCallback(address callback, bytes calldata data, uint256 value) internal {
        callbacks[nextCallbackId++] = Callback({callback: callback, data: data, value: value});
    }

    function register(bytes calldata callback) external payable { 
        addCallback(msg.sender, callback, msg.value);
    }

    function executeNextCallback() external onlySelf {
        if (nextCallbackId > lastUnusedCallbackId) {
            return; // todo: change to revert if possible
        }

        uint256 callbackId = lastUnusedCallbackId++;
        require(callbackId < lastUnusedCallbackId, "Paranoia- todo: delete");
        Callback memory callback = callbacks[callbackId];
        (bool success, ) = callback.callback.call{value: callback.value}(callback.data);
        if (success) {
            delete callbacks[callbackId];
        }
    }
}