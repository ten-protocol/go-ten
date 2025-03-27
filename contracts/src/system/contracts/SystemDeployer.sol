// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.0;

import "./TransactionPostProcessor.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {Fees} from "./Fees.sol";
import {MessageBus} from "../../cross_chain_messaging/common/MessageBus.sol";
import {PublicCallbacks} from "./PublicCallbacks.sol";

/**
 * @title SystemDeployer
 * @dev Contract that deploys the system contracts
 * 
 * TODO stefan to add docs
 */
contract SystemDeployer {
    event SystemContractDeployed(string name, address contractAddress);

    constructor(address eoaAdmin) {
       deployAnalyzer(eoaAdmin);
       address feesProxy = deployFees(eoaAdmin, 0);
       deployMessageBus(eoaAdmin, feesProxy);
       deployPublicCallbacks(eoaAdmin);
    }

    function deployAnalyzer(address eoaAdmin) internal {
        TransactionPostProcessor transactionsPostProcessor = new TransactionPostProcessor();
        bytes memory callData = abi.encodeWithSelector(transactionsPostProcessor.initialize.selector, eoaAdmin);
        address transactionsPostProcessorProxy = deployProxy(address(transactionsPostProcessor), eoaAdmin, callData);
        
        emit SystemContractDeployed("TransactionsPostProcessor", transactionsPostProcessorProxy);
    }

    function deployMessageBus(address eoaAdmin, address feesAddress) internal {
        MessageBus messageBus = new MessageBus();
        bytes memory callData = abi.encodeWithSelector(messageBus.initialize.selector, eoaAdmin, feesAddress);
        address messageBusProxy = deployProxy(address(messageBus), eoaAdmin, callData);

        emit SystemContractDeployed("MessageBus", messageBusProxy);
    }

    function deployPublicCallbacks(address eoaAdmin) internal {
        PublicCallbacks publicCallbacks = new PublicCallbacks();
        bytes memory callData = abi.encodeWithSelector(publicCallbacks.initialize.selector);
        address publicCallbacksProxy = deployProxy(address(publicCallbacks), eoaAdmin, callData);

        emit SystemContractDeployed("PublicCallbacks", publicCallbacksProxy);
    }

    function deployFees(address eoaAdmin, uint256 initialMessageFeePerByte) internal returns (address) {
        Fees fees = new Fees();
        bytes memory callData = abi.encodeWithSelector(fees.initialize.selector, initialMessageFeePerByte, eoaAdmin);
        address feesProxy = deployProxy(address(fees), eoaAdmin, callData);

        emit SystemContractDeployed("Fees", feesProxy);
        return feesProxy;
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