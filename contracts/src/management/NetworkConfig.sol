// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "../lib/Storage.sol";

// Stores all the network addresses in one place
contract NetworkConfig is Initializable, OwnableUpgradeable {

    struct FixedAddresses {
        address crossChain;
        address messageBus;
        address networkEnclaveRegistry;
        address rollupContract;
    }

    struct NamedAddress {
        string name;
        address addr;
    }

    struct Addresses {
        address crossChain;
        address messageBus;
        address networkEnclaveRegistry;
        address rollupContract;
        NamedAddress[] additionalContracts;  // Dynamic address storage
    }

    // fixed storage slots for fixed contracts
    bytes32 public constant CROSS_CHAIN_SLOT = bytes32(uint256(keccak256("networkconfig.crossChain")) - 1);
    bytes32 public constant MESSAGE_BUS_SLOT = bytes32(uint256(keccak256("networkconfig.messageBus")) - 1);
    bytes32 public constant NETWORK_ENCLAVE_REGISTRY_SLOT = bytes32(uint256(keccak256("networkconfig.networkEnclaveRegistry")) - 1);
    bytes32 public constant ROLLUP_CONTRACT_SLOT = bytes32(uint256(keccak256("networkconfig.rollupContract")) - 1);

    // simple storage for additional addresses
    string[] private addressNames;
    mapping(string => address) public additionalAddresses;

    event NetworkContractAddressAdded(string name, address addr);

    function initialize( NetworkConfig.FixedAddresses memory _addresses, address owner) public initializer {
        __Ownable_init(owner);

        Storage.setAddress(CROSS_CHAIN_SLOT, _addresses.crossChain);
        Storage.setAddress(MESSAGE_BUS_SLOT, _addresses.messageBus);
        Storage.setAddress(NETWORK_ENCLAVE_REGISTRY_SLOT, _addresses.networkEnclaveRegistry);
        Storage.setAddress(ROLLUP_CONTRACT_SLOT, _addresses.rollupContract);
    }

    function crossChainContractAddress() public view returns (address addr_) {
        addr_ = Storage.getAddress(CROSS_CHAIN_SLOT);
    }

    function messageBusContractAddress() public view returns (address addr_) {
        addr_ = Storage.getAddress(MESSAGE_BUS_SLOT);
    }

    function networkEnclaveRegistryContractAddress() public view returns (address addr_) {
        addr_ = Storage.getAddress(NETWORK_ENCLAVE_REGISTRY_SLOT);
    }

    function rollupContractAddress() public view returns (address addr_) {
        addr_ = Storage.getAddress(ROLLUP_CONTRACT_SLOT);
    }

    // stores a new address in the simple address mapping
    function addAddress(string calldata name, address addr) external onlyOwner {
        require(addr != address(0), "Invalid address");
        require(additionalAddresses[name] == address(0), "Address name already exists");
        additionalAddresses[name] = addr;
        addressNames.push(name);
        emit NetworkContractAddressAdded(name, addr);
    }

    // returns all stored addresses
    function addresses() external view returns (Addresses memory) {
        NamedAddress[] memory additional = new NamedAddress[](addressNames.length);
        for(uint i = 0; i < addressNames.length; i++) {
            additional[i] = NamedAddress({
                name: addressNames[i],
                addr: additionalAddresses[addressNames[i]]
            });
        }

        return Addresses({
            networkEnclaveRegistry: networkEnclaveRegistryContractAddress(),
            crossChain: crossChainContractAddress(),
            messageBus: messageBusContractAddress(),
            rollupContract: rollupContractAddress(),
            additionalContracts: additional
        });
    }
}