// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {MessageBus} from "../messaging/MessageBus.sol";
import "./TransactionPostProcessor.sol";

contract SystemDeployer {
    event SystemContractDeployed(string name, address contractAddress);

    constructor(address eoaAdmin) {
       deployAnalyzer(eoaAdmin);
       deployMessageBus(eoaAdmin);
    }

    function deployAnalyzer(address eoaAdmin) internal {
        TransactionPostProcessor transactionsPostProcessor = new TransactionPostProcessor();
        bytes memory callData = abi.encodeWithSelector(transactionsPostProcessor.initialize.selector, eoaAdmin);
        address transactionsPostProcessorProxy = deployProxy(address(transactionsPostProcessor), eoaAdmin, callData);
        
        emit SystemContractDeployed("TransactionsPostProcessor", transactionsPostProcessorProxy);
    }

    function deployMessageBus(address eoaAdmin) internal {
        MessageBus messageBus = new MessageBus();
        bytes memory callData = abi.encodeWithSelector(messageBus.initialize.selector, eoaAdmin);
        address messageBusProxy = deployProxy(address(messageBus), eoaAdmin, callData);

        emit SystemContractDeployed("MessageBus", messageBusProxy);
    }

    function deployProxy(address _logic, address _admin, bytes memory _data) internal returns (address proxyAddress) {
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            _logic,   // Address of the logic contract
            _admin,   // Address of the admin (who can upgrade the proxy)
            _data     // Initialization data (can be empty if not required)
        );

        proxyAddress = address(proxy);  // Store proxy address for reference
    }
}