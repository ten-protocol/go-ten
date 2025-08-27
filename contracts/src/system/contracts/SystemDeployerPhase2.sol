// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.28;

import {CrossChainMessenger} from "../../cross_chain_messaging/common/CrossChainMessenger.sol";
import {EthereumBridge} from "../../reference_bridge/L2/contracts/EthereumBridge.sol";
import {MessageBus} from "../../cross_chain_messaging/common/MessageBus.sol";
import {TenSystemCalls} from "./TenSystemCalls.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

/**
 * @title SystemDeployerPhase2
 * @dev Deploys messaging and bridge contracts (Phase 2)
 */
contract SystemDeployerPhase2 {
    event SystemContractDeployed(string name, address contractAddress);

    struct DeployedContracts {
        address messageBus;
        address crossChainMessenger;
        address ethereumBridge;
        address tenSystemCalls;
    }

    DeployedContracts public deployedContracts;

    constructor(address eoaAdmin, address feesAddress, address remoteBridgeAddress) {
        require(eoaAdmin != address(0), "Invalid EOA admin address");
        require(feesAddress != address(0), "Invalid fees address");
        require(remoteBridgeAddress != address(0), "Invalid remote bridge address");
        deployedContracts.messageBus = deployMessageBus(eoaAdmin, feesAddress);
        deployedContracts.crossChainMessenger = deployCrossChainMessenger(eoaAdmin, deployedContracts.messageBus);
        deployedContracts.ethereumBridge = deployEthereumBridge(eoaAdmin, deployedContracts.crossChainMessenger, remoteBridgeAddress);
        deployedContracts.tenSystemCalls = deployTenSystemCalls(eoaAdmin);
    }

    function deployMessageBus(address eoaAdmin, address feesAddress) internal returns (address) {
        MessageBus messageBus = new MessageBus();
        bytes memory callData = abi.encodeWithSelector(messageBus.initialize.selector, eoaAdmin, eoaAdmin, feesAddress);
        address messageBusProxy = deployProxy(address(messageBus), eoaAdmin, callData);

        emit SystemContractDeployed("MessageBus", messageBusProxy);
        return messageBusProxy;
    }

    function deployCrossChainMessenger(address eoaAdmin, address messageBusAddress) internal returns (address) {
        CrossChainMessenger crossChainMessenger = new CrossChainMessenger();
        bytes memory callData = abi.encodeWithSelector(crossChainMessenger.initialize.selector, messageBusAddress);
        address crossChainMessengerProxy = deployProxy(address(crossChainMessenger), eoaAdmin, callData);

        emit SystemContractDeployed("CrossChainMessenger", crossChainMessengerProxy);
        return crossChainMessengerProxy;
    }

    function deployEthereumBridge(address eoaAdmin, address crossChainMessengerAddress, address remoteBridgeAddress) internal returns (address) {
        EthereumBridge ethereumBridge = new EthereumBridge();
        bytes memory callData = abi.encodeWithSelector(ethereumBridge.initialize.selector, crossChainMessengerAddress, remoteBridgeAddress);
        address ethereumBridgeProxy = deployProxy(address(ethereumBridge), eoaAdmin, callData);

        emit SystemContractDeployed("EthereumBridge", ethereumBridgeProxy);
        return ethereumBridgeProxy;
    }

    function deployTenSystemCalls(address eoaAdmin) internal returns (address) {
        TenSystemCalls tenSystemCalls = new TenSystemCalls();
        bytes memory callData = abi.encodeWithSelector(tenSystemCalls.initialize.selector);
        address tenSystemCallsProxy = deployProxy(address(tenSystemCalls), eoaAdmin, callData);

        emit SystemContractDeployed("TenSystemCalls", tenSystemCallsProxy);
        return tenSystemCallsProxy;
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
