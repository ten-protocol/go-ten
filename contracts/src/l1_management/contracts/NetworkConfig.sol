// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

import "../../lib/Storage.sol";
import "../../common/UnrenouncableOwnable2Step.sol";
/**
 * @title NetworkConfig
 * @dev Contract for managing network configuration and addresses
 * Implements a storage mechanism for fixed and dynamic addresses
 * Allows for adding and retrieving additional addresses
 */
contract NetworkConfig is Initializable, UnrenouncableOwnable2Step {

    /**
     * @dev Struct for system addresses that we don't expect to change
     */
    struct FixedAddresses {
        address crossChain;
        address messageBus;
        address networkEnclaveRegistry;
        address dataAvailabilityRegistry;
    }

    /**
     * @dev Struct for named addresses
     */
    struct NamedAddress {
        string name;
        address addr;
    }

    /**
     * @dev Struct for all addresses
     */
    struct Addresses {
        address crossChain;
        address messageBus;
        address networkEnclaveRegistry;
        address dataAvailabilityRegistry;
        address l1Bridge;
        address l2Bridge;
        address l1CrossChainMessenger;
        address l2CrossChainMessenger;
        NamedAddress[] additionalContracts;  // Dynamic address storage
    }

    /**
     * @dev Struct for contract version information
     */
    struct ContractVersion {
        string name;
        string version;
        address implementation;
    }

    // storage slots for fixed contracts
    bytes32 public constant CROSS_CHAIN_SLOT = bytes32(uint256(keccak256("networkconfig.crossChain")) - 1);
    bytes32 public constant MESSAGE_BUS_SLOT = bytes32(uint256(keccak256("networkconfig.messageBus")) - 1);
    bytes32 public constant NETWORK_ENCLAVE_REGISTRY_SLOT = bytes32(uint256(keccak256("networkconfig.networkEnclaveRegistry")) - 1);
    bytes32 public constant DATA_AVAILABILITY_REGISTRY_SLOT = bytes32(uint256(keccak256("networkconfig.dataAvailabilityRegistry")) - 1);

    // storage slots for contracts that may need to be redeployed 
    bytes32 public constant L1_BRIDGE_SLOT = bytes32(uint256(keccak256("networkconfig.l1Bridge")) - 1);
    bytes32 public constant L2_BRIDGE_SLOT = bytes32(uint256(keccak256("networkconfig.l2Bridge")) - 1);
    bytes32 public constant L1_CROSS_CHAIN_MESSENGER = bytes32(uint256(keccak256("networkconfig.l1CrossChainMessenger")) - 1);
    bytes32 public constant L2_CROSS_CHAIN_MESSENGER = bytes32(uint256(keccak256("networkconfig.l2CrossChainMessenger")) - 1);

    // simple storage for additional addresses
    string[] public addressNames;
    mapping(string contractName => address contractAddress) public additionalAddresses;

    /**
     * @dev Storage slot for the fork manager
     */
    bytes32 public constant FORK_MANAGER_SLOT = bytes32(uint256(keccak256("networkconfig.forkManager")) - 1);

    /**
     * @dev Event emitted when a network contract address is added
     * @param name The name of the contract
     * @param addr The address of the contract
     */
    event NetworkContractAddressAdded(string name, address addr);
    
    /**
     * @dev Event emitted when an additional contract address is added
     * @param name The name of the contract
     * @param addr The address of the contract
     */
    event AdditionalContractAddressAdded(string name, address addr);

    /**
     * @dev Event emitted when an additional contract address is removed
     * @param name The name of the contract
     */
    event AdditionalContractAddressRemoved(string name);

    /**
     * @dev Event emitted when a hardfork upgrade occurs
     * @param forkName The name of the hardfork
     */
    event HardforkUpgrade(
        string indexed forkName
    );

    /**
     * @dev Initializes the contract with addresses and owner
     * @param _addresses The fixed addresses
     * @param _owner Address of the contract owner
     */
    function initialize(NetworkConfig.FixedAddresses memory _addresses, address _owner) public initializer {
        require(_owner != address(0), "Owner cannot be 0x0");   
        require(_addresses.crossChain != address(0), "Cross chain cannot be 0x0");
        require(_addresses.messageBus != address(0), "Message bus cannot be 0x0");
        require(_addresses.networkEnclaveRegistry != address(0), "Network enclave registry cannot be 0x0");
        require(_addresses.dataAvailabilityRegistry != address(0), "Data availability registry cannot be 0x0");

        __UnrenouncableOwnable2Step_init(_owner);
        Storage.setAddress(CROSS_CHAIN_SLOT, _addresses.crossChain);
        Storage.setAddress(MESSAGE_BUS_SLOT, _addresses.messageBus);
        Storage.setAddress(NETWORK_ENCLAVE_REGISTRY_SLOT, _addresses.networkEnclaveRegistry);
        Storage.setAddress(DATA_AVAILABILITY_REGISTRY_SLOT, _addresses.dataAvailabilityRegistry);
    }

    /**
     * @dev Gets the cross chain contract address
     */
    function crossChainContractAddress() public view returns (address addr_) {
        addr_ = Storage.getAddress(CROSS_CHAIN_SLOT);
    }

    /**
     * @dev Gets the message bus contract address
     */
    function messageBusContractAddress() public view returns (address addr_) {
        addr_ = Storage.getAddress(MESSAGE_BUS_SLOT);
    }

    /**
     * @dev Gets the network enclave registry contract address
     */
    function networkEnclaveRegistryContractAddress() public view returns (address addr_) {
        addr_ = Storage.getAddress(NETWORK_ENCLAVE_REGISTRY_SLOT);
    }

    /**
     * @dev Gets the data availability registry contract address
     */     
    function daRegistryContractAddress() public view returns (address addr_) {
        addr_ = Storage.getAddress(DATA_AVAILABILITY_REGISTRY_SLOT);
    }

    /**
     * @dev Gets the L1 bridge contract address
     */
    function l1BridgeAddress() public view returns (address addr_) {
        addr_ = Storage.getAddress(L1_BRIDGE_SLOT);
    }

    /**
     * @dev Gets the L2 bridge contract address
     */
    function l2BridgeAddress() public view returns (address addr_) {
        addr_ = Storage.getAddress(L2_BRIDGE_SLOT);
    }

    /**
     * @dev Gets the L1 cross chain messenger contract address
     */
    function l1CrossChainMessengerAddress() public view returns (address addr_) {
        addr_ = Storage.getAddress(L1_CROSS_CHAIN_MESSENGER);
    }

    /**
     * @dev Gets the L2 cross chain messenger contract address
     */
    function l2CrossChainMessengerAddress() public view returns (address addr_) {
        addr_ = Storage.getAddress(L2_CROSS_CHAIN_MESSENGER);
    }

    /**
     * @dev Sets the L1 bridge contract address
     * @param _addr The address of the L1 bridge contract
     */
    function setL1BridgeAddress(address _addr) external onlyOwner {
        require(_addr != address(0), "Invalid address");
        Storage.setAddress(L1_BRIDGE_SLOT, _addr);
        emit NetworkContractAddressAdded("l1Bridge", _addr);
    }

    /**
     * @dev Sets the L2 bridge contract address
     * @param _addr The address of the L2 bridge contract
     */
    function setL2BridgeAddress(address _addr) external onlyOwner {
        require(_addr != address(0), "Invalid address");
        Storage.setAddress(L2_BRIDGE_SLOT, _addr);
        emit NetworkContractAddressAdded("l2Bridge", _addr);
    }

    /**
     * @dev Sets the L1 cross chain messenger contract address
     * @param _addr The address of the L1 cross chain messenger contract
     */
    function setL1CrossChainMessengerAddress(address _addr) external onlyOwner {
        require(_addr != address(0), "Invalid address");
        Storage.setAddress(L1_CROSS_CHAIN_MESSENGER, _addr);
        emit NetworkContractAddressAdded("l1CrossChainMessenger", _addr);
    }

    /**
     * @dev Sets the L2 cross chain messenger contract address
     * @param _addr The address of the L2 cross chain messenger contract
     */
    function setL2CrossChainMessengerAddress(address _addr) external onlyOwner {
        require(_addr != address(0), "Invalid address");
        Storage.setAddress(L2_CROSS_CHAIN_MESSENGER, _addr);
        emit NetworkContractAddressAdded("l2CrossChainMessenger", _addr);
    }

    /**
     * @dev Adds an additional contract address
     * @param name The name of the contract
     * @param addr The address of the contract
     */
    function addAdditionalAddress(string calldata name, address addr) external onlyOwner {
        require(addr != address(0), "Invalid address");
        require(additionalAddresses[name] == address(0), "Address already exists");
        require(bytes(name).length > 0, "Name cannot be empty");
        if (additionalAddresses[name] == address(0)) {
            addressNames.push(name);
        }
        additionalAddresses[name] = addr;
        emit AdditionalContractAddressAdded(name, addr);
    }

    function removeAdditionalAddress(string calldata name) external onlyOwner {
        require(additionalAddresses[name] != address(0), "Address does not exist");
        delete additionalAddresses[name];
        emit AdditionalContractAddressRemoved(name);
    }

    
    /**
     * @dev Gets the additional contract names
     * @return string[] The names of the additional contracts
     */
    function getAdditionalContractNames() public view returns (string[] memory) {
        return addressNames;
    }

    /**
     * @dev Gets all stored addresses
     * @return Addresses The addresses
     */
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
            dataAvailabilityRegistry: daRegistryContractAddress(),
            l1Bridge: l1BridgeAddress(),
            l2Bridge: l2BridgeAddress(),
            l1CrossChainMessenger: l1CrossChainMessengerAddress(),
            l2CrossChainMessenger: l2CrossChainMessengerAddress(),
            additionalContracts: additional
        });
    }

    /**
     * @dev Emits a hardfork upgrade event that can be subscribed to by the L2
     */
    function recordHardfork(
        string calldata hardforkName
    ) external onlyOwner {
        emit HardforkUpgrade(hardforkName);
    }
}