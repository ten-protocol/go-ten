// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.28;

import "./TransactionPostProcessor.sol";
import "./Fees.sol";
import "./PublicCallbacks.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

/**
 * @title SystemDeployerPhase1
 * @dev Deploys core system contracts (Phase 1)
 */
contract SystemDeployerPhase1 {
    event SystemContractDeployed(string name, address contractAddress);

    struct DeployedContracts {
        address transactionPostProcessor;
        address fees;
        address publicCallbacks;
    }

    DeployedContracts public deployedContracts;

    constructor(address eoaAdmin) {
        require(eoaAdmin != address(0), "Invalid EOA admin address");
        deployedContracts.transactionPostProcessor = deployAnalyzer(eoaAdmin);
        deployedContracts.fees = deployFees(eoaAdmin, 0);
        deployedContracts.publicCallbacks = deployPublicCallbacks(eoaAdmin);
    }

    function deployAnalyzer(address eoaAdmin) internal returns (address) {
        TransactionPostProcessor transactionsPostProcessor = new TransactionPostProcessor();
        bytes memory callData = abi.encodeWithSelector(transactionsPostProcessor.initialize.selector, eoaAdmin);
        address transactionsPostProcessorProxy = deployProxy(address(transactionsPostProcessor), eoaAdmin, callData);
        
        emit SystemContractDeployed("TransactionsPostProcessor", transactionsPostProcessorProxy);
        return transactionsPostProcessorProxy;
    }

    function deployFees(address eoaAdmin, uint256 initialMessageFeePerByte) internal returns (address) {
        Fees fees = new Fees();
        bytes memory callData = abi.encodeWithSelector(fees.initialize.selector, initialMessageFeePerByte, eoaAdmin);
        address feesProxy = deployProxy(address(fees), eoaAdmin, callData);

        emit SystemContractDeployed("Fees", feesProxy);
        return feesProxy;
    }

    function deployPublicCallbacks(address eoaAdmin) internal returns (address) {
        PublicCallbacks publicCallbacks = new PublicCallbacks();
        bytes memory callData = abi.encodeWithSelector(publicCallbacks.initialize.selector);
        address publicCallbacksProxy = deployProxy(address(publicCallbacks), eoaAdmin, callData);

        emit SystemContractDeployed("PublicCallbacks", publicCallbacksProxy);
        return publicCallbacksProxy;
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
