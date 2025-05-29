// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.28;

import "./TransactionPostProcessor.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {Fees} from "./Fees.sol";
import {CrossChainMessenger} from "../../cross_chain_messaging/common/CrossChainMessenger.sol";
import {EthereumBridge} from "../../reference_bridge/L2/contracts/EthereumBridge.sol";
import {MessageBus} from "../../cross_chain_messaging/common/MessageBus.sol";
import {PublicCallbacks} from "./PublicCallbacks.sol";
import {TenSystemCalls} from "./TenSystemCalls.sol";
/**
 * @title SystemDeployer
 * @dev Contract that deploys the system contracts
 * 
 * Auto executed contract on the L2 at the second batch, used to deploy the other system contracts.
 * The eoaAdmin is the owner of the proxies and can upgrade them. 
 * depends on the remoteBridgeAddress in order to configure the cross chain functionality.
 */
contract SystemDeployer {
    event SystemContractDeployed(string name, address contractAddress);

    constructor(address eoaAdmin, address remoteBridgeAddress) {
        require(eoaAdmin != address(0), "Invalid EOA admin address");
        require(remoteBridgeAddress != address(0), "Invalid remote bridge address");
        deployAnalyzer(eoaAdmin);
        address feesProxy = deployFees(eoaAdmin, 0);
        address messageBusProxy = deployMessageBus(eoaAdmin, feesProxy);
        deployPublicCallbacks(eoaAdmin);
        address crossChainMessengerProxy = deployCrossChainMessenger(eoaAdmin, messageBusProxy);
        deployEthereumBridge(eoaAdmin, crossChainMessengerProxy, remoteBridgeAddress);
        deployTenSystemCalls(eoaAdmin);
    }

    function deployAnalyzer(address eoaAdmin) internal {
        TransactionPostProcessor transactionsPostProcessor = new TransactionPostProcessor();
        bytes memory callData = abi.encodeWithSelector(transactionsPostProcessor.initialize.selector, eoaAdmin);
        address transactionsPostProcessorProxy = deployProxy(address(transactionsPostProcessor), eoaAdmin, callData);
        
        emit SystemContractDeployed("TransactionsPostProcessor", transactionsPostProcessorProxy);
    }

    function deployMessageBus(address eoaAdmin, address feesAddress) internal returns (address) {
        MessageBus messageBus = new MessageBus();
        bytes memory callData = abi.encodeWithSelector(messageBus.initialize.selector, eoaAdmin, feesAddress);
        address messageBusProxy = deployProxy(address(messageBus), eoaAdmin, callData);

        emit SystemContractDeployed("MessageBus", messageBusProxy);
        return messageBusProxy;
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

    function deployCrossChainMessenger(address eoaAdmin, address messageBusAddress) internal returns (address) {
        CrossChainMessenger crossChainMessenger = new CrossChainMessenger();
        bytes memory callData = abi.encodeWithSelector(crossChainMessenger.initialize.selector, messageBusAddress);
        address crossChainMessengerProxy = deployProxy(address(crossChainMessenger), eoaAdmin, callData);

        emit SystemContractDeployed("CrossChainMessenger", crossChainMessengerProxy);
        return crossChainMessengerProxy;
    }

    function deployEthereumBridge(address eoaAdmin, address crossChainMessengerAddress, address remoteBridgeAddress) internal {
        EthereumBridge ethereumBridge = new EthereumBridge();
        bytes memory callData = abi.encodeWithSelector(ethereumBridge.initialize.selector, crossChainMessengerAddress, remoteBridgeAddress);
        address ethereumBridgeProxy = deployProxy(address(ethereumBridge), eoaAdmin, callData);

        emit SystemContractDeployed("EthereumBridge", ethereumBridgeProxy);
    }

    function deployTenSystemCalls(address eoaAdmin) internal {
        TenSystemCalls tenSystemCalls = new TenSystemCalls();
        bytes memory callData = abi.encodeWithSelector(tenSystemCalls.initialize.selector);
        address tenSystemCallsProxy = deployProxy(address(tenSystemCalls), eoaAdmin, callData);

        emit SystemContractDeployed("TenSystemCalls", tenSystemCallsProxy);
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