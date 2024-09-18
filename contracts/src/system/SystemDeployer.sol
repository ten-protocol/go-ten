// SPDX-License-Identifier: Apache 2
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "./TransactionsAnalyzer.sol";

contract SystemDeployer {
    event SystemContractDeployed(string name, address contractAddress);

    constructor(address eoaAdmin) {
       deployAnalyzer(eoaAdmin);
    }

    function deployAnalyzer(address eoaAdmin) internal {
        TransactionsAnalyzer transactionsAnalyzer = new TransactionsAnalyzer();
        bytes memory callData = abi.encodeWithSelector(transactionsAnalyzer.initialize.selector, eoaAdmin, msg.sender);
        address transactionsAnalyzerProxy = deployProxy(address(transactionsAnalyzer), eoaAdmin, callData);
        
        emit SystemContractDeployed("TransactionsAnalyzer", transactionsAnalyzerProxy);
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