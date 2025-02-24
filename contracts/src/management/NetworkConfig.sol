// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "../libraries/Storage.sol";

// Stores all the network addresses in one place
contract NetworkConfig is Initializable, OwnableUpgradeable {

    struct Addresses {
        address crossChain;
        address messageBus;
        address networkEnclaveRegistry;
        address rollupContract;
    }

    bytes32 public constant CROSS_CHAIN_SLOT = bytes32(uint256(keccak256("networkconfig.crossChain")) - 1);

    bytes32 public constant MESSAGE_BUS_SLOT = bytes32(uint256(keccak256("networkconfig.messageBus")) - 1);

    bytes32 public constant NETWORK_ENCLAVE_REGISTRY_SLOT = bytes32(uint256(keccak256("networkconfig.networkEnclaveRegistry")) - 1);

    bytes32 public constant ROLLUP_CONTRACT_SLOT = bytes32(uint256(keccak256("networkconfig.rollupContract")) - 1);

    function initialize( NetworkConfig.Addresses memory _addresses) public initializer {
        __Ownable_init(msg.sender);

        Storage.setAddress(CROSS_CHAIN_SLOT, _addresses.crossChain);
        Storage.setAddress(MESSAGE_BUS_SLOT, _addresses.messageBus);
        Storage.setAddress(NETWORK_ENCLAVE_REGISTRY_SLOT, _addresses.networkEnclaveRegistry);
        Storage.setAddress(ROLLUP_CONTRACT_SLOT, _addresses.rollupContract);
//        emit LogManagementContractCreated(address(messageBus));
    }

    function crossChain() public view returns (address addr_) {
        addr_ = Storage.getAddress(CROSS_CHAIN_SLOT);
    }

    function messageBus() public view returns (address addr_) {
        addr_ = Storage.getAddress(MESSAGE_BUS_SLOT);
    }

    function networkEnclaveRegistry() public view returns (address addr_) {
        addr_ = Storage.getAddress(NETWORK_ENCLAVE_REGISTRY_SLOT);
    }

    function rollupContract() public view returns (address addr_) {
        addr_ = Storage.getAddress(ROLLUP_CONTRACT_SLOT);
    }

    function addresses() external view returns (Addresses memory) {
        return Addresses({
            crossChain: crossChain(),
            messageBus: messageBus(),
            networkEnclaveRegistry: networkEnclaveRegistry(),
            rollupContract: rollupContract()
        });
    }
}