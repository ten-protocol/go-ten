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
        address l1Bridge;
        address l2Bridge;
        address l1crossChainMessenger;
        address l2crossChainMessenger;
        NamedAddress[] additionalContracts;  // Dynamic address storage
    }

    // storage slots for fixed contracts
    bytes32 public constant CROSS_CHAIN_SLOT = bytes32(uint256(keccak256("networkconfig.crossChain")) - 1);
    bytes32 public constant MESSAGE_BUS_SLOT = bytes32(uint256(keccak256("networkconfig.messageBus")) - 1);
    bytes32 public constant NETWORK_ENCLAVE_REGISTRY_SLOT = bytes32(uint256(keccak256("networkconfig.networkEnclaveRegistry")) - 1);
    bytes32 public constant ROLLUP_CONTRACT_SLOT = bytes32(uint256(keccak256("networkconfig.rollupContract")) - 1);

    // storage slots for contracts that may need to be redeployed 
    bytes32 public constant L1_BRIDGE_SLOT = bytes32(uint256(keccak256("networkconfig.l1Bridge")) - 1);
    bytes32 public constant L2_BRIDGE_SLOT = bytes32(uint256(keccak256("networkconfig.l2Bridge")) - 1);
    bytes32 public constant L1_CROSS_CHAIN_MESSENGER = bytes32(uint256(keccak256("networkconfig.l1CrossChainMessenger")) - 1);
    bytes32 public constant L2_CROSS_CHAIN_MESSENGER = bytes32(uint256(keccak256("networkconfig.l2CrossChainMessenger")) - 1);

    // simple storage for additional addresses
    string[] public addressNames;
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

    // Add getters for bridge and messenger addresses
    function l1BridgeAddress() public view returns (address addr_) {
        addr_ = Storage.getAddress(L1_BRIDGE_SLOT);
    }

    function l2BridgeAddress() public view returns (address addr_) {
        addr_ = Storage.getAddress(L2_BRIDGE_SLOT);
    }

    function l1CrossChainMessengerAddress() public view returns (address addr_) {
        addr_ = Storage.getAddress(L1_CROSS_CHAIN_MESSENGER);
    }

    function l2CrossChainMessengerAddress() public view returns (address addr_) {
        addr_ = Storage.getAddress(L2_CROSS_CHAIN_MESSENGER);
    }

    // Setters for bridge and messenger addresses
    function setL1BridgeAddress(address _addr) external onlyOwner {
        require(_addr != address(0), "Invalid address");
        Storage.setAddress(L1_BRIDGE_SLOT, _addr);
        emit NetworkContractAddressAdded("l1Bridge", _addr);
    }

    function setL2BridgeAddress(address _addr) external onlyOwner {
        require(_addr != address(0), "Invalid address");
        Storage.setAddress(L2_BRIDGE_SLOT, _addr);
        emit NetworkContractAddressAdded("l2Bridge", _addr);
    }

    function setL1CrossChainMessengerAddress(address _addr) external onlyOwner {
        require(_addr != address(0), "Invalid address");
        Storage.setAddress(L1_CROSS_CHAIN_MESSENGER, _addr);
        emit NetworkContractAddressAdded("l1CrossChainMessenger", _addr);
    }

    function setL2CrossChainMessengerAddress(address _addr) external onlyOwner {
        require(_addr != address(0), "Invalid address");
        Storage.setAddress(L2_CROSS_CHAIN_MESSENGER, _addr);
        emit NetworkContractAddressAdded("l2CrossChainMessenger", _addr);
    }

    // stores or updates an address in the simple address mapping
    function addAdditionalAddress(string calldata name, address addr) external onlyOwner {
        require(addr != address(0), "Invalid address");
        if (additionalAddresses[name] == address(0)) {
            addressNames.push(name);
        }
        additionalAddresses[name] = addr;
        emit NetworkContractAddressAdded(name, addr);
    }


    function getAdditionaContractNames() public view returns (string[] memory) {
        return addressNames;
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
            l1Bridge: l1BridgeAddress(),
            l2Bridge: l2BridgeAddress(),
            l1crossChainMessenger: l1CrossChainMessengerAddress(),
            l2crossChainMessenger: l2CrossChainMessengerAddress(),
            additionalContracts: additional
        });
    }
}